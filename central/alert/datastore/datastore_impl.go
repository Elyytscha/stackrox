package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/stackrox/rox/central/alert/convert"
	"github.com/stackrox/rox/central/alert/datastore/internal/index"
	"github.com/stackrox/rox/central/alert/datastore/internal/search"
	"github.com/stackrox/rox/central/alert/datastore/internal/store"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/batcher"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/debug"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sac"
	searchCommon "github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/paginated"
	"github.com/stackrox/rox/pkg/txn"
)

var (
	log = logging.LoggerForModule()

	alertSAC = sac.ForResource(resources.Alert)
)

const (
	alertBatchSize = 1000
)

// datastoreImpl is a transaction script with methods that provide the domain logic for CRUD uses cases for Alert
// objects.
type datastoreImpl struct {
	storage    store.Store
	indexer    index.Indexer
	searcher   search.Searcher
	keyedMutex *concurrency.KeyedMutex
}

func (ds *datastoreImpl) Search(ctx context.Context, q *v1.Query) ([]searchCommon.Result, error) {
	return ds.searcher.Search(ctx, q)
}

func (ds *datastoreImpl) SearchListAlerts(ctx context.Context, q *v1.Query) ([]*storage.ListAlert, error) {
	return ds.searcher.SearchListAlerts(ctx, q)
}

// SearchAlerts returns search results for the given request.
func (ds *datastoreImpl) SearchAlerts(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	return ds.searcher.SearchAlerts(ctx, q)
}

// SearchRawAlerts returns search results for the given request in the form of a slice of alerts.
func (ds *datastoreImpl) SearchRawAlerts(ctx context.Context, q *v1.Query) ([]*storage.Alert, error) {
	return ds.searcher.SearchRawAlerts(ctx, q)
}

func (ds *datastoreImpl) ListAlerts(ctx context.Context, request *v1.ListAlertsRequest) ([]*storage.ListAlert, error) {
	var q *v1.Query
	if request.GetQuery() == "" {
		q = searchCommon.EmptyQuery()
	} else {
		var err error
		q, err = searchCommon.ParseRawQuery(request.GetQuery())
		if err != nil {
			return nil, err
		}
	}

	paginated.FillPagination(q, request.GetPagination(), alertBatchSize)

	alerts, err := ds.SearchListAlerts(ctx, q)
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

// GetAlert returns an alert by id.
func (ds *datastoreImpl) GetAlert(ctx context.Context, id string) (*storage.Alert, bool, error) {
	alert, exists, err := ds.storage.GetAlert(id)
	if err != nil || !exists {
		return nil, false, err
	}

	if ok, err := alertSAC.ReadAllowed(ctx, sac.KeyForNSScopedObj(alert.GetDeployment())...); err != nil || !ok {
		return nil, false, err
	}
	return alert, true, nil
}

// CountAlerts returns the number of alerts that are active
func (ds *datastoreImpl) CountAlerts(ctx context.Context) (int, error) {
	activeQuery := searchCommon.NewQueryBuilder().AddStrings(searchCommon.ViolationState, storage.ViolationState_ACTIVE.String()).ProtoQuery()
	results, err := ds.Search(ctx, activeQuery)
	if err != nil {
		return 0, err
	}
	return len(results), nil
}

// AddAlert inserts an alert into storage and into the indexer
func (ds *datastoreImpl) AddAlert(ctx context.Context, alert *storage.Alert) error {
	if ok, err := alertSAC.WriteAllowed(ctx, sac.KeyForNSScopedObj(alert.GetDeployment())...); err != nil || !ok {
		return errors.New("permission denied")
	}

	ds.keyedMutex.Lock(alert.GetId())
	defer ds.keyedMutex.Unlock(alert.GetId())

	if err := ds.storage.AddAlert(alert); err != nil {
		return err
	}
	return ds.indexer.AddListAlert(convert.AlertToListAlert(alert))
}

// UpdateAlert updates an alert in storage and in the indexer
func (ds *datastoreImpl) UpdateAlert(ctx context.Context, alert *storage.Alert) error {
	if ok, err := alertSAC.WriteAllowed(ctx, sac.KeyForNSScopedObj(alert.GetDeployment())...); err != nil || !ok {
		return errors.New("permission denied")
	}

	ds.keyedMutex.Lock(alert.GetId())
	defer ds.keyedMutex.Unlock(alert.GetId())

	oldAlert, exists, err := ds.GetAlert(ctx, alert.GetId())
	if err != nil {
		log.Infof("Error in get alert: %v", err)
		return err
	}
	if exists {
		if !hasSameScope(alert.GetDeployment(), oldAlert.GetDeployment()) {
			return errors.New("cannot change the cluster or namespace of an existing alert")
		}
	}

	return ds.updateAlertNoLock(alert)
}

func (ds *datastoreImpl) MarkAlertStale(ctx context.Context, id string) error {
	ds.keyedMutex.Lock(id)
	defer ds.keyedMutex.Unlock(id)

	alert, exists, err := ds.GetAlert(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("alert with id '%s' does not exist", id)
	}

	if ok, err := alertSAC.WriteAllowed(ctx, sac.KeyForNSScopedObj(alert.GetDeployment())...); err != nil || !ok {
		return errors.New("permission denied")
	}
	alert.State = storage.ViolationState_RESOLVED
	return ds.updateAlertNoLock(alert)
}

func (ds *datastoreImpl) DeleteAlerts(ctx context.Context, ids ...string) error {
	if ok, err := alertSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}

	errorList := errorhelpers.NewErrorList("deleting alert")
	if err := ds.storage.DeleteAlerts(ids...); err != nil {
		errorList.AddError(err)
	}
	if err := ds.indexer.DeleteListAlerts(ids); err != nil {
		errorList.AddError(err)
	}
	return errorList.ToError()
}

