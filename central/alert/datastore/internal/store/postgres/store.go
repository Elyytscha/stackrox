// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/sac"
	pgSearch "github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/sync"
	"gorm.io/gorm"
)

const (
	baseTable = "alerts"
	storeName = "Alert"

	batchAfter = 100

	// using copyFrom, we may not even want to batch.  It would probably be simpler
	// to deal with failures if we just sent it all.  Something to think about as we
	// proceed and move into more e2e and larger performance testing
	batchSize = 10000
)

var (
	log            = logging.LoggerForModule()
	schema         = pkgSchema.AlertsSchema
	targetResource = resources.Alert
)

// Store is the interface to interact with the storage for storage.Alert
type Store interface {
	Upsert(ctx context.Context, obj *storage.Alert) error
	UpsertMany(ctx context.Context, objs []*storage.Alert) error
	Delete(ctx context.Context, id string) error
	DeleteByQuery(ctx context.Context, q *v1.Query) error
	DeleteMany(ctx context.Context, identifiers []string) error

	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)

	Get(ctx context.Context, id string) (*storage.Alert, bool, error)
	GetByQuery(ctx context.Context, query *v1.Query) ([]*storage.Alert, error)
	GetMany(ctx context.Context, identifiers []string) ([]*storage.Alert, []int, error)
	GetIDs(ctx context.Context) ([]string, error)

	Walk(ctx context.Context, fn func(obj *storage.Alert) error) error
}

type storeImpl struct {
	*pgSearch.GenericStore[storage.Alert, *storage.Alert]
	db    postgres.DB
	mutex sync.RWMutex
}

// New returns a new Store instance using the provided sql instance.
func New(db postgres.DB) Store {
	return &storeImpl{
		db: db,
		GenericStore: pgSearch.NewGenericStore[storage.Alert, *storage.Alert](
			db,
			schema,
			pkGetter,
			metricsSetAcquireDBConnDuration,
			metricsSetPostgresOperationDurationTime,
			targetResource,
		),
	}
}

// region Helper functions

func pkGetter(obj *storage.Alert) string {
	return obj.GetId()
}

func metricsSetPostgresOperationDurationTime(start time.Time, op ops.Op) {
	metrics.SetPostgresOperationDurationTime(start, op, storeName)
}

func metricsSetAcquireDBConnDuration(start time.Time, op ops.Op) {
	metrics.SetAcquireDBConnDuration(start, op, storeName)
}

