// Code generated by pg-bindings generator. DO NOT EDIT.
package postgres

import (
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	metrics "github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	ops "github.com/stackrox/rox/pkg/metrics"
	search "github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
	"github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/search/postgres/mapping"
)

func init() {
	mapping.RegisterCategoryToTable(v1.SearchCategory_POLICY_CATEGORIES, schema)
}

// NewIndexer returns new indexer for `storage.PolicyCategory`.
func NewIndexer(db *pgxpool.Pool) *indexerImpl {
	return &indexerImpl{
		db: db,
	}
}

type indexerImpl struct {
	db *pgxpool.Pool
}

func (b *indexerImpl) Count(q *aux.Query, opts ...blevesearch.SearchOption) (int, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Count, "PolicyCategory")

	return postgres.RunCountRequest(v1.SearchCategory_POLICY_CATEGORIES, q, b.db)
}

func (b *indexerImpl) Search(q *aux.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "PolicyCategory")

	return postgres.RunSearchRequest(v1.SearchCategory_POLICY_CATEGORIES, q, b.db)
}

//// Stubs for satisfying interfaces

func (b *indexerImpl) AddPolicyCategory(deployment *storage.PolicyCategory) error {
	return nil
}

func (b *indexerImpl) AddPolicyCategories(_ []*storage.PolicyCategory) error {
	return nil
}

func (b *indexerImpl) DeletePolicyCategory(id string) error {
	return nil
}

func (b *indexerImpl) DeletePolicyCategories(_ []string) error {
	return nil
}

func (b *indexerImpl) MarkInitialIndexingComplete() error {
	return nil
}

func (b *indexerImpl) NeedsInitialIndexing() (bool, error) {
	return false, nil
}
