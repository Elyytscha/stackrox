package datastore

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/alert/datastore/internal/commentsstore"
	"github.com/stackrox/rox/central/alert/datastore/internal/index"
	"github.com/stackrox/rox/central/alert/datastore/internal/search"
	"github.com/stackrox/rox/central/alert/datastore/internal/store"
	"github.com/stackrox/rox/central/analystnotes"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/alert/convert"
	"github.com/stackrox/rox/pkg/batcher"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/debug"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sac"
	searchCommon "github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/paginated"
	"github.com/stackrox/rox/pkg/sliceutils"
	"github.com/stackrox/rox/pkg/sync"
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
	storage         store.Store
	commentsStorage commentsstore.Store
	indexer         index.Indexer
	searcher        search.Searcher
	keyedMutex      *concurrency.KeyedMutex
}

func (ds *datastoreImpl) Search(ctx context.Context, q *v1.Query) ([]searchCommon.Result, error) {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "Search")

	return ds.searcher.Search(ctx, q)
}

func (ds *datastoreImpl) SearchListAlerts(ctx context.Context, q *v1.Query) ([]*storage.ListAlert, error) {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "SearchListAlerts")

	return ds.searcher.SearchListAlerts(ctx, q)
}

// SearchAlerts returns search results for the given request.
func (ds *datastoreImpl) SearchAlerts(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "SearchAlerts")

	return ds.searcher.SearchAlerts(ctx, q)
}

// SearchRawAlerts returns search results for the given request in the form of a slice of alerts.
func (ds *datastoreImpl) SearchRawAlerts(ctx context.Context, q *v1.Query) ([]*storage.Alert, error) {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "SearchRawAlerts")

	return ds.searcher.SearchRawAlerts(ctx, q)
}

func (ds *datastoreImpl) ListAlerts(ctx context.Context, request *v1.ListAlertsRequest) ([]*storage.ListAlert, error) {
	var q *v1.Query
	if request.GetQuery() == "" {
		q = searchCommon.EmptyQuery()
	} else {
		var err error
		q, err = searchCommon.ParseQuery(request.GetQuery())
		if err != nil {
			return nil, err
		}
	}

	paginated.FillPagination(q, request.GetPagination(), math.MaxInt32)

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
func (ds *datastoreImpl) UpsertAlert(ctx context.Context, alert *storage.Alert) error {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "UpsertAlert")

	if ok, err := alertSAC.WriteAllowed(ctx, sac.KeyForNSScopedObj(alert.GetDeployment())...); err != nil || !ok {
		return errors.New("permission denied")
	}

	ds.keyedMutex.Lock(alert.GetId())
	defer ds.keyedMutex.Unlock(alert.GetId())

	if err := ds.storage.UpsertAlert(alert); err != nil {
		return err
	}
	if err := ds.indexer.AddListAlert(convert.AlertToListAlert(alert)); err != nil {
		return err
	}
	return ds.storage.AckKeysIndexed(alert.GetId())
}

// UpdateAlert updates an alert in storage and in the indexer
func (ds *datastoreImpl) UpdateAlertBatch(ctx context.Context, alert *storage.Alert, waitGroup *sync.WaitGroup, c chan error) {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "UpdateAlertBatch")

	defer waitGroup.Done()
	ds.keyedMutex.Lock(alert.GetId())
	defer ds.keyedMutex.Unlock(alert.GetId())
	oldAlert, exists, err := ds.GetAlert(ctx, alert.GetId())
	if err != nil {
		log.Errorf("error in get alert: %+v", err)
		c <- err
		return
	}
	if exists {
		if !hasSameScope(alert.GetDeployment(), oldAlert.GetDeployment()) {
			c <- fmt.Errorf("cannot change the cluster or namespace of an existing alert %q", alert.GetId())
			return
		}
	}
	err = ds.updateAlertNoLock(alert)
	if err != nil {
		c <- err
	}
}

// UpdateAlert updates an alert in storage and in the indexer
func (ds *datastoreImpl) UpsertAlerts(ctx context.Context, alertBatch []*storage.Alert) error {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "UpsertAlerts")

	var waitGroup sync.WaitGroup
	c := make(chan error, len(alertBatch))
	for _, alert := range alertBatch {
		waitGroup.Add(1)
		go ds.UpdateAlertBatch(ctx, alert, &waitGroup, c)
	}
	waitGroup.Wait()
	close(c)
	if len(c) > 0 {
		errorList := errorhelpers.NewErrorList(fmt.Sprintf("found %d errors while resolving alerts", len(c)))
		for err := range c {
			errorList.AddError(err)
		}
		return errorList.ToError()
	}
	return nil
}

