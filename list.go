package gohelm

import (
	"fmt"

	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
)

func (c *Client) ListReleasesByStatus(status int) ([]*release.Release, error) {

	var allReleases []*release.Release
	var offset string

	sv := services.NewReleaseServiceClient(c.Conn)

	for {
		// Get Helm releases with specific status
		realReq := &services.ListReleasesRequest{
			StatusCodes: []release.Status_Code{release.Status_DELETED},
			Limit:       20,
			Offset:      offset,
		}
		res, err := sv.ListReleases(c.Context, realReq)
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

		if len(allReleases) > 200 {
			break
		}

		if rec.Count < 10 {
			break
		}
	}

	return allReleases, nil
}
