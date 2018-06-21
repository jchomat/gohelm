package gohelm

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc/metadata"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
)

func (c *Client) ListReleasesByStatus(ctx context.Context, status []release.Status_Code) ([]*release.Release, error) {

	// Get Helm releases with specific status
	realReq := &services.ListReleasesRequest{
		StatusCodes: status,
	}
	sv := services.NewReleaseServiceClient(c.Conn)

	// Config helm version header
	md := metadata.Pairs("x-helm-api-client", c.Version)
	helmCtx := metadata.NewOutgoingContext(ctx, md)

	res, err := sv.ListReleases(helmCtx, realReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list release: %s", err)
	}

	var releases []*release.Release
	for {
		rec, err := res.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read release list grpc response: %s", err)
		}

		releases = append(releases, rec.Releases...)

		select {
		case <-helmCtx.Done():
			break
		default:
			continue
		}
	}

	return releases, nil
}
