package search

import (
	"context"

	componentCVEEdgeIndexer "github.com/stackrox/rox/central/componentcveedge/index"
	cveIndexer "github.com/stackrox/rox/central/cve/index"
	"github.com/stackrox/rox/central/image/datastore/internal/store"
	imageIndexer "github.com/stackrox/rox/central/image/index"
	componentIndexer "github.com/stackrox/rox/central/imagecomponent/index"
	imageComponentEdgeIndexer "github.com/stackrox/rox/central/imagecomponentedge/index"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/idspace"
)

// Searcher provides search functionality on existing alerts
//go:generate mockgen-wrapper
type Searcher interface {
	SearchImages(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error)
	SearchRawImages(ctx context.Context, q *v1.Query) ([]*storage.Image, error)
	SearchListImages(ctx context.Context, q *v1.Query) ([]*storage.ListImage, error)

	Search(ctx context.Context, q *v1.Query) ([]search.Result, error)
}

// New returns a new instance of Searcher for the given storage and indexer.
func New(storage store.Store, graphProvider idspace.GraphProvider,
	cveIndexer cveIndexer.Indexer,
	componentCVEEdgeIndexer componentCVEEdgeIndexer.Indexer,
	componentIndexer componentIndexer.Indexer,
	imageComponentEdgeIndexer imageComponentEdgeIndexer.Indexer,
	imageIndexer imageIndexer.Indexer) Searcher {
	return &searcherImpl{
		storage: storage,
		indexer: imageIndexer,
		searcher: formatSearcher(graphProvider,
			cveIndexer,
			componentCVEEdgeIndexer,
			componentIndexer,
			imageComponentEdgeIndexer,
			imageIndexer),
	}
}
