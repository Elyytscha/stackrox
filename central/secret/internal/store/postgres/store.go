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
	baseTable = "secrets"

	batchAfter = 100

	// using copyFrom, we may not even want to batch.  It would probably be simpler
	// to deal with failures if we just sent it all.  Something to think about as we
	// proceed and move into more e2e and larger performance testing
	batchSize = 10000
)

var (
	log            = logging.LoggerForModule()
	schema         = pkgSchema.SecretsSchema
	targetResource = resources.Secret
)

// Store is the interface to interact with the storage for storage.Secret
type Store interface {
	Upsert(ctx context.Context, obj *storage.Secret) error
	UpsertMany(ctx context.Context, objs []*storage.Secret) error
	Delete(ctx context.Context, id string) error
	DeleteByQuery(ctx context.Context, q *v1.Query) error
	DeleteMany(ctx context.Context, identifiers []string) error

	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)

	Get(ctx context.Context, id string) (*storage.Secret, bool, error)
	GetByQuery(ctx context.Context, query *v1.Query) ([]*storage.Secret, error)
	GetMany(ctx context.Context, identifiers []string) ([]*storage.Secret, []int, error)
	GetIDs(ctx context.Context) ([]string, error)

	Walk(ctx context.Context, fn func(obj *storage.Secret) error) error
}

type storeImpl struct {
	*pgSearch.GenericSingleIDStore[storage.Secret, *storage.Secret]
	db    postgres.DB
	mutex sync.RWMutex
}

// New returns a new Store instance using the provided sql instance.
func New(db postgres.DB) Store {
	return &storeImpl{
		GenericSingleIDStore: pgSearch.NewGenericSingleIDStore[storage.Secret, *storage.Secret](
			db,
			"Secret",
			targetResource,
			schema,
			metrics.SetPostgresOperationDurationTime,
			pkGetter,
		),
		db: db,
	}
}

//// Helper functions

func pkGetter(obj *storage.Secret) string {
	return obj.GetId()
}

func insertIntoSecrets(ctx context.Context, batch *pgx.Batch, obj *storage.Secret) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(obj.GetId()),
		obj.GetName(),
		pgutils.NilOrUUID(obj.GetClusterId()),
		obj.GetClusterName(),
		obj.GetNamespace(),
		pgutils.NilOrTime(obj.GetCreatedAt()),
		serialized,
	}

	finalStr := "INSERT INTO secrets (Id, Name, ClusterId, ClusterName, Namespace, CreatedAt, serialized) VALUES($1, $2, $3, $4, $5, $6, $7) ON CONFLICT(Id) DO UPDATE SET Id = EXCLUDED.Id, Name = EXCLUDED.Name, ClusterId = EXCLUDED.ClusterId, ClusterName = EXCLUDED.ClusterName, Namespace = EXCLUDED.Namespace, CreatedAt = EXCLUDED.CreatedAt, serialized = EXCLUDED.serialized"
	batch.Queue(finalStr, values...)

	var query string

	for childIndex, child := range obj.GetFiles() {
		if err := insertIntoSecretsFiles(ctx, batch, child, obj.GetId(), childIndex); err != nil {
			return err
		}
	}

	query = "delete from secrets_files where secrets_Id = $1 AND idx >= $2"
	batch.Queue(query, pgutils.NilOrUUID(obj.GetId()), len(obj.GetFiles()))
	return nil
}

func insertIntoSecretsFiles(ctx context.Context, batch *pgx.Batch, obj *storage.SecretDataFile, secretID string, idx int) error {

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(secretID),
		idx,
		obj.GetType(),
		pgutils.NilOrTime(obj.GetCert().GetEndDate()),
	}

	finalStr := "INSERT INTO secrets_files (secrets_Id, idx, Type, Cert_EndDate) VALUES($1, $2, $3, $4) ON CONFLICT(secrets_Id, idx) DO UPDATE SET secrets_Id = EXCLUDED.secrets_Id, idx = EXCLUDED.idx, Type = EXCLUDED.Type, Cert_EndDate = EXCLUDED.Cert_EndDate"
	batch.Queue(finalStr, values...)

	var query string

	for childIndex, child := range obj.GetImagePullSecret().GetRegistries() {
		if err := insertIntoSecretsFilesRegistries(ctx, batch, child, secretID, idx, childIndex); err != nil {
			return err
		}
	}

	query = "delete from secrets_files_registries where secrets_Id = $1 AND secrets_files_idx = $2 AND idx >= $3"
	batch.Queue(query, pgutils.NilOrUUID(secretID), idx, len(obj.GetImagePullSecret().GetRegistries()))
	return nil
}

func insertIntoSecretsFilesRegistries(_ context.Context, batch *pgx.Batch, obj *storage.ImagePullSecret_Registry, secretID string, secretFileIdx int, idx int) error {

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(secretID),
		secretFileIdx,
		idx,
		obj.GetName(),
	}

	finalStr := "INSERT INTO secrets_files_registries (secrets_Id, secrets_files_idx, idx, Name) VALUES($1, $2, $3, $4) ON CONFLICT(secrets_Id, secrets_files_idx, idx) DO UPDATE SET secrets_Id = EXCLUDED.secrets_Id, secrets_files_idx = EXCLUDED.secrets_files_idx, idx = EXCLUDED.idx, Name = EXCLUDED.Name"
	batch.Queue(finalStr, values...)

	return nil
}

