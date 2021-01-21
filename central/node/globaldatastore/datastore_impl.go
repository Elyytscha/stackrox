package globaldatastore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/node/datastore"
	"github.com/stackrox/rox/central/node/globalstore"
	"github.com/stackrox/rox/central/node/index"
	"github.com/stackrox/rox/central/node/index/mappings"
	"github.com/stackrox/rox/central/ranking"
	riskDS "github.com/stackrox/rox/central/risk/datastore"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
)

var (
	log = logging.LoggerForModule()

	nodesSAC             = sac.ForResource(resources.Node)
	nodesSACSearchHelper = nodesSAC.MustCreateSearchHelper(mappings.OptionsMap)
)

type globalDataStore struct {
	globalStore globalstore.GlobalStore

	indexer index.Indexer

	risks               riskDS.DataStore
	nodeRanker          *ranking.Ranker
	nodeComponentRanker *ranking.Ranker
}

// New creates and returns a new GlobalDataStore.
func New(globalStore globalstore.GlobalStore, indexer index.Indexer, risks riskDS.DataStore,
	nodeRanker *ranking.Ranker, nodeComponentRanker *ranking.Ranker) (GlobalDataStore, error) {
	gds := &globalDataStore{
		globalStore: globalStore,
		indexer:     indexer,

		risks:               risks,
		nodeRanker:          nodeRanker,
		nodeComponentRanker: nodeComponentRanker,
	}
	if err := gds.buildIndex(); err != nil {
		return nil, err
	}

	if features.HostScanning.Enabled() {
		if err := gds.initializeRankers(); err != nil {
			return nil, err
		}
	}

	return gds, nil
}

func (s *globalDataStore) buildIndex() error {
	nodeStores, err := s.globalStore.GetAllClusterNodeStores()
	if err != nil {
		return err
	}

	for _, nodeStore := range nodeStores {
		nodes, err := nodeStore.ListNodes()
		if err != nil {
			return err
		}
		if err := s.indexer.AddNodes(nodes); err != nil {
			return err
		}
	}
	return nil
}

func (s *globalDataStore) initializeRankers() error {
	stores, err := s.globalStore.GetAllClusterNodeStores()
	if err != nil {
		return errors.Wrap(err, "initializing node rankers")
	}

	for _, store := range stores {
		nodes, err := store.ListNodes()
		if err != nil {
			return errors.Wrap(err, "retrieving nodes from store for ranker initialization")
		}

		for _, node := range nodes {
			s.nodeRanker.Add(node.GetId(), node.GetRiskScore())
		}
	}

	return nil
}

func (s *globalDataStore) GetAllClusterNodeStores(ctx context.Context, writeAccess bool) (map[string]datastore.DataStore, error) {
	stores, err := s.globalStore.GetAllClusterNodeStores()
	if err != nil {
		return nil, err
	}

	accessMode := storage.Access_READ_ACCESS
	if writeAccess {
		accessMode = storage.Access_READ_WRITE_ACCESS
	}

	if ok, err := nodesSAC.AccessAllowed(ctx, accessMode); err != nil {
		return nil, err
	} else if !ok {
		scopeChecker := nodesSAC.ScopeChecker(ctx, accessMode)
		// Pass 1: Mark requests for all clusters as pending
		for clusterID := range stores {
			scopeChecker.TryAllowed(sac.ClusterScopeKey(clusterID))
		}
		if err := scopeChecker.PerformChecks(ctx); err != nil {
			return nil, err
		}
		// Pass 2: Filter out clusters for which we have no access.
		for clusterID := range stores {
			if scopeChecker.TryAllowed(sac.ClusterScopeKey(clusterID)) != sac.Allow {
				delete(stores, clusterID)
			}
		}
	}

	dataStores := make(map[string]datastore.DataStore, len(stores))
	for clusterID, store := range stores {
		dataStores[clusterID] = datastore.New(store, s.indexer, writeAccess, s.risks, s.nodeRanker, s.nodeComponentRanker)
	}

	return dataStores, nil
}

