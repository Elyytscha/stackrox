package datastore

import (
	"context"

	"github.com/blevesearch/bleve"
	bolt "github.com/etcd-io/bbolt"
	"github.com/stackrox/rox/central/alert/datastore/internal/index"
	"github.com/stackrox/rox/central/alert/datastore/internal/search"
	"github.com/stackrox/rox/central/alert/datastore/internal/store"
	"github.com/stackrox/rox/central/globaldb"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	searchPkg "github.com/stackrox/rox/pkg/search"
)

// DataStore is a transaction script with methods that provide the domain logic for CRUD uses cases for Alert objects.
//go:generate mockgen-wrapper
type DataStore interface {
	Search(ctx context.Context, q *v1.Query) ([]searchPkg.Result, error)
	SearchAlerts(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error)
	SearchRawAlerts(ctx context.Context, q *v1.Query) ([]*storage.Alert, error)
	SearchListAlerts(ctx context.Context, q *v1.Query) ([]*storage.ListAlert, error)

	ListAlerts(ctx context.Context, request *v1.ListAlertsRequest) ([]*storage.ListAlert, error)

	GetAlert(ctx context.Context, id string) (*storage.Alert, bool, error)
	CountAlerts(ctx context.Context) (int, error)
	AddAlert(ctx context.Context, alert *storage.Alert) error
	UpdateAlert(ctx context.Context, alert *storage.Alert) error
	MarkAlertStale(ctx context.Context, id string) error

	DeleteAlerts(ctx context.Context, ids ...string) error
}

// New returns a new soleInstance of DataStore using the input store, indexer, and searcher.
func New(storage store.Store, indexer index.Indexer, searcher search.Searcher) (DataStore, error) {
	ds := &datastoreImpl{
		storage:    storage,
		indexer:    indexer,
		searcher:   searcher,
		keyedMutex: concurrency.NewKeyedMutex(globaldb.DefaultDataStorePoolSize),
	}
	if err := ds.buildIndex(); err != nil {
		return nil, err
	}
	return ds, nil
}

// NewWithDb returns a new soleInstance of DataStore using the input indexer, and searcher.
func NewWithDb(db *bolt.DB, bIndex bleve.Index) DataStore {
	store := store.New(db)
	indexer := index.New(bIndex)
	searcher := search.New(store, indexer)

	return &datastoreImpl{
		storage:    store,
		indexer:    indexer,
		searcher:   searcher,
		keyedMutex: concurrency.NewKeyedMutex(globaldb.DefaultDataStorePoolSize),
	}
}
