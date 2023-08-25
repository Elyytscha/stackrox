// Code generated by pg-bindings generator. DO NOT EDIT.
package postgres

import (
	"context"
	"time"

	metrics "github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres"
	search "github.com/stackrox/rox/pkg/search"
	pgSearch "github.com/stackrox/rox/pkg/search/postgres"
)

// NewIndexer returns new indexer for `storage.NodeComponentEdge`.
func NewIndexer(db postgres.DB) *indexerImpl {
	return &indexerImpl{
		db: db,
	}
}

type indexerImpl struct {
	db postgres.DB
}

func (b *indexerImpl) Count(ctx context.Context, q *v1.Query) (int, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Count, "NodeComponentEdge")

	return pgSearch.RunCountRequest(ctx, v1.SearchCategory_NODE_COMPONENT_EDGE, q, b.db)
}

func (b *indexerImpl) Search(ctx context.Context, q *v1.Query) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "NodeComponentEdge")

	return pgSearch.RunSearchRequest(ctx, v1.SearchCategory_NODE_COMPONENT_EDGE, q, b.db)
}