func (s *storeImpl) copyFromSecrets(ctx context.Context, tx *postgres.Tx, objs ...*storage.Secret) error {

	inputRows := [][]interface{}{}

	var err error

	// This is a copy so first we must delete the rows and re-add them
	// Which is essentially the desired behaviour of an upsert.
	var deletes []string

	copyCols := []string{

		"id",

		"name",

		"clusterid",

		"clustername",

		"namespace",

		"createdat",

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

			obj.GetName(),

			pgutils.NilOrUUID(obj.GetClusterId()),

			obj.GetClusterName(),

			obj.GetNamespace(),

			pgutils.NilOrTime(obj.GetCreatedAt()),

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

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"secrets"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	for idx, obj := range objs {
		_ = idx // idx may or may not be used depending on how nested we are, so avoid compile-time errors.

		if err = s.copyFromSecretsFiles(ctx, tx, obj.GetId(), obj.GetFiles()...); err != nil {
			return err
		}
	}

	return err
}

func (s *storeImpl) copyFromSecretsFiles(ctx context.Context, tx *postgres.Tx, secretID string, objs ...*storage.SecretDataFile) error {

	inputRows := [][]interface{}{}

	var err error

	copyCols := []string{

		"secrets_id",

		"idx",

		"type",

		"cert_enddate",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{

			pgutils.NilOrUUID(secretID),

			idx,

			obj.GetType(),

			pgutils.NilOrTime(obj.GetCert().GetEndDate()),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"secrets_files"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	for idx, obj := range objs {
		_ = idx // idx may or may not be used depending on how nested we are, so avoid compile-time errors.

		if err = s.copyFromSecretsFilesRegistries(ctx, tx, secretID, idx, obj.GetImagePullSecret().GetRegistries()...); err != nil {
			return err
		}
	}

	return err
}

func (s *storeImpl) copyFromSecretsFilesRegistries(ctx context.Context, tx *postgres.Tx, secretID string, secretFileIdx int, objs ...*storage.ImagePullSecret_Registry) error {

	inputRows := [][]interface{}{}

	var err error

	copyCols := []string{

		"secrets_id",

		"secrets_files_idx",

		"idx",

		"name",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		inputRows = append(inputRows, []interface{}{

			pgutils.NilOrUUID(secretID),

			secretFileIdx,

			idx,

			obj.GetName(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"secrets_files_registries"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return err
}

func (s *storeImpl) acquireConn(ctx context.Context, op ops.Op, typ string) (*postgres.Conn, func(), error) {
	defer metrics.SetAcquireDBConnDuration(time.Now(), op, typ)
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	return conn, conn.Release, nil
}

func (s *storeImpl) copyFrom(ctx context.Context, objs ...*storage.Secret) error {
	conn, release, err := s.acquireConn(ctx, ops.Get, "Secret")
	if err != nil {
		return err
	}
	defer release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	if err := s.copyFromSecrets(ctx, tx, objs...); err != nil {
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

func (s *storeImpl) upsert(ctx context.Context, objs ...*storage.Secret) error {
	conn, release, err := s.acquireConn(ctx, ops.Get, "Secret")
	if err != nil {
		return err
	}
	defer release()

	for _, obj := range objs {
		batch := &pgx.Batch{}
		if err := insertIntoSecrets(ctx, batch, obj); err != nil {
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

//// Helper functions - END

//// Interface functions

// Upsert saves the current state of an object in storage.
func (s *storeImpl) Upsert(ctx context.Context, obj *storage.Secret) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Upsert, "Secret")

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
func (s *storeImpl) UpsertMany(ctx context.Context, objs []*storage.Secret) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.UpdateMany, "Secret")

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
			return errors.Wrapf(sac.ErrResourceAccessDenied, "modifying secrets with IDs [%s] was denied", strings.Join(deniedIDs, ", "))
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

//// Stubs for satisfying legacy interfaces

//// Interface functions - END

//// Used for testing

// CreateTableAndNewStore returns a new Store instance for testing.
func CreateTableAndNewStore(ctx context.Context, db postgres.DB, gormDB *gorm.DB) Store {
	pkgSchema.ApplySchemaForTable(ctx, gormDB, baseTable)
	return New(db)
}

// Destroy drops the tables associated with the target object type.
func Destroy(ctx context.Context, db postgres.DB) {
	dropTableSecrets(ctx, db)
}

func dropTableSecrets(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS secrets CASCADE")
	dropTableSecretsFiles(ctx, db)

}

func dropTableSecretsFiles(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS secrets_files CASCADE")
	dropTableSecretsFilesRegistries(ctx, db)

}

func dropTableSecretsFilesRegistries(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS secrets_files_registries CASCADE")

}

//// Used for testing - END
