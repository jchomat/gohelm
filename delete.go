package gohelm

import (
	"log"

	tiller "k8s.io/helm/pkg/proto/hapi/services"
)

func (c *Client) DeleteReleasesInNamespace(namespace string) error {
	// Get Helm release in current namespace
	realReq := &tiller.ListReleasesRequest{
		SortBy:    tiller.ListSort_NAME,
		SortOrder: tiller.ListSort_ASC,
		Namespace: namespace,
	}
	sv := tiller.NewReleaseServiceClient(c.Conn)
	res, err := sv.ListReleases(c.Context, realReq)
	if err != nil {
		return err
	}
	rec, err := res.Recv()
	if err != nil {
		return err
	}

	if rec.Count == 0 {
		log.Printf("No helm release to delete in namespace %s\n", namespace)
	}

	// Delete releases in current namespace
	for _, release := range rec.Releases {
		uniReq := &tiller.UninstallReleaseRequest{
			Name:  release.GetName(),
			Purge: true,
		}
		_, err := sv.UninstallRelease(c.Context, uniReq)
		if err != nil {
			log.Printf("Failed to uninstall release %s: %s\n", release.GetName(), err)
			continue
		}
		log.Printf("Release %s uninstalled\n", release.GetName())
	}

	return nil
}