func insertIntoAlerts(_ context.Context, batch *pgx.Batch, obj *storage.Alert) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(obj.GetId()),
		obj.GetPolicy().GetId(),
		obj.GetPolicy().GetName(),
		obj.GetPolicy().GetDescription(),
		obj.GetPolicy().GetDisabled(),
		obj.GetPolicy().GetCategories(),
		obj.GetPolicy().GetSeverity(),
		obj.GetPolicy().GetEnforcementActions(),
		pgutils.NilOrTime(obj.GetPolicy().GetLastUpdated()),
		obj.GetPolicy().GetSORTName(),
		obj.GetPolicy().GetSORTLifecycleStage(),
		obj.GetPolicy().GetSORTEnforcement(),
		obj.GetLifecycleStage(),
		pgutils.NilOrUUID(obj.GetClusterId()),
		obj.GetClusterName(),
		obj.GetNamespace(),
		pgutils.NilOrUUID(obj.GetNamespaceId()),
		pgutils.NilOrUUID(obj.GetDeployment().GetId()),
		obj.GetDeployment().GetName(),
		obj.GetDeployment().GetInactive(),
		obj.GetImage().GetId(),
		obj.GetImage().GetName().GetRegistry(),
		obj.GetImage().GetName().GetRemote(),
		obj.GetImage().GetName().GetTag(),
		obj.GetImage().GetName().GetFullName(),
		obj.GetResource().GetResourceType(),
		obj.GetResource().GetName(),
		obj.GetEnforcement().GetAction(),
		pgutils.NilOrTime(obj.GetTime()),
		obj.GetState(),
		serialized,
	}

	finalStr := "INSERT INTO alerts (Id, Policy_Id, Policy_Name, Policy_Description, Policy_Disabled, Policy_Categories, Policy_Severity, Policy_EnforcementActions, Policy_LastUpdated, Policy_SORTName, Policy_SORTLifecycleStage, Policy_SORTEnforcement, LifecycleStage, ClusterId, ClusterName, Namespace, NamespaceId, Deployment_Id, Deployment_Name, Deployment_Inactive, Image_Id, Image_Name_Registry, Image_Name_Remote, Image_Name_Tag, Image_Name_FullName, Resource_ResourceType, Resource_Name, Enforcement_Action, Time, State, serialized) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31) ON CONFLICT(Id) DO UPDATE SET Id = EXCLUDED.Id, Policy_Id = EXCLUDED.Policy_Id, Policy_Name = EXCLUDED.Policy_Name, Policy_Description = EXCLUDED.Policy_Description, Policy_Disabled = EXCLUDED.Policy_Disabled, Policy_Categories = EXCLUDED.Policy_Categories, Policy_Severity = EXCLUDED.Policy_Severity, Policy_EnforcementActions = EXCLUDED.Policy_EnforcementActions, Policy_LastUpdated = EXCLUDED.Policy_LastUpdated, Policy_SORTName = EXCLUDED.Policy_SORTName, Policy_SORTLifecycleStage = EXCLUDED.Policy_SORTLifecycleStage, Policy_SORTEnforcement = EXCLUDED.Policy_SORTEnforcement, LifecycleStage = EXCLUDED.LifecycleStage, ClusterId = EXCLUDED.ClusterId, ClusterName = EXCLUDED.ClusterName, Namespace = EXCLUDED.Namespace, NamespaceId = EXCLUDED.NamespaceId, Deployment_Id = EXCLUDED.Deployment_Id, Deployment_Name = EXCLUDED.Deployment_Name, Deployment_Inactive = EXCLUDED.Deployment_Inactive, Image_Id = EXCLUDED.Image_Id, Image_Name_Registry = EXCLUDED.Image_Name_Registry, Image_Name_Remote = EXCLUDED.Image_Name_Remote, Image_Name_Tag = EXCLUDED.Image_Name_Tag, Image_Name_FullName = EXCLUDED.Image_Name_FullName, Resource_ResourceType = EXCLUDED.Resource_ResourceType, Resource_Name = EXCLUDED.Resource_Name, Enforcement_Action = EXCLUDED.Enforcement_Action, Time = EXCLUDED.Time, State = EXCLUDED.State, serialized = EXCLUDED.serialized"
	batch.Queue(finalStr, values...)

	return nil
}

