// Code generated by blevebindings generator. DO NOT EDIT.

package index

import (
	"bytes"
	"time"

	bleve "github.com/blevesearch/bleve"
	metrics "github.com/stackrox/rox/central/metrics"
	mappings "github.com/stackrox/rox/central/rbac/k8srole/mappings"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/aux"
	storage "github.com/stackrox/rox/generated/storage"
	batcher "github.com/stackrox/rox/pkg/batcher"
	ops "github.com/stackrox/rox/pkg/metrics"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

const batchSize = 5000

const resourceName = "K8SRole"

type indexerImpl struct {
	index bleve.Index
}

type k8SRoleWrapper struct {
	*storage.K8SRole `json:"k8s_role"`
	Type             string `json:"type"`
}

func (b *indexerImpl) AddK8SRole(k8srole *storage.K8SRole) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Add, "K8SRole")
	if err := b.index.Index(k8srole.GetId(), &k8SRoleWrapper{
		K8SRole: k8srole,
		Type:    v1.SearchCategory_ROLES.String(),
	}); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) AddK8SRoles(k8sroles []*storage.K8SRole) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.AddMany, "K8SRole")
	batchManager := batcher.New(len(k8sroles), batchSize)
	for {
		start, end, ok := batchManager.Next()
		if !ok {
			break
		}
		if err := b.processBatch(k8sroles[start:end]); err != nil {
			return err
		}
	}
	return nil
}

func (b *indexerImpl) processBatch(k8sroles []*storage.K8SRole) error {
	batch := b.index.NewBatch()
	for _, k8srole := range k8sroles {
		if err := batch.Index(k8srole.GetId(), &k8SRoleWrapper{
			K8SRole: k8srole,
			Type:    v1.SearchCategory_ROLES.String(),
		}); err != nil {
			return err
		}
	}
	return b.index.Batch(batch)
}

func (b *indexerImpl) Count(q *aux.Query, opts ...blevesearch.SearchOption) (int, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Count, "K8SRole")
	return blevesearch.RunCountRequest(v1.SearchCategory_ROLES, q, b.index, mappings.OptionsMap, opts...)
}

func (b *indexerImpl) DeleteK8SRole(id string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Remove, "K8SRole")
	if err := b.index.Delete(id); err != nil {
		return err
	}
	return nil
}

func (b *indexerImpl) DeleteK8SRoles(ids []string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.RemoveMany, "K8SRole")
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
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "K8SRole")
	return blevesearch.RunSearchRequest(v1.SearchCategory_ROLES, q, b.index, mappings.OptionsMap, opts...)
}
