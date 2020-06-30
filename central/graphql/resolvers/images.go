package resolvers

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/graphql/resolvers/loaders"
	"github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	pkgMetrics "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/scoped"
	"github.com/stackrox/rox/pkg/utils"
)

func init() {
	schema := getBuilder()
	utils.Must(
		schema.AddQuery("images(query: String, pagination: Pagination): [Image!]!"),
		schema.AddQuery("imageCount(query: String): Int!"),
		schema.AddQuery("image(sha:ID!): Image"),
		schema.AddExtraResolver("Image", "deployments(query: String, pagination: Pagination): [Deployment!]!"),
		schema.AddExtraResolver("Image", "deploymentCount(query: String): Int!"),
		schema.AddExtraResolver("Image", "topVuln(query: String): EmbeddedVulnerability"),
		schema.AddExtraResolver("Image", "vulns(query: String, pagination: Pagination): [EmbeddedVulnerability]!"),
		schema.AddExtraResolver("Image", "vulnCount(query: String): Int!"),
		schema.AddExtraResolver("Image", "vulnCounter(query: String): VulnerabilityCounter!"),
		schema.AddExtraResolver("EmbeddedImageScanComponent", "layerIndex: Int"),
		schema.AddExtraResolver("Image", "components(query: String, pagination: Pagination): [EmbeddedImageScanComponent!]!"),
		schema.AddExtraResolver("Image", `componentCount(query: String): Int!`),
		schema.AddExtraResolver("Image", `unusedVarSink(query: String): Int`),
		schema.AddExtraResolver("Image", "plottedVulns(query: String): PlottedVulnerabilities!"),
	)
}

// Images returns GraphQL resolvers for all images
func (resolver *Resolver) Images(ctx context.Context, args PaginatedQuery) ([]*imageResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "Images")
	if err := readImages(ctx); err != nil {
		return nil, err
	}

	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}
	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return nil, err
	}
	return resolver.wrapImages(imageLoader.FromQuery(ctx, q))
}

// ImageCount returns count of all images across deployments
func (resolver *Resolver) ImageCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "ImageCount")
	if err := readImages(ctx); err != nil {
		return 0, err
	}

	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return 0, err
	}
	return imageLoader.CountFromQuery(ctx, q)
}

// Image returns a graphql resolver for the identified image, if it exists
func (resolver *Resolver) Image(ctx context.Context, args struct{ Sha graphql.ID }) (*imageResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "Image")
	if err := readImages(ctx); err != nil {
		return nil, err
	}

	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return nil, err
	}
	image, err := imageLoader.FromID(ctx, string(args.Sha))
	return resolver.wrapImage(image, image != nil, err)
}

// Deployments returns the deployments which use this image for the identified image, if it exists
func (resolver *imageResolver) Deployments(ctx context.Context, args PaginatedQuery) ([]*deploymentResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Images, "Deployments")
	if err := readDeployments(ctx); err != nil {
		return nil, err
	}

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getImageRawQuery())

	return resolver.root.Deployments(ctx, PaginatedQuery{Pagination: args.Pagination, Query: &query})
}

// Deployments returns the deployments which use this image for the identified image, if it exists
func (resolver *imageResolver) DeploymentCount(ctx context.Context, args RawQuery) (int32, error) {
	if err := readDeployments(ctx); err != nil {
		return 0, err
	}

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getImageRawQuery())

	return resolver.root.DeploymentCount(ctx, RawQuery{Query: &query})
}

// TopVuln returns the first vulnerability with the top CVSS score.
func (resolver *imageResolver) TopVuln(ctx context.Context, args RawQuery) (VulnerabilityResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Images, "TopVulnerability")
	return resolver.topVulnV2(ctx, args)
}

