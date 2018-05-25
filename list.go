package gohelm

import (
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
			Limit:       10,
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

		allReleases = append(allReleases, rec.Releases...)

		if rec.Count < 10 {
			break
		}
	}

	return allReleases, nil
}