func (s *globalDataStore) GetClusterNodeStore(ctx context.Context, clusterID string, writeAccess bool) (datastore.DataStore, error) {
	accessMode := storage.Access_READ_ACCESS
	if writeAccess {
		accessMode = storage.Access_READ_WRITE_ACCESS
	}

	if ok, err := nodesSAC.AccessAllowed(ctx, accessMode, sac.ClusterScopeKey(clusterID)); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("permission denied")
	}
	store, err := s.globalStore.GetClusterNodeStore(clusterID, writeAccess)
	if err != nil {
		return nil, err
	}
	return datastore.New(store, s.indexer, writeAccess, s.risks, s.nodeRanker, s.nodeComponentRanker), nil
}

func (s *globalDataStore) RemoveClusterNodeStores(ctx context.Context, clusterIDs ...string) error {
	// Stop early, otherwise, the clusterID query turns into an empty query that matches all
	if len(clusterIDs) == 0 {
		return nil
	}

	if ok, err := nodesSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}

	if err := s.globalStore.RemoveClusterNodeStores(clusterIDs...); err != nil {
		return errors.Wrap(err, "removing cluster node stores")
	}

	q := search.NewQueryBuilder().AddExactMatches(search.ClusterID, clusterIDs...).ProtoQuery()
	results, err := s.indexer.Search(q)
	if err != nil {
		return errors.Wrap(err, "searching for nodes")
	}
	if err := s.indexer.DeleteNodes(search.ResultsToIDs(results)); err != nil {
		return errors.Wrap(err, "deleting nodes from indexer")
	}
	return nil
}

func (s *globalDataStore) CountAllNodes(ctx context.Context) (int, error) {
	if ok, err := nodesSAC.ReadAllowed(ctx); err != nil {
		return 0, err
	} else if ok {
		return s.globalStore.CountAllNodes()
	}

	searchResults, err := s.Search(ctx, search.EmptyQuery())
	if err != nil {
		return 0, err
	}
	return len(searchResults), nil
}

// SearchResults returns any node matches to the query
func (s *globalDataStore) SearchResults(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	// We do the filtering in the search, so OK to operate on store directly.
	stores, err := s.globalStore.GetAllClusterNodeStores()
	if err != nil {
		return nil, err
	}

	results, err := s.Search(ctx, q)
	if err != nil {
		return nil, err
	}

	searchResults := make([]*v1.SearchResult, 0, len(results))
	for _, r := range results {
		var node *storage.Node
		for _, s := range stores {
			node, err = s.GetNode(r.ID)
			if err == nil {
				break
			}
		}
		if node == nil {
			continue
		}
		searchResults = append(searchResults, &v1.SearchResult{
			Id:             r.ID,
			Name:           node.Name,
			Category:       v1.SearchCategory_NODES,
			FieldToMatches: search.GetProtoMatchesMap(r.Matches),
			Score:          r.Score,
			Location:       fmt.Sprintf("%s/%s", node.GetClusterName(), node.GetName()),
		})
	}
	return searchResults, nil
}

// SearchRawNodes returns nodes that match a query
func (s *globalDataStore) SearchRawNodes(ctx context.Context, q *v1.Query) ([]*storage.Node, error) {
	// We do the filtering in the search, so OK to operate on store directly.
	stores, err := s.globalStore.GetAllClusterNodeStores()
	if err != nil {
		return nil, err
	}

	results, err := s.Search(ctx, q)
	if err != nil {
		return nil, err
	}

	nodes := make([]*storage.Node, 0, len(results))
	for _, r := range results {
		var node *storage.Node
		for _, store := range stores {
			node, err = store.GetNode(r.ID)
			if err == nil && node != nil {
				nodes = append(nodes, node)
				break
			}
		}
	}
	return nodes, nil
}

// Search returns any node matches to the query
func (s *globalDataStore) Search(ctx context.Context, q *v1.Query) ([]search.Result, error) {
	return nodesSACSearchHelper.Apply(s.indexer.Search)(ctx, q)
}
