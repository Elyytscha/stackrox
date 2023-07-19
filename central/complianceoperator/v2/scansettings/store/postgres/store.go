// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	pgSearch "github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/sync"
	"gorm.io/gorm"
)

const (
	baseTable = "compliance_operator_scan_setting_v2"
	storeName = "ComplianceOperatorScanSettingV2"

	// using copyFrom, we may not even want to batch.  It would probably be simpler
	// to deal with failures if we just sent it all.  Something to think about as we
	// proceed and move into more e2e and larger performance testing
	batchSize = 10000
)

var (
	log            = logging.LoggerForModule()
	schema         = pkgSchema.ComplianceOperatorScanSettingV2Schema
	targetResource = resources.ComplianceOperator
)

type storeType = storage.ComplianceOperatorScanSettingV2

// Store is the interface to interact with the storage for storage.ComplianceOperatorScanSettingV2
type Store interface {
	Upsert(ctx context.Context, obj *storeType) error
	UpsertMany(ctx context.Context, objs []*storeType) error
	Delete(ctx context.Context, name string) error
	DeleteByQuery(ctx context.Context, q *v1.Query) error
	DeleteMany(ctx context.Context, identifiers []string) error

	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, name string) (bool, error)

	Get(ctx context.Context, name string) (*storeType, bool, error)
	GetByQuery(ctx context.Context, query *v1.Query) ([]*storeType, error)
	GetMany(ctx context.Context, identifiers []string) ([]*storeType, []int, error)
	GetIDs(ctx context.Context) ([]string, error)
	GetAll(ctx context.Context) ([]*storeType, error)

	Walk(ctx context.Context, fn func(obj *storeType) error) error
}

type storeImpl struct {
	*pgSearch.GenericStore[storeType, *storeType]
	mutex sync.RWMutex
}

// New returns a new Store instance using the provided sql instance.
func New(db postgres.DB) Store {
	return &storeImpl{
		GenericStore: pgSearch.NewGenericStore[storeType, *storeType](
			db,
			schema,
			pkGetter,
			insertIntoComplianceOperatorScanSettingV2,
			copyFromComplianceOperatorScanSettingV2,
			metricsSetAcquireDBConnDuration,
			metricsSetPostgresOperationDurationTime,
			pgSearch.GloballyScopedUpsertChecker[storeType, *storeType](targetResource),
			targetResource,
		),
	}
}

// region Helper functions

func pkGetter(obj *storeType) string {
	return obj.GetName()
}

func metricsSetPostgresOperationDurationTime(start time.Time, op ops.Op) {
	metrics.SetPostgresOperationDurationTime(start, op, storeName)
}

func metricsSetAcquireDBConnDuration(start time.Time, op ops.Op) {
	metrics.SetAcquireDBConnDuration(start, op, storeName)
}

func insertIntoComplianceOperatorScanSettingV2(ctx context.Context, batch *pgx.Batch, obj *storage.ComplianceOperatorScanSettingV2) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start
		obj.GetName(),
		obj.GetCreatedBy().GetName(),
		serialized,
	}

	finalStr := "INSERT INTO compliance_operator_scan_setting_v2 (Name, CreatedBy_Name, serialized) VALUES($1, $2, $3) ON CONFLICT(Name) DO UPDATE SET Name = EXCLUDED.Name, CreatedBy_Name = EXCLUDED.CreatedBy_Name, serialized = EXCLUDED.serialized"
	batch.Queue(finalStr, values...)

	var query string

	for childIndex, child := range obj.GetProfiles() {
		if err := insertIntoComplianceOperatorScanSettingV2Profiles(ctx, batch, child, obj.GetName(), childIndex); err != nil {
			return err
		}
	}

	query = "delete from compliance_operator_scan_setting_v2_profiles where compliance_operator_scan_setting_v2_Name = $1 AND idx >= $2"
	batch.Queue(query, obj.GetName(), len(obj.GetProfiles()))
	for childIndex, child := range obj.GetClusters() {
		if err := insertIntoComplianceOperatorScanSettingV2Clusters(ctx, batch, child, obj.GetName(), childIndex); err != nil {
			return err
		}
	}

	query = "delete from compliance_operator_scan_setting_v2_clusters where compliance_operator_scan_setting_v2_Name = $1 AND idx >= $2"
	batch.Queue(query, obj.GetName(), len(obj.GetClusters()))
	return nil
}

func insertIntoComplianceOperatorScanSettingV2Profiles(_ context.Context, batch *pgx.Batch, obj *storage.ComplianceOperatorScanSettingV2_Profile, complianceOperatorScanSettingV2Name string, idx int) error {

	values := []interface{}{
		// parent primary keys start
		complianceOperatorScanSettingV2Name,
		idx,
		obj.GetProfileId(),
	}

	finalStr := "INSERT INTO compliance_operator_scan_setting_v2_profiles (compliance_operator_scan_setting_v2_Name, idx, ProfileId) VALUES($1, $2, $3) ON CONFLICT(compliance_operator_scan_setting_v2_Name, idx) DO UPDATE SET compliance_operator_scan_setting_v2_Name = EXCLUDED.compliance_operator_scan_setting_v2_Name, idx = EXCLUDED.idx, ProfileId = EXCLUDED.ProfileId"
	batch.Queue(finalStr, values...)

	return nil
}

