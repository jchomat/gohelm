package gohelm

import (
	"fmt"

	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
)

func (c *Client) DeleteRelease(name string) error {
	sv := services.NewReleaseServiceClient(c.Conn)

	uniReq := &services.UninstallReleaseRequest{
		Name:  release.GetName(),
		Purge: true,
	}
	_, err := sv.UninstallRelease(c.Context, uniReq)

	return err
}

func (c *Client) DeleteReleasesInNamespace(namespace string) error {
	// Get Helm release in current namespace
	realReq := &services.ListReleasesRequest{
		SortBy:    services.ListSort_NAME,
		SortOrder: services.ListSort_ASC,
		Namespace: namespace,
	}
	sv := services.NewReleaseServiceClient(c.Conn)
	res, err := sv.ListReleases(c.Context, realReq)
	if err != nil {
		return err
	}
	rec, err := res.Recv()
	if err != nil {
		return err
	}

	if rec.Count == 0 {
		return nil
	}

	// Delete releases in current namespace
	for _, release := range rec.Releases {
		uniReq := &services.UninstallReleaseRequest{
			Name:  release.GetName(),
			Purge: true,
		}
		_, err := sv.UninstallRelease(c.Context, uniReq)
		if err != nil {
			fmt.Errorf("Failed to uninstall release %s: %s\n", release.GetName(), err)
			continue
		}
	}

	return nil
}
