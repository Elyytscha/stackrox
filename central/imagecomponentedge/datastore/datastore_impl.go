package datastore

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/imagecomponentedge/index"
	"github.com/stackrox/rox/central/imagecomponentedge/search"
	"github.com/stackrox/rox/central/imagecomponentedge/store"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/sac"
	searchPkg "github.com/stackrox/rox/pkg/search"
)

var (
	imagesSAC = sac.ForResource(resources.Image)
)

type datastoreImpl struct {
	storage  store.Store
	indexer  index.Indexer
	searcher search.Searcher
}

func (ds *datastoreImpl) Search(ctx context.Context, q *v1.Query) ([]searchPkg.Result, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, err
	}
	return ds.searcher.Search(ctx, q)
}

func (ds *datastoreImpl) SearchEdges(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, err
	}
	return ds.searcher.SearchEdges(ctx, q)
}

func (ds *datastoreImpl) SearchRawEdges(ctx context.Context, q *v1.Query) ([]*storage.ImageComponentEdge, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, err
	}
	edges, err := ds.searcher.SearchRawEdges(ctx, q)
	if err != nil {
		return nil, err
	}
	return edges, nil
}

func (ds *datastoreImpl) Count(ctx context.Context) (int, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); !ok || err != nil {
		return 0, err
	}
	return ds.storage.Count()
}

func (ds *datastoreImpl) Get(ctx context.Context, id string) (*storage.ImageComponentEdge, bool, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, false, err
	}
	edge, found, err := ds.storage.Get(id)
	if err != nil || !found {
		return nil, false, err
	}
	return edge, true, nil
}

func (ds *datastoreImpl) Exists(ctx context.Context, id string) (bool, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); err != nil || !ok {
		return false, err
	}
	found, err := ds.storage.Exists(id)
	if err != nil || !found {
		return false, err
	}
	return true, nil
}

func (ds *datastoreImpl) GetBatch(ctx context.Context, ids []string) ([]*storage.ImageComponentEdge, error) {
	if ok, err := imagesSAC.ReadAllowed(ctx); err != nil || !ok {
		return nil, err
	}
	edges, _, err := ds.storage.GetBatch(ids)
	if err != nil {
		return nil, err
	}
	return edges, nil
}

func (ds *datastoreImpl) Upsert(ctx context.Context, edge *storage.ImageComponentEdge) error {
	if edge.GetId() == "" {
		return errors.New("cannot upsert an image/component edge without an id")
	}
	if ok, err := imagesSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}

	if err := ds.indexer.AddImageComponentEdge(edge); err != nil {
		return err
	}
	return ds.storage.Upsert(edge)
}

func (ds *datastoreImpl) Delete(ctx context.Context, ids ...string) error {
	if ok, err := imagesSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}

	errorList := errorhelpers.NewErrorList("deleting image/component edges")
	for _, id := range ids {
		if err := ds.storage.Delete(id); err != nil {
			errorList.AddError(err)
			continue
		}
		if err := ds.indexer.DeleteImageComponentEdge(id); err != nil {
			errorList.AddError(err)
		}
	}
	return errorList.ToError()
}
