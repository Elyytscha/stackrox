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
	mapping.RegisterCategoryToTable(v1.SearchCategory_SEARCH_UNSET, schema)
}

// NewIndexer returns new indexer for `storage.TestMultiKeyStruct`.
func NewIndexer(db *pgxpool.Pool) *indexerImpl {
	return &indexerImpl{
		db: db,
	}
}

type indexerImpl struct {
	db *pgxpool.Pool
}

func (b *indexerImpl) Count(q *aux.Query, opts ...blevesearch.SearchOption) (int, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Count, "TestMultiKeyStruct")

	return postgres.RunCountRequest(v1.SearchCategory_SEARCH_UNSET, q, b.db)
}

func (b *indexerImpl) Search(q *aux.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "TestMultiKeyStruct")

	return postgres.RunSearchRequest(v1.SearchCategory_SEARCH_UNSET, q, b.db)
}

//// Stubs for satisfying interfaces

func (b *indexerImpl) AddTestMultiKeyStruct(deployment *storage.TestMultiKeyStruct) error {
	return nil
}

func (b *indexerImpl) AddTestMultiKeyStructs(_ []*storage.TestMultiKeyStruct) error {
	return nil
}

func (b *indexerImpl) DeleteTestMultiKeyStruct(id string) error {
	return nil
}

func (b *indexerImpl) DeleteTestMultiKeyStructs(_ []string) error {
	return nil
}

func (b *indexerImpl) MarkInitialIndexingComplete() error {
	return nil
}

func (b *indexerImpl) NeedsInitialIndexing() (bool, error) {
	return false, nil
}