func (resolver *imageResolver) topVulnV2(ctx context.Context, args RawQuery) (VulnerabilityResolver, error) {
	query, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}

	if resolver.data.GetSetTopCvss() == nil {
		return nil, nil
	}

	query = search.NewConjunctionQuery(query, resolver.getImageQuery())
	query.Pagination = &v1.QueryPagination{
		SortOptions: []*v1.QuerySortOption{
			{
				Field:    search.CVSS.String(),
				Reversed: true,
			},
			{
				Field:    search.CVE.String(),
				Reversed: true,
			},
		},
		Limit:  1,
		Offset: 0,
	}

	vulnLoader, err := loaders.GetCVELoader(ctx)
	if err != nil {
		return nil, err
	}
	vulns, err := vulnLoader.FromQuery(ctx, query)
	if err != nil {
		return nil, err
	} else if len(vulns) == 0 {
		return nil, err
	} else if len(vulns) > 1 {
		return nil, errors.New("multiple vulnerabilities matched for top image vulnerability")
	}
	return &cVEResolver{root: resolver.root, data: vulns[0]}, nil

}

// Vulns returns all of the vulnerabilities in the image.
func (resolver *imageResolver) Vulns(ctx context.Context, args PaginatedQuery) ([]VulnerabilityResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Images, "Vulnerabilities")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getImageRawQuery())

	return resolver.root.Vulnerabilities(scoped.Context(ctx, scoped.Scope{
		Level: v1.SearchCategory_IMAGES,
		ID:    resolver.data.GetId(),
	}), PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

// VulnCount returns the number of vulnerabilities the image has.
func (resolver *imageResolver) VulnCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Images, "VulnerabilityCount")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getImageRawQuery())

	return resolver.root.VulnerabilityCount(scoped.Context(ctx, scoped.Scope{
		Level: v1.SearchCategory_IMAGES,
		ID:    resolver.data.GetId(),
	}), RawQuery{Query: &query})
}

// VulnCounter resolves the number of different types of vulnerabilities contained in an image component.
func (resolver *imageResolver) VulnCounter(ctx context.Context, args RawQuery) (*VulnerabilityCounterResolver, error) {
	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getImageRawQuery())

	return resolver.root.VulnCounter(scoped.Context(ctx, scoped.Scope{
		Level: v1.SearchCategory_IMAGES,
		ID:    resolver.data.GetId(),
	}), RawQuery{Query: &query})
}

// Vulns returns all of the vulnerabilities in the image.
func (resolver *imageResolver) Components(ctx context.Context, args PaginatedQuery) ([]ComponentResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Images, "ImageComponents")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getImageRawQuery())

	return resolver.root.Components(scoped.Context(ctx, scoped.Scope{
		Level: v1.SearchCategory_IMAGES,
		ID:    resolver.data.GetId(),
	}), PaginatedQuery{Query: &query, Pagination: args.Pagination})
}

func (resolver *imageResolver) ComponentCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Cluster, "ComponentCount")

	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getImageRawQuery())

	return resolver.root.ComponentCount(scoped.Context(ctx, scoped.Scope{
		Level: v1.SearchCategory_IMAGES,
		ID:    resolver.data.GetId(),
	}), RawQuery{Query: &query})
}

func (resolver *Resolver) getImage(ctx context.Context, id string) *storage.Image {
	imageLoader, err := loaders.GetImageLoader(ctx)
	if err != nil {
		return nil
	}
	image, err := imageLoader.FromID(ctx, id)
	if err != nil {
		return nil
	}
	return image
}

func (resolver *imageResolver) getImageRawQuery() string {
	return search.NewQueryBuilder().AddExactMatches(search.ImageSHA, resolver.data.GetId()).Query()
}

func (resolver *imageResolver) getImageQuery() *v1.Query {
	return search.NewQueryBuilder().AddExactMatches(search.ImageSHA, resolver.data.GetId()).ProtoQuery()
}

func (resolver *imageResolver) PlottedVulns(ctx context.Context, args RawQuery) (*PlottedVulnerabilitiesResolver, error) {
	query := search.AddRawQueriesAsConjunction(args.String(), resolver.getImageRawQuery())
	return newPlottedVulnerabilitiesResolver(ctx, resolver.root, RawQuery{Query: &query})
}

func (resolver *imageResolver) UnusedVarSink(ctx context.Context, args RawQuery) *int32 {
	return nil
}
