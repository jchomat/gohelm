package gohelm

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
)

func (c *Client) ListReleasesByStatus(ctx context.Context, status []release.Status_Code) ([]*release.Release, error) {

	var allReleases []*release.Release
	var offset string

	// Config helm version header
	md := metadata.Pairs("x-helm-api-client", c.Version)
	helmCtx := metadata.NewOutgoingContext(ctx, md)

	sv := services.NewReleaseServiceClient(c.Conn)

	for {
		// Get Helm releases with specific status
		realReq := &services.ListReleasesRequest{
			StatusCodes: status,
			Limit:       10,
			Offset:      offset,
		}
		res, err := sv.ListReleases(helmCtx, realReq)
		if err != nil {
			return allReleases, err
		}
		rec, err := res.Recv()
		if err != nil {
			return allReleases, err
		}

		allReleases = append(allReleases, rec.Releases...)

		if rec.Count < 10 {
			break
		}

		offset = rec.Releases[len(rec.Releases)-1].GetName()
		fmt.Printf("Offset: %s\n", offset)

		select {
		case <-helmCtx.Done():
			break
		default:
			continue
		}
	}

	return allReleases, nil
}
