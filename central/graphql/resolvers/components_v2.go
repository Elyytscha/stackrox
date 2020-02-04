package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/graphql/resolvers/loaders"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
)

// Top Level Resolvers.
///////////////////////

func (resolver *Resolver) componentV2(ctx context.Context, args idQuery) (ComponentResolver, error) {
	return resolver.wrapImageComponent(resolver.ImageComponentDataStore.Get(ctx, string(*args.ID)))
}

func (resolver *Resolver) componentsV2(ctx context.Context, args PaginatedQuery) ([]ComponentResolver, error) {
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}
	return resolver.componentsV2Query(ctx, query)
}

func (resolver *Resolver) componentsV2Query(ctx context.Context, query *v1.Query) ([]ComponentResolver, error) {
	componentLoader, err := loaders.GetComponentLoader(ctx)
	if err != nil {
		return nil, err
	}

	compRes, err := resolver.wrapImageComponents(componentLoader.FromQuery(ctx, query))
	if err != nil {
		return nil, err
	}

	ret := make([]ComponentResolver, 0, len(compRes))
	for _, resolver := range compRes {
		ret = append(ret, resolver)
	}
	return ret, err
}

func (resolver *Resolver) componentCountV2(ctx context.Context, args RawQuery) (int32, error) {
	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	return resolver.componentCountV2Query(ctx, q)
}

func (resolver *Resolver) componentCountV2Query(ctx context.Context, query *v1.Query) (int32, error) {
	componentLoader, err := loaders.GetComponentLoader(ctx)
	if err != nil {
		return 0, err
	}

	return componentLoader.CountFromQuery(ctx, query)
}

// Resolvers on Component Object.
/////////////////////////////////

// ID returns a unique identifier for the component. Need to implement this on top of 'Id' so that we can implement
// the same interface as the non-generated embedded resolver used in v1.
func (eicr *imageComponentResolver) ID(ctx context.Context) graphql.ID {
	return graphql.ID(eicr.data.GetId())
}

// LastScanned is the last time the component was scanned in an image.
func (eicr *imageComponentResolver) LastScanned(ctx context.Context) (*graphql.Time, error) {
	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return nil, err
	}

	componentQuery := eicr.componentQuery()
	componentQuery.Pagination = &v1.QueryPagination{
		Limit:  1,
		Offset: 0,
		SortOptions: []*v1.QuerySortOption{
			{
				Field:    search.ImageScanTime.String(),
				Reversed: true,
			},
		},
	}

	images, err := imageLoader.FromQuery(ctx, componentQuery)
	if err != nil {
		return nil, err
	} else if len(images) == 0 {
		return nil, nil
	} else if len(images) > 1 {
		return nil, errors.New("multiple images matched for last scanned component query")
	}

	return timestamp(images[0].GetScan().GetScanTime())
}

// TopVuln returns the first vulnerability with the top CVSS score.
func (eicr *imageComponentResolver) TopVuln(ctx context.Context) (VulnerabilityResolver, error) {
	vulnLoader, err := loaders.GetCVELoader(ctx)
	if err != nil {
		return nil, err
	}
	vulns, err := vulnLoader.FromQuery(ctx, eicr.componentQuery())
	if err != nil {
		return nil, err
	}
	var topVuln *storage.CVE
	for _, vuln := range vulns {
		if topVuln == nil || vuln.GetCvss() > topVuln.GetCvss() {
			topVuln = vuln
		}
	}
	if topVuln == nil {
		return nil, err
	}
	return &cVEResolver{
		root: eicr.root,
		data: topVuln,
	}, nil
}

// Vulns resolves the vulnerabilities contained in the image component.
func (eicr *imageComponentResolver) Vulns(ctx context.Context, args PaginatedQuery) ([]VulnerabilityResolver, error) {
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}
	pagination := query.GetPagination()
	query, err = search.AddAsConjunction(eicr.componentQuery(), query)
	if err != nil {
		return nil, err
	}
	query.Pagination = pagination
	return eicr.root.vulnerabilitiesV2Query(ctx, query)
}

