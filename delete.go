package gohelm

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc/metadata"
	"k8s.io/helm/pkg/proto/hapi/release"
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

func (c *Client) DeleteReleasesInNamespace(ctx context.Context, namespace string) error {
	// Get Helm release in current namespace
	realReq := &services.ListReleasesRequest{
		SortBy:    services.ListSort_NAME,
		SortOrder: services.ListSort_ASC,
		Namespace: namespace,
	}
	sv := services.NewReleaseServiceClient(c.Conn)

	// Config helm version header
	md := metadata.Pairs("x-helm-api-client", c.Version)
	helmCtx := metadata.NewOutgoingContext(ctx, md)

	res, err := sv.ListReleases(helmCtx, realReq)
	if err != nil {
		return fmt.Errorf("failed to list release: %s", err)
	}

	var releases []*release.Release
	for {
		rec, err := res.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to read release list grpc response: %s", err)
		}

		releases = append(releases, rec.GetReleases())
	}

	if len(releases) == 0 {
		return nil
	}

	// Delete releases in current namespace
	for _, release := range releases {
		uniReq := &services.UninstallReleaseRequest{
			Name:  release.GetName(),
			Purge: true,
		}
		log.Printf("Will delete %s release", release.GetName())
		//		_, err := sv.UninstallRelease(helmCtx, uniReq)
		//		if err != nil {
		//			fmt.Errorf("Failed to uninstall release %s: %s\n", release.GetName(), err)
		//			continue
		//		}
	}

	return nil
}
