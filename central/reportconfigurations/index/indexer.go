// Code generated by blevebindings generator. DO NOT EDIT.

package index

import (
	bleve "github.com/blevesearch/bleve"
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

type Indexer interface {
	AddReportConfiguration(reportconfiguration *storage.ReportConfiguration) error
	AddReportConfigurations(reportconfigurations []*storage.ReportConfiguration) error
	Count(q *aux.Query, opts ...blevesearch.SearchOption) (int, error)
	DeleteReportConfiguration(id string) error
	DeleteReportConfigurations(ids []string) error
	MarkInitialIndexingComplete() error
	NeedsInitialIndexing() (bool, error)
	Search(q *aux.Query, opts ...blevesearch.SearchOption) ([]search.Result, error)
}

func New(index bleve.Index) Indexer {
	return &indexerImpl{index: index}
}
