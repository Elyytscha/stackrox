// Code generated by blevebindings generator. DO NOT EDIT.

package index

import (
	"bytes"
	"time"

	bleve "github.com/blevesearch/bleve"
	metrics "github.com/stackrox/rox/central/metrics"
	mappings "github.com/stackrox/rox/central/processbaseline/index/mappings"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	batcher "github.com/stackrox/rox/pkg/batcher"
	ops "github.com/stackrox/rox/pkg/metrics"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

const batchSize = 5000

const resourceName = "ProcessBaseline"

type indexerImpl struct {
	index bleve.Index
}

type processBaselineWrapper struct {
	*storage.ProcessBaseline `json:"process_baseline"`
	Type                     string `json:"type"`
}

func (b *indexerImpl) AddProcessBaseline(processbaseline *storage.ProcessBaseline) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Add, "ProcessBaseline")
	if err := b.index.Index(processbaseline.GetId(), &processBaselineWrapper{
		ProcessBaseline: processbaseline,
		Type:            v1.SearchCategory_PROCESS_BASELINES.String(),
	}); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) AddProcessBaselines(processbaselines []*storage.ProcessBaseline) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.AddMany, "ProcessBaseline")
	batchManager := batcher.New(len(processbaselines), batchSize)
	for {
		start, end, ok := batchManager.Next()
		if !ok {
			break
		}
		if err := b.processBatch(processbaselines[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (b *indexerImpl) processBatch(processbaselines []*storage.ProcessBaseline) error {
	batch := b.index.NewBatch()
	for _, processbaseline := range processbaselines {
		if err := batch.Index(processbaseline.GetId(), &processBaselineWrapper{
			ProcessBaseline: processbaseline,
			Type:            v1.SearchCategory_PROCESS_BASELINES.String(),
		}); err != nil {
			return err
		}
	}
	return b.index.Batch(batch)
}

func (b *indexerImpl) Count(q *aux.Query, opts ...blevesearch.SearchOption) (int, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Count, "ProcessBaseline")
	return blevesearch.RunCountRequest(v1.SearchCategory_PROCESS_BASELINES, q, b.index, mappings.OptionsMap, opts...)
}

func (b *indexerImpl) DeleteProcessBaseline(id string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Remove, "ProcessBaseline")
	if err := b.index.Delete(id); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) DeleteProcessBaselines(ids []string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.RemoveMany, "ProcessBaseline")
	batch := b.index.NewBatch()
	for _, id := range ids {
		batch.Delete(id)
	}
	if err := b.index.Batch(batch); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) MarkInitialIndexingComplete() error {
	return b.index.SetInternal([]byte(resourceName), []byte("old"))
}

func (b *indexerImpl) NeedsInitialIndexing() (bool, error) {
	data, err := b.index.GetInternal([]byte(resourceName))
	if err != nil {
		return false, err
	}
	return !bytes.Equal([]byte("old"), data), nil
}

func (b *indexerImpl) Search(q *aux.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "ProcessBaseline")
	return blevesearch.RunSearchRequest(v1.SearchCategory_PROCESS_BASELINES, q, b.index, mappings.OptionsMap, opts...)
}