// VulnCount resolves the number of vulnerabilities contained in the image component.
func (eicr *imageComponentResolver) VulnCount(ctx context.Context, args RawQuery) (int32, error) {
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	query, err = search.AddAsConjunction(eicr.componentQuery(), query)
	if err != nil {
		return 0, err
	}
	return eicr.root.vulnerabilityCountV2Query(ctx, query)
}

// VulnCounter resolves the number of different types of vulnerabilities contained in an image component.
func (eicr *imageComponentResolver) VulnCounter(ctx context.Context, args RawQuery) (*VulnerabilityCounterResolver, error) {
	vulnLoader, err := loaders.GetCVELoader(ctx)
	if err != nil {
		return nil, err
	}
	q1 := search.NewConjunctionQuery(eicr.componentQuery(), search.NewQueryBuilder().AddBools(search.Fixable, true).ProtoQuery())
	fixableVulns, err := vulnLoader.FromQuery(ctx, q1)
	if err != nil {
		return nil, err
	}
	allVulns, err := vulnLoader.FromQuery(ctx, eicr.componentQuery())
	if err != nil {
		return nil, err
	}
	return mapVCEsToVulnerabilityCounter(fixableVulns, allVulns), nil
}

// Images are the images that contain the Component.
func (eicr *imageComponentResolver) Images(ctx context.Context, args PaginatedQuery) ([]*imageResolver, error) {
	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return nil, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}
	pagination := query.GetPagination()
	query, err = search.AddAsConjunction(eicr.componentQuery(), query)
	if err != nil {
		return nil, err
	}
	query.Pagination = pagination
	return eicr.root.wrapImages(imageLoader.FromQuery(ctx, query))
}

// ImageCount is the number of images that contain the Component.
func (eicr *imageComponentResolver) ImageCount(ctx context.Context, args RawQuery) (int32, error) {
	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return 0, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	query, err = search.AddAsConjunction(eicr.componentQuery(), query)
	if err != nil {
		return 0, err
	}
	return imageLoader.CountFromQuery(ctx, query)
}

// Deployments are the deployments that contain the Component.
func (eicr *imageComponentResolver) Deployments(ctx context.Context, args PaginatedQuery) ([]*deploymentResolver, error) {
	if err := readDeployments(ctx); err != nil {
		return nil, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	pagination := query.GetPagination()
	query, err = search.AddAsConjunction(eicr.componentQuery(), query)
	if err != nil {
		return nil, err
	}
	query.Pagination = pagination

	deploymentLoader, err := loaders.GetDeploymentLoader(ctx)
	if err != nil {
		return nil, err
	}
	return eicr.root.wrapDeployments(deploymentLoader.FromQuery(ctx, query))
}

// DeploymentCount is the number of deployments that contain the Component.
func (eicr *imageComponentResolver) DeploymentCount(ctx context.Context, args RawQuery) (int32, error) {
	if err := readDeployments(ctx); err != nil {
		return 0, err
	}
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	query, err = search.AddAsConjunction(eicr.componentQuery(), query)
	if err != nil {
		return 0, err
	}
	deploymentLoader, err := loaders.GetDeploymentLoader(ctx)
	if err != nil {
		return 0, err
	}
	return deploymentLoader.CountFromQuery(ctx, query)
}

// Helper functions.
////////////////////

func (eicr *imageComponentResolver) componentQuery() *v1.Query {
	return search.NewQueryBuilder().
		AddExactMatches(search.ComponentName, eicr.data.GetName()).
		AddExactMatches(search.ComponentVersion, eicr.data.GetVersion()).
		ProtoQuery()
}

// These return dummy values, as they should not be accessed from the top level component resolver, but the embedded
// version instead.

// Location returns the location of the component.
func (eicr *imageComponentResolver) Location(ctx context.Context) string {
	return ""
}

// LayerIndex is the index in the parent image.
func (eicr *imageComponentResolver) LayerIndex() *int32 {
	return nil
}
