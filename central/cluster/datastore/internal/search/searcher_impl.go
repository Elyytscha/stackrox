package search

import (
	"context"
	"fmt"

	"github.com/stackrox/rox/central/cluster/index"
	"github.com/stackrox/rox/central/cluster/index/mappings"
	"github.com/stackrox/rox/central/cluster/store"
	cveSAC "github.com/stackrox/rox/central/cve/sac"
	"github.com/stackrox/rox/central/dackbox"
	deploymentSAC "github.com/stackrox/rox/central/deployment/sac"
	nsSAC "github.com/stackrox/rox/central/namespace/sac"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/dackbox/graph"
	"github.com/stackrox/rox/pkg/derivedfields/counter"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
	"github.com/stackrox/rox/pkg/search/derivedfields"
	"github.com/stackrox/rox/pkg/search/paginated"
)

var (
	defaultSortOption = &v1.QuerySortOption{
		Field:    search.Cluster.String(),
		Reversed: false,
	}

	clusterSearchHelper = sac.ForResource(resources.Cluster).MustCreateSearchHelper(mappings.OptionsMap)
)

type searcherImpl struct {
	storage           store.Store
	indexer           index.Indexer
	formattedSearcher search.Searcher
}

func (ds *searcherImpl) SearchResults(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	clusters, results, err := ds.searchClusters(ctx, q)
	if err != nil {
		return nil, err
	}
	protoResults := make([]*v1.SearchResult, 0, len(clusters))
	for i, alert := range clusters {
		protoResults = append(protoResults, convertCluster(alert, results[i]))
	}
	return protoResults, nil
}

func (ds *searcherImpl) SearchClusters(ctx context.Context, q *v1.Query) ([]*storage.Cluster, error) {
	clusters, _, err := ds.searchClusters(ctx, q)
	return clusters, err
}

func (ds *searcherImpl) searchClusters(ctx context.Context, q *v1.Query) ([]*storage.Cluster, []search.Result, error) {
	results, err := ds.Search(ctx, q)
	if err != nil {
		return nil, nil, err
	}

	clusters, missingIndices, err := ds.storage.GetSelectedClusters(search.ResultsToIDs(results))
	if err != nil {
		return nil, nil, err
	}
	results = search.RemoveMissingResults(results, missingIndices)
	return clusters, results, nil
}

func (ds *searcherImpl) Search(ctx context.Context, q *v1.Query) ([]search.Result, error) {
	return ds.formattedSearcher.Search(ctx, q)
}

func convertCluster(cluster *storage.Cluster, result search.Result) *v1.SearchResult {
	return &v1.SearchResult{
		Category:       v1.SearchCategory_CLUSTERS,
		Id:             cluster.GetId(),
		Name:           cluster.GetName(),
		FieldToMatches: search.GetProtoMatchesMap(result.Matches),
		Score:          result.Score,
		Location:       fmt.Sprintf("/%s", cluster.GetName()),
	}
}

// Helper functions which format our searching.
///////////////////////////////////////////////

func formatSearcher(unsafeSearcher blevesearch.UnsafeSearcher, graphProvider graph.Provider) search.Searcher {
	filteredSearcher := clusterSearchHelper.FilteredSearcher(unsafeSearcher) // Make the UnsafeSearcher safe.
	derivedFieldSortedSearcher := wrapDerivedFieldSearcher(graphProvider, filteredSearcher)
	paginatedSearcher := paginated.Paginated(derivedFieldSortedSearcher)
	defaultSortedSearcher := paginated.WithDefaultSortOption(paginatedSearcher, defaultSortOption)
	return defaultSortedSearcher
}

func wrapDerivedFieldSearcher(graphProvider graph.Provider, searcher search.Searcher) search.Searcher {
	if !features.Dackbox.Enabled() {
		return searcher
	}
	return derivedfields.CountSortedSearcher(searcher, map[string]counter.DerivedFieldCounter{
		search.NamespaceCount.String():  counter.NewGraphBasedDerivedFieldCounter(graphProvider, dackbox.ClusterToNamespace, nsSAC.GetSACFilter()),
		search.DeploymentCount.String(): counter.NewGraphBasedDerivedFieldCounter(graphProvider, dackbox.ClusterToDeployment, deploymentSAC.GetSACFilter()),
		search.CVECount.String():        counter.NewGraphBasedDerivedFieldCounter(graphProvider, dackbox.ClusterToCVE, cveSAC.GetSACFilter()),
	})
}
