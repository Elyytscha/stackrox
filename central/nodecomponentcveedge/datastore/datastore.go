package datastore

import (
	"context"

	"github.com/stackrox/rox/central/nodecomponentcveedge/datastore/index"
	"github.com/stackrox/rox/central/nodecomponentcveedge/datastore/search"
	"github.com/stackrox/rox/central/nodecomponentcveedge/datastore/store/postgres"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/aux"
	"github.com/stackrox/rox/generated/storage"
	searchPkg "github.com/stackrox/rox/pkg/search"
)

// DataStore is an intermediary to Component/CVE edge storage.
//go:generate mockgen-wrapper
type DataStore interface {
	Search(ctx context.Context, q *aux.Query) ([]searchPkg.Result, error)
	SearchEdges(ctx context.Context, q *aux.Query) ([]*v1.SearchResult, error)
	SearchRawEdges(ctx context.Context, q *aux.Query) ([]*storage.NodeComponentCVEEdge, error)

	Exists(ctx context.Context, id string) (bool, error)
	Get(ctx context.Context, id string) (*storage.NodeComponentCVEEdge, bool, error)
	Count(ctx context.Context, q *aux.Query) (int, error)
}

// New returns a new instance of a DataStore.
func New(storage postgres.Store, indexer index.Indexer, searcher search.Searcher) (DataStore, error) {
	ds := &datastoreImpl{
		storage:  storage,
		indexer:  indexer,
		searcher: searcher,
	}
	return ds, nil
}