func (s *storeImpl) copyFromAlerts(ctx context.Context, tx *postgres.Tx, objs ...*storage.Alert) error {

	inputRows := [][]interface{}{}

	var err error

	// This is a copy so first we must delete the rows and re-add them
	// Which is essentially the desired behaviour of an upsert.
	var deletes []string

	copyCols := []string{

		"id",

		"policy_id",

		"policy_name",

		"policy_description",

		"policy_disabled",

		"policy_categories",

		"policy_severity",

		"policy_enforcementactions",

		"policy_lastupdated",

		"policy_sortname",

		"policy_sortlifecyclestage",

		"policy_sortenforcement",

		"lifecyclestage",

		"clusterid",

		"clustername",

		"namespace",

		"namespaceid",

		"deployment_id",

		"deployment_name",

		"deployment_inactive",

		"image_id",

		"image_name_registry",

		"image_name_remote",

		"image_name_tag",

		"image_name_fullname",

		"resource_resourcetype",

		"resource_name",

		"enforcement_action",

		"time",

		"state",

		"serialized",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		serialized, marshalErr := obj.Marshal()
		if marshalErr != nil {
			return marshalErr
		}

		inputRows = append(inputRows, []interface{}{

			pgutils.NilOrUUID(obj.GetId()),

			obj.GetPolicy().GetId(),

			obj.GetPolicy().GetName(),

			obj.GetPolicy().GetDescription(),

			obj.GetPolicy().GetDisabled(),

			obj.GetPolicy().GetCategories(),

			obj.GetPolicy().GetSeverity(),

			obj.GetPolicy().GetEnforcementActions(),

			pgutils.NilOrTime(obj.GetPolicy().GetLastUpdated()),

			obj.GetPolicy().GetSORTName(),

			obj.GetPolicy().GetSORTLifecycleStage(),

			obj.GetPolicy().GetSORTEnforcement(),

			obj.GetLifecycleStage(),

			pgutils.NilOrUUID(obj.GetClusterId()),

			obj.GetClusterName(),

			obj.GetNamespace(),

			pgutils.NilOrUUID(obj.GetNamespaceId()),

			pgutils.NilOrUUID(obj.GetDeployment().GetId()),

			obj.GetDeployment().GetName(),

			obj.GetDeployment().GetInactive(),

			obj.GetImage().GetId(),

			obj.GetImage().GetName().GetRegistry(),

			obj.GetImage().GetName().GetRemote(),

			obj.GetImage().GetName().GetTag(),

			obj.GetImage().GetName().GetFullName(),

			obj.GetResource().GetResourceType(),

			obj.GetResource().GetName(),

			obj.GetEnforcement().GetAction(),

			pgutils.NilOrTime(obj.GetTime()),

			obj.GetState(),

			serialized,
		})

		// Add the ID to be deleted.
		deletes = append(deletes, obj.GetId())

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if err := s.DeleteMany(ctx, deletes); err != nil {
				return err
			}
			// clear the inserts and vals for the next batch
			deletes = nil

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"alerts"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return err
}

func (s *storeImpl) copyFrom(ctx context.Context, objs ...*storage.Alert) error {
	conn, err := s.AcquireConn(ctx, ops.Get)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	if err := s.copyFromAlerts(ctx, tx, objs...); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *storeImpl) upsert(ctx context.Context, objs ...*storage.Alert) error {
	conn, err := s.AcquireConn(ctx, ops.Get)
	if err != nil {
		return err
	}
	defer conn.Release()

	for _, obj := range objs {
		batch := &pgx.Batch{}
		if err := insertIntoAlerts(ctx, batch, obj); err != nil {
			return err
		}
		batchResults := conn.SendBatch(ctx, batch)
		var result *multierror.Error
		for i := 0; i < batch.Len(); i++ {
			_, err := batchResults.Exec()
			result = multierror.Append(result, err)
		}
		if err := batchResults.Close(); err != nil {
			return err
		}
		if err := result.ErrorOrNil(); err != nil {
			return err
		}
	}
	return nil
}

// endregion Helper functions

//// Interface functions

// Upsert saves the current state of an object in storage.
func (s *storeImpl) Upsert(ctx context.Context, obj *storage.Alert) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Upsert, "Alert")

	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_WRITE_ACCESS).Resource(targetResource).
		ClusterID(obj.GetClusterId()).Namespace(obj.GetNamespace())
	if !scopeChecker.IsAllowed() {
		return sac.ErrResourceAccessDenied
	}

	return pgutils.Retry(func() error {
		return s.upsert(ctx, obj)
	})
}

// UpsertMany saves the state of multiple objects in the storage.
func (s *storeImpl) UpsertMany(ctx context.Context, objs []*storage.Alert) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.UpdateMany, "Alert")

	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_WRITE_ACCESS).Resource(targetResource)
	if !scopeChecker.IsAllowed() {
		var deniedIDs []string
		for _, obj := range objs {
			subScopeChecker := scopeChecker.ClusterID(obj.GetClusterId()).Namespace(obj.GetNamespace())
			if !subScopeChecker.IsAllowed() {
				deniedIDs = append(deniedIDs, obj.GetId())
			}
		}
		if len(deniedIDs) != 0 {
			return errors.Wrapf(sac.ErrResourceAccessDenied, "modifying alerts with IDs [%s] was denied", strings.Join(deniedIDs, ", "))
		}
	}

	return pgutils.Retry(func() error {
		// Lock since copyFrom requires a delete first before being executed.  If multiple processes are updating
		// same subset of rows, both deletes could occur before the copyFrom resulting in unique constraint
		// violations
		if len(objs) < batchAfter {
			s.mutex.RLock()
			defer s.mutex.RUnlock()

			return s.upsert(ctx, objs...)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()

		return s.copyFrom(ctx, objs...)
	})
}

//// Interface functions - END

//// Used for testing

// CreateTableAndNewStore returns a new Store instance for testing.
func CreateTableAndNewStore(ctx context.Context, db postgres.DB, gormDB *gorm.DB) Store {
	pkgSchema.ApplySchemaForTable(ctx, gormDB, baseTable)
	return New(db)
}

// Destroy drops the tables associated with the target object type.
func Destroy(ctx context.Context, db postgres.DB) {
	dropTableAlerts(ctx, db)
}

func dropTableAlerts(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS alerts CASCADE")

}

//// Used for testing - END
