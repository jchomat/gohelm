package gohelm

import (
	"context"
	"fmt"

	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
)

func (c *Client) ListReleasesByStatus(ctx context.Context, status []release.Status_Code) ([]*release.Release, error) {

	var allReleases []*release.Release
	var offset string

	sv := services.NewReleaseServiceClient(c.Conn)

	for {
		// Get Helm releases with specific status
		realReq := &services.ListReleasesRequest{
			StatusCodes: status,
			Limit:       10,
			Offset:      offset,
		}
		res, err := sv.ListReleases(ctx, realReq)
		if err != nil {
			return allReleases, err
		}
		rec, err := res.Recv()
		if err != nil {
			return allReleases, err
		}

		offset = rec.Releases[len(rec.Releases)-1].GetName()

		allReleases = append(allReleases, rec.Releases...)

		fmt.Printf("Offset: %s\n", offset)

		if rec.Count < 10 {
			break
		}

		select {
		case <-ctx.Done():
			break
		default:
			continue
		}
	}

	return allReleases, nil
}