func (ds *datastoreImpl) updateAlertNoLock(alert *storage.Alert) error {
	// Checks pass then update.
	if err := ds.storage.UpdateAlert(alert); err != nil {
		return err
	}
	return ds.indexer.AddListAlert(convert.AlertToListAlert(alert))
}

func hasSameScope(o1, o2 sac.NamespaceScopedObject) bool {
	return o1.GetClusterId() == o2.GetClusterId() && o1.GetNamespace() == o2.GetNamespace()
}

func (ds *datastoreImpl) buildIndex() error {
	defer debug.FreeOSMemory()

	log.Infof("[STARTUP] Determining if alert db/indexer reconciliation is needed")
	indexerTxNum := ds.indexer.GetTxnCount()

	dbTxNum, err := ds.storage.GetTxnCount()
	if err != nil {
		return err
	}

	if !txn.ReconciliationNeeded(dbTxNum, indexerTxNum) {
		log.Infof("[STARTUP] Reconciliation for alerts is not needed")
		return nil
	}

	log.Info("[STARTUP] Reconciling DB and Indexer by indexing alerts")

	if err := ds.indexer.ResetIndex(); err != nil {
		return err
	}

	defer debug.FreeOSMemory()

	alertIDs, err := ds.storage.GetAlertIDs()
	if err != nil {
		return err
	}

	alertBatcher := batcher.New(len(alertIDs), alertBatchSize)
	for start, end, valid := alertBatcher.Next(); valid; start, end, valid = alertBatcher.Next() {
		listAlerts, _, err := ds.storage.GetListAlerts(alertIDs[start:end])
		if err != nil {
			return err
		}
		if err := ds.indexer.AddListAlerts(listAlerts); err != nil {
			return err
		}
	}

	log.Info("[STARTUP] Successfully indexed all alerts")

	if err := ds.storage.IncTxnCount(); err != nil {
		return err
	}
	if err := ds.indexer.SetTxnCount(dbTxNum + 1); err != nil {
		return err
	}

	return nil
}