func insertIntoComplianceOperatorScanSettingV2Clusters(_ context.Context, batch *pgx.Batch, obj *storage.ComplianceOperatorScanSettingV2_ClusterScanStatus, complianceOperatorScanSettingV2Name string, idx int) error {

	values := []interface{}{
		// parent primary keys start
		complianceOperatorScanSettingV2Name,
		idx,
		pgutils.NilOrUUID(obj.GetClusterId()),
	}

	finalStr := "INSERT INTO compliance_operator_scan_setting_v2_clusters (compliance_operator_scan_setting_v2_Name, idx, ClusterId) VALUES($1, $2, $3) ON CONFLICT(compliance_operator_scan_setting_v2_Name, idx) DO UPDATE SET compliance_operator_scan_setting_v2_Name = EXCLUDED.compliance_operator_scan_setting_v2_Name, idx = EXCLUDED.idx, ClusterId = EXCLUDED.ClusterId"
	batch.Queue(finalStr, values...)

	return nil
}

func copyFromComplianceOperatorScanSettingV2(ctx context.Context, s pgSearch.Deleter, tx *postgres.Tx, objs ...*storage.ComplianceOperatorScanSettingV2) error {
	inputRows := make([][]interface{}, 0, batchSize)

	// This is a copy so first we must delete the rows and re-add them
	// Which is essentially the desired behaviour of an upsert.
	deletes := make([]string, 0, batchSize)

	copyCols := []string{
		"name",
		"createdby_name",
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
			obj.GetName(),
			obj.GetCreatedBy().GetName(),
			serialized,
		})

		// Add the ID to be deleted.
		deletes = append(deletes, obj.GetName())

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if err := s.DeleteMany(ctx, deletes); err != nil {
				return err
			}
			// clear the inserts and vals for the next batch
			deletes = deletes[:0]

			if _, err := tx.CopyFrom(ctx, pgx.Identifier{"compliance_operator_scan_setting_v2"}, copyCols, pgx.CopyFromRows(inputRows)); err != nil {
				return err
			}
			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	for idx, obj := range objs {
		_ = idx // idx may or may not be used depending on how nested we are, so avoid compile-time errors.

		if err := copyFromComplianceOperatorScanSettingV2Profiles(ctx, s, tx, obj.GetName(), obj.GetProfiles()...); err != nil {
			return err
		}
		if err := copyFromComplianceOperatorScanSettingV2Clusters(ctx, s, tx, obj.GetName(), obj.GetClusters()...); err != nil {
			return err
		}
	}

	return nil
}

func copyFromComplianceOperatorScanSettingV2Profiles(ctx context.Context, s pgSearch.Deleter, tx *postgres.Tx, complianceOperatorScanSettingV2Name string, objs ...*storage.ComplianceOperatorScanSettingV2_Profile) error {
	inputRows := make([][]interface{}, 0, batchSize)

	copyCols := []string{
		"compliance_operator_scan_setting_v2_name",
		"idx",
		"profileid",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{
			complianceOperatorScanSettingV2Name,
			idx,
			obj.GetProfileId(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if _, err := tx.CopyFrom(ctx, pgx.Identifier{"compliance_operator_scan_setting_v2_profiles"}, copyCols, pgx.CopyFromRows(inputRows)); err != nil {
				return err
			}
			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return nil
}

func copyFromComplianceOperatorScanSettingV2Clusters(ctx context.Context, s pgSearch.Deleter, tx *postgres.Tx, complianceOperatorScanSettingV2Name string, objs ...*storage.ComplianceOperatorScanSettingV2_ClusterScanStatus) error {
	inputRows := make([][]interface{}, 0, batchSize)

	copyCols := []string{
		"compliance_operator_scan_setting_v2_name",
		"idx",
		"clusterid",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{
			complianceOperatorScanSettingV2Name,
			idx,
			pgutils.NilOrUUID(obj.GetClusterId()),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if _, err := tx.CopyFrom(ctx, pgx.Identifier{"compliance_operator_scan_setting_v2_clusters"}, copyCols, pgx.CopyFromRows(inputRows)); err != nil {
				return err
			}
			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return nil
}

// endregion Helper functions

// region Used for testing

// CreateTableAndNewStore returns a new Store instance for testing.
func CreateTableAndNewStore(ctx context.Context, db postgres.DB, gormDB *gorm.DB) Store {
	pkgSchema.ApplySchemaForTable(ctx, gormDB, baseTable)
	return New(db)
}

// Destroy drops the tables associated with the target object type.
func Destroy(ctx context.Context, db postgres.DB) {
	dropTableComplianceOperatorScanSettingV2(ctx, db)
}

func dropTableComplianceOperatorScanSettingV2(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS compliance_operator_scan_setting_v2 CASCADE")
	dropTableComplianceOperatorScanSettingV2Profiles(ctx, db)
	dropTableComplianceOperatorScanSettingV2Clusters(ctx, db)

}

func dropTableComplianceOperatorScanSettingV2Profiles(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS compliance_operator_scan_setting_v2_profiles CASCADE")

}

func dropTableComplianceOperatorScanSettingV2Clusters(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS compliance_operator_scan_setting_v2_clusters CASCADE")

}

// endregion Used for testing