func (ds *datastoreImpl) MarkAlertStale(ctx context.Context, id string) error {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "MarkAlertStale")

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
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "DeleteAlerts")

	if ok, err := alertSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}

	errorList := errorhelpers.NewErrorList("deleting alert")
	if err := ds.storage.DeleteAlerts(ids...); err != nil {
		errorList.AddError(err)
	} else {
		for _, id := range ids {
			err := ds.commentsStorage.RemoveAlertComments(id)
			if err != nil {
				err = errors.Wrapf(err, "deleting comments of alert %q", id)
				errorList.AddError(err)
			}
		}
	}
	if err := ds.indexer.DeleteListAlerts(ids); err != nil {
		errorList.AddError(err)
	}
	if err := ds.storage.AckKeysIndexed(ids...); err != nil {
		errorList.AddError(err)
	}

	return errorList.ToError()
}

func (ds *datastoreImpl) GetAlertComments(ctx context.Context, alertID string) (comments []*storage.Comment, err error) {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "GetAlertComments")
	_, exists, err := ds.GetAlert(ctx, alertID)
	if err != nil || !exists {
		return nil, err
	}

	return ds.commentsStorage.GetCommentsForAlert(alertID)
}

func (ds *datastoreImpl) AddAlertComment(ctx context.Context, request *storage.Comment) (string, error) {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "AddAlertComment")
	alert, exists, err := ds.GetAlert(ctx, request.GetResourceId())
	if err != nil || !exists {
		return "", err
	}
	if ok, err := alertSAC.WriteAllowed(ctx, sac.KeyForNSScopedObj(alert.GetDeployment())...); err != nil || !ok {
		return "", errors.New("permission denied")
	}

	request.User = analystnotes.UserFromContext(ctx)
	return ds.commentsStorage.AddAlertComment(request)
}

func (ds *datastoreImpl) UpdateAlertComment(ctx context.Context, request *storage.Comment) error {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "UpdateAlertComment")
	alert, exists, err := ds.GetAlert(ctx, request.GetResourceId())
	if err != nil || !exists {
		return err
	}
	if ok, err := alertSAC.WriteAllowed(ctx, sac.KeyForNSScopedObj(alert.GetDeployment())...); err != nil || !ok {
		return errors.New("permission denied")
	}

	user := analystnotes.UserFromContext(ctx)
	existingComment, err := ds.commentsStorage.GetComment(request.GetResourceId(), request.GetCommentId())
	if err != nil {
		return errors.Wrap(err, "failed to get the alert comment")
	}
	if existingComment == nil {
		return errors.Errorf("cannot update comment %q for alert %q: it does not exist", request.GetCommentId(), request.GetResourceId())
	}
	if !analystnotes.CommentIsModifiableUser(user, existingComment) {
		return errors.New("user cannot modify comment: permission denied")
	}
	request.User = user
	return ds.commentsStorage.UpdateAlertComment(request)
}

func (ds *datastoreImpl) RemoveAlertComment(ctx context.Context, alertID, commentID string) error {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "RemoveAlertComment")
	alert, exists, err := ds.GetAlert(ctx, alertID)
	if err != nil || !exists {
		return err
	}
	if ok, err := alertSAC.WriteAllowed(ctx, sac.KeyForNSScopedObj(alert.GetDeployment())...); err != nil || !ok {
		return errors.New("permission denied")
	}

	existingComment, err := ds.commentsStorage.GetComment(alertID, commentID)
	if err != nil {
		return errors.Wrap(err, "failed to get the alert comment")
	}
	if existingComment == nil {
		// Comment has already been deleted, all good
		return nil
	}
	if !analystnotes.CommentIsDeletable(ctx, existingComment) {
		return errors.New("user cannot delete comment: permission denied")
	}

	return ds.commentsStorage.RemoveAlertComment(alertID, commentID)
}

func (ds *datastoreImpl) AddAlertTags(ctx context.Context, resourceID string, tags []string) ([]string, error) {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "AddAlertTags")

	alert, exists, err := ds.storage.GetAlert(resourceID)
	if err != nil {
		return nil, errors.Wrapf(err, "error fetching alert %q from the DB", resourceID)
	}
	if !exists {
		return nil, fmt.Errorf("cannot add tags to alert %q that no longer exists", resourceID)
	}
	if ok, err := alertSAC.WriteAllowed(ctx, sac.KeyForNSScopedObj(alert.GetDeployment())...); err != nil || !ok {
		return nil, errors.New("permission denied")
	}

	allTags := sliceutils.StringUnion(alert.GetTags(), tags)
	sort.Strings(allTags)
	alert.Tags = allTags
	if err := ds.updateAlertNoLock(alert); err != nil {
		return nil, errors.Wrapf(err, "error upserting alert %q", alert.GetId())
	}

	return alert.GetTags(), nil
}

