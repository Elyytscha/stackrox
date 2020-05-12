package datastore

import (
	"context"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger"
	"github.com/stackrox/rox/central/analystnotes"
	componentCVEEdgeIndexer "github.com/stackrox/rox/central/componentcveedge/index"
	cveIndexer "github.com/stackrox/rox/central/cve/index"
	"github.com/stackrox/rox/central/deployment/datastore/internal/processtagsstore"
	"github.com/stackrox/rox/central/deployment/datastore/internal/search"
	"github.com/stackrox/rox/central/deployment/index"
	"github.com/stackrox/rox/central/deployment/store"
	badgerStore "github.com/stackrox/rox/central/deployment/store/badger"
	"github.com/stackrox/rox/central/deployment/store/cache"
	dackBoxStore "github.com/stackrox/rox/central/deployment/store/dackbox"
	imageDS "github.com/stackrox/rox/central/image/datastore"
	imageIndexer "github.com/stackrox/rox/central/image/index"
	componentIndexer "github.com/stackrox/rox/central/imagecomponent/index"
	imageComponentEdgeIndexer "github.com/stackrox/rox/central/imagecomponentedge/index"
	nfDS "github.com/stackrox/rox/central/networkflow/datastore"
	piDS "github.com/stackrox/rox/central/processindicator/datastore"
	pwDS "github.com/stackrox/rox/central/processwhitelist/datastore"
	"github.com/stackrox/rox/central/ranking"
	riskDS "github.com/stackrox/rox/central/risk/datastore"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/dackbox"
	"github.com/stackrox/rox/pkg/dackbox/graph"
	"github.com/stackrox/rox/pkg/expiringcache"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/process/filter"
	pkgSearch "github.com/stackrox/rox/pkg/search"
)

// DataStore is an intermediary to AlertStorage.
//go:generate mockgen-wrapper
type DataStore interface {
	Search(ctx context.Context, q *v1.Query) ([]pkgSearch.Result, error)
	SearchDeployments(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error)
	SearchRawDeployments(ctx context.Context, q *v1.Query) ([]*storage.Deployment, error)
	SearchListDeployments(ctx context.Context, q *v1.Query) ([]*storage.ListDeployment, error)

	ListDeployment(ctx context.Context, id string) (*storage.ListDeployment, bool, error)

	GetDeployment(ctx context.Context, id string) (*storage.Deployment, bool, error)
	GetDeployments(ctx context.Context, ids []string) ([]*storage.Deployment, error)
	CountDeployments(ctx context.Context) (int, error)
	// UpsertDeployment adds or updates a deployment. If the deployment exists, the tags in the deployment are taken from
	// the stored deployment.
	UpsertDeployment(ctx context.Context, deployment *storage.Deployment) error

	// UpsertDeploymentIntoStoreOnly does not index the data on insertion. If the deployment exists, the tags in the deployment
	// are taken from the stored deployment.
	UpsertDeploymentIntoStoreOnly(ctx context.Context, deployment *storage.Deployment) error

	AddTagsToProcessKey(ctx context.Context, key *analystnotes.ProcessNoteKey, tags []string) error
	RemoveTagsFromProcessKey(ctx context.Context, key *analystnotes.ProcessNoteKey, tags []string) error
	GetTagsForProcessKey(ctx context.Context, key *analystnotes.ProcessNoteKey) ([]string, error)

	RemoveDeployment(ctx context.Context, clusterID, id string) error

	GetImagesForDeployment(ctx context.Context, deployment *storage.Deployment) ([]*storage.Image, error)
	GetDeploymentIDs() ([]string, error)
}

func newDataStore(storage store.Store, graphProvider graph.Provider, processTagsStore processtagsstore.Store, bleveIndex bleve.Index, processIndex bleve.Index,
	images imageDS.DataStore, indicators piDS.DataStore, whitelists pwDS.DataStore, networkFlows nfDS.ClusterDataStore,
	risks riskDS.DataStore, deletedDeploymentCache expiringcache.Cache, processFilter filter.Filter,
	clusterRanker *ranking.Ranker, nsRanker *ranking.Ranker, deploymentRanker *ranking.Ranker) (DataStore, error) {
	var searcher search.Searcher
	indexer := index.New(bleveIndex, processIndex)

	storage = cache.NewCachedStore(storage)
	if features.Dackbox.Enabled() {
		searcher = search.New(storage,
			graphProvider,
			cveIndexer.New(bleveIndex),
			componentCVEEdgeIndexer.New(bleveIndex),
			componentIndexer.New(bleveIndex),
			imageComponentEdgeIndexer.New(bleveIndex),
			imageIndexer.New(bleveIndex),
			indexer)
	} else {
		searcher = search.New(storage, nil, nil, nil, nil, nil, nil, indexer)
	}

	ds, err := newDatastoreImpl(storage, processTagsStore, indexer, searcher, images, indicators, whitelists, networkFlows, risks,
		deletedDeploymentCache, processFilter, clusterRanker, nsRanker, deploymentRanker)

	if err != nil {
		return nil, err
	}

	ds.initializeRanker()
	return ds, nil
}

// NewBadger creates a deployment datastore based on BadgerDB
func NewBadger(dacky *dackbox.DackBox, keyFence concurrency.KeyFence, db *badger.DB, processTagsStore processtagsstore.Store, bleveIndex bleve.Index, processIndex bleve.Index,
	images imageDS.DataStore, indicators piDS.DataStore, whitelists pwDS.DataStore, networkFlows nfDS.ClusterDataStore,
	risks riskDS.DataStore, deletedDeploymentCache expiringcache.Cache, processFilter filter.Filter,
	clusterRanker *ranking.Ranker, nsRanker *ranking.Ranker, deploymentRanker *ranking.Ranker) (DataStore, error) {
	var err error
	var storage store.Store
	if features.Dackbox.Enabled() {
		storage, err = dackBoxStore.New(dacky, keyFence)
		if err != nil {
			return nil, err
		}
	} else {
		storage, err = badgerStore.New(db)
		if err != nil {
			return nil, err
		}
	}

	return newDataStore(storage, dacky, processTagsStore, bleveIndex, processIndex, images, indicators, whitelists, networkFlows, risks, deletedDeploymentCache, processFilter, clusterRanker, nsRanker, deploymentRanker)
}
