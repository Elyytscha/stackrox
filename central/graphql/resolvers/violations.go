package resolvers

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	pkgMetrics "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/utils"
)

func init() {
	schema := getBuilder()
	utils.Must(
		schema.AddQuery("violations(query: String, pagination: Pagination): [Alert!]!"),
		schema.AddQuery("violationCount(query: String): Int!"),
		schema.AddQuery("violation(id: ID!): Alert"),
	)
}

// Violations returns a list of all violations, or those that match the requested query
func (resolver *Resolver) Violations(ctx context.Context, args paginatedQuery) ([]*alertResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "Violations")
	if err := readAlerts(ctx); err != nil {
		return nil, err
	}
	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}
	return resolver.wrapListAlerts(
		resolver.ViolationsDataStore.SearchListAlerts(ctx, q))
}

// ViolationCount returns count of all violations, or those that match the requested query
func (resolver *Resolver) ViolationCount(ctx context.Context, args rawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "ViolationCount")
	if err := readAlerts(ctx); err != nil {
		return 0, err
	}
	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	results, err := resolver.ViolationsDataStore.Search(ctx, q)
	if err != nil {
		return 0, err
	}
	return int32(len(results)), nil
}

// Violation returns the violation with the requested ID
func (resolver *Resolver) Violation(ctx context.Context, args struct{ graphql.ID }) (*alertResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "Violation")
	if err := readAlerts(ctx); err != nil {
		return nil, err
	}
	return resolver.wrapAlert(
		resolver.ViolationsDataStore.GetAlert(ctx, string(args.ID)))
}

func (resolver *Resolver) getAlert(ctx context.Context, id string) *storage.Alert {
	alert, ok, err := resolver.ViolationsDataStore.GetAlert(ctx, id)
	if err != nil || !ok {
		return nil
	}
	return alert
}

func getLatestViolationTime(ctx context.Context, root *Resolver, q *v1.Query) (*graphql.Time, error) {
	if err := readAlerts(ctx); err != nil {
		return nil, err
	}

	q.Pagination = &v1.QueryPagination{
		SortOptions: []*v1.QuerySortOption{
			{
				Field:    search.ViolationTime.String(),
				Reversed: true,
			},
		},
		Limit:  1,
		Offset: 0,
	}

	alerts, err := root.ViolationsDataStore.SearchRawAlerts(ctx, q)
	if err != nil || len(alerts) == 0 || alerts[0] == nil {
		return nil, err
	}

	return timestamp(alerts[0].GetTime())
}