func (ds *datastoreImpl) RemoveAlertTags(ctx context.Context, resourceID string, tags []string) error {
	defer metrics.SetDatastoreFunctionDuration(time.Now(), "Alert", "DeleteAlertTags")

	alert, exists, err := ds.storage.GetAlert(resourceID)
	if err != nil {
		return errors.Wrapf(err, "error fetching alert %q from the DB", resourceID)
	}
	if !exists {
		return fmt.Errorf("cannot add tags to alert %q that no longer exists", resourceID)
	}
	if ok, err := alertSAC.WriteAllowed(ctx, sac.KeyForNSScopedObj(alert.GetDeployment())...); err != nil || !ok {
		return errors.New("permission denied")
	}

	remainingTags := sliceutils.StringDifference(alert.GetTags(), tags)
	sort.Strings(remainingTags)

	if len(remainingTags) == 0 {
		alert.Tags = nil
	} else {
		alert.Tags = remainingTags
	}
	if err := ds.updateAlertNoLock(alert); err != nil {
		return fmt.Errorf("error upserting alert %q", alert.GetId())
	}

	return nil
}

func (ds *datastoreImpl) updateAlertNoLock(alert *storage.Alert) error {
	// Checks pass then update.
	if err := ds.storage.UpsertAlert(alert); err != nil {
		return err
	}
	if err := ds.indexer.AddListAlert(convert.AlertToListAlert(alert)); err != nil {
		return err
	}
	return ds.storage.AckKeysIndexed(alert.GetId())
}

func hasSameScope(o1, o2 sac.NamespaceScopedObject) bool {
	return o1.GetClusterId() == o2.GetClusterId() && o1.GetNamespace() == o2.GetNamespace()
}

func (ds *datastoreImpl) fullReindex() error {
	log.Info("[STARTUP] Reindexing all alerts")

	alertIDs, err := ds.storage.GetAlertIDs()
	if err != nil {
		return err
	}
	log.Infof("[STARTUP] Found %d alerts to index", len(alertIDs))
	alertBatcher := batcher.New(len(alertIDs), alertBatchSize)
	for start, end, valid := alertBatcher.Next(); valid; start, end, valid = alertBatcher.Next() {
		listAlerts, _, err := ds.storage.GetListAlerts(alertIDs[start:end])
		if err != nil {
			return err
		}
		if err := ds.indexer.AddListAlerts(listAlerts); err != nil {
			return err
		}
		if end%(alertBatchSize*10) == 0 {
			log.Infof("[STARTUP] Successfully indexed %d/%d alerts", end, len(alertIDs))
		}
	}
	log.Infof("[STARTUP] Successfully indexed %d alerts", len(alertIDs))

	// Clear the keys because we just re-indexed everything
	keys, err := ds.storage.GetKeysToIndex()
	if err != nil {
		return err
	}
	if err := ds.storage.AckKeysIndexed(keys...); err != nil {
		return err
	}

	// Write out that initial indexing is complete
	if err := ds.indexer.MarkInitialIndexingComplete(); err != nil {
		return err
	}

	return nil
}

func (ds *datastoreImpl) buildIndex() error {
	defer debug.FreeOSMemory()

	needsFullIndexing, err := ds.indexer.NeedsInitialIndexing()
	if err != nil {
		return err
	}
	if needsFullIndexing {
		return ds.fullReindex()
	}

	log.Info("[STARTUP] Determining if alert db/indexer reconciliation is needed")
	keysToIndex, err := ds.storage.GetKeysToIndex()
	if err != nil {
		return errors.Wrap(err, "error retrieving keys to index from store")
	}

	log.Infof("[STARTUP] Found %d Alerts to index", len(keysToIndex))

	defer debug.FreeOSMemory()

	alertBatcher := batcher.New(len(keysToIndex), alertBatchSize)
	for start, end, valid := alertBatcher.Next(); valid; start, end, valid = alertBatcher.Next() {
		listAlerts, missingIndices, err := ds.storage.GetListAlerts(keysToIndex[start:end])
		if err != nil {
			return err
		}
		if err := ds.indexer.AddListAlerts(listAlerts); err != nil {
			return err
		}

		if len(missingIndices) > 0 {
			idsToRemove := make([]string, 0, len(missingIndices))
			for _, missingIdx := range missingIndices {
				idsToRemove = append(idsToRemove, keysToIndex[start:end][missingIdx])
			}
			if err := ds.indexer.DeleteListAlerts(idsToRemove); err != nil {
				return err
			}
		}
		// Ack keys so that even if central restarts, we don't need to reindex them again
		if err := ds.storage.AckKeysIndexed(keysToIndex[start:end]...); err != nil {
			return err
		}

		log.Infof("[STARTUP] Successfully indexed %d/%d alerts", end, len(keysToIndex))
	}

	log.Info("[STARTUP] Successfully indexed all out of sync alerts")
	return nil
}

func (ds *datastoreImpl) WalkAll(ctx context.Context, fn func(*storage.ListAlert) error) error {
	if ok, err := alertSAC.ReadAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return errors.New("permission denied")
	}

	return ds.storage.WalkAll(fn)
}
