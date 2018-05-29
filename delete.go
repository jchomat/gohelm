package gohelm

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
	"k8s.io/helm/pkg/proto/hapi/services"
)

func (c *Client) DeleteRelease(ctx context.Context, name string) error {
	sv := services.NewReleaseServiceClient(c.Conn)

	// Config helm version header
	md := metadata.Pairs("x-helm-api-client", c.Version)
	helmCtx := metadata.NewOutgoingContext(ctx, md)

	uniReq := &services.UninstallReleaseRequest{
		Name:  name,
		Purge: true,
	}
	_, err := sv.UninstallRelease(helmCtx, uniReq)

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
