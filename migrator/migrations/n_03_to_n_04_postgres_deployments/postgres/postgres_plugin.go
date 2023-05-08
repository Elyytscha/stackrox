// Code originally generated by pg-bindings generator.

package postgres

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	pkgSchema "github.com/stackrox/rox/migrator/migrations/frozenschema/v73"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"github.com/stackrox/rox/pkg/search"
	pgSearch "github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

const (
	baseTable = "deployments"

	batchAfter = 100

	// using copyFrom, we may not even want to batch.  It would probably be simpler
	// to deal with failures if we just sent it all.  Something to think about as we
	// proceed and move into more e2e and larger performance testing
	batchSize = 10000

	cursorBatchSize = 50
	deleteBatchSize = 5000
)

var (
	log    = logging.LoggerForModule()
	schema = pkgSchema.DeploymentsSchema
)

// Store for migration
type Store interface {
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)
	Get(ctx context.Context, id string) (*storage.Deployment, bool, error)
	GetByQuery(ctx context.Context, query *v1.Query) ([]*storage.Deployment, error)
	Upsert(ctx context.Context, obj *storage.Deployment) error
	UpsertMany(ctx context.Context, objs []*storage.Deployment) error
	Delete(ctx context.Context, id string) error
	DeleteByQuery(ctx context.Context, q *v1.Query) error
	GetIDs(ctx context.Context) ([]string, error)
	GetMany(ctx context.Context, ids []string) ([]*storage.Deployment, []int, error)
	DeleteMany(ctx context.Context, ids []string) error

	Walk(ctx context.Context, fn func(obj *storage.Deployment) error) error
}

type storeImpl struct {
	db    postgres.DB
	mutex sync.Mutex
}

// New returns a new Store instance using the provided sql instance.
func New(db postgres.DB) Store {
	return &storeImpl{
		db: db,
	}
}

func insertIntoDeployments(ctx context.Context, batch *pgx.Batch, obj *storage.Deployment) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	if pgutils.NilOrUUID(obj.GetId()) == nil {
		utils.Should(errors.Errorf("Id is not a valid uuid -- %q", obj.GetId()))
		return nil
	}

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(obj.GetId()),
		obj.GetName(),
		obj.GetType(),
		obj.GetNamespace(),
		pgutils.NilOrUUID(obj.GetNamespaceId()),
		obj.GetOrchestratorComponent(),
		obj.GetLabels(),
		obj.GetPodLabels(),
		pgutils.NilOrTime(obj.GetCreated()),
		pgutils.NilOrUUID(obj.GetClusterId()),
		obj.GetClusterName(),
		obj.GetAnnotations(),
		obj.GetPriority(),
		obj.GetImagePullSecrets(),
		obj.GetServiceAccount(),
		obj.GetServiceAccountPermissionLevel(),
		obj.GetRiskScore(),
		serialized,
	}

	finalStr := "INSERT INTO deployments (Id, Name, Type, Namespace, NamespaceId, OrchestratorComponent, Labels, PodLabels, Created, ClusterId, ClusterName, Annotations, Priority, ImagePullSecrets, ServiceAccount, ServiceAccountPermissionLevel, RiskScore, serialized) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) ON CONFLICT(Id) DO UPDATE SET Id = EXCLUDED.Id, Name = EXCLUDED.Name, Type = EXCLUDED.Type, Namespace = EXCLUDED.Namespace, NamespaceId = EXCLUDED.NamespaceId, OrchestratorComponent = EXCLUDED.OrchestratorComponent, Labels = EXCLUDED.Labels, PodLabels = EXCLUDED.PodLabels, Created = EXCLUDED.Created, ClusterId = EXCLUDED.ClusterId, ClusterName = EXCLUDED.ClusterName, Annotations = EXCLUDED.Annotations, Priority = EXCLUDED.Priority, ImagePullSecrets = EXCLUDED.ImagePullSecrets, ServiceAccount = EXCLUDED.ServiceAccount, ServiceAccountPermissionLevel = EXCLUDED.ServiceAccountPermissionLevel, RiskScore = EXCLUDED.RiskScore, serialized = EXCLUDED.serialized"
	batch.Queue(finalStr, values...)

	var query string

	for childIdx, child := range obj.GetContainers() {
		if err := insertIntoDeploymentsContainers(ctx, batch, child, obj.GetId(), childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_containers where deployments_Id = $1 AND idx >= $2"
	batch.Queue(query, pgutils.NilOrUUID(obj.GetId()), len(obj.GetContainers()))
	for childIdx, child := range obj.GetPorts() {
		if err := insertIntoDeploymentsPorts(ctx, batch, child, obj.GetId(), childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_ports where deployments_Id = $1 AND idx >= $2"
	batch.Queue(query, pgutils.NilOrUUID(obj.GetId()), len(obj.GetPorts()))
	return nil
}

func insertIntoDeploymentsContainers(ctx context.Context, batch *pgx.Batch, obj *storage.Container, deploymentID string, idx int) error {
	if pgutils.NilOrUUID(deploymentID) == nil {
		utils.Should(errors.Errorf("deploymentID is not a valid uuid -- %q", deploymentID))
		return nil
	}

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(deploymentID),
		idx,
		obj.GetImage().GetId(),
		obj.GetImage().GetName().GetRegistry(),
		obj.GetImage().GetName().GetRemote(),
		obj.GetImage().GetName().GetTag(),
		obj.GetImage().GetName().GetFullName(),
		obj.GetSecurityContext().GetPrivileged(),
		obj.GetSecurityContext().GetDropCapabilities(),
		obj.GetSecurityContext().GetAddCapabilities(),
		obj.GetSecurityContext().GetReadOnlyRootFilesystem(),
		obj.GetResources().GetCpuCoresRequest(),
		obj.GetResources().GetCpuCoresLimit(),
		obj.GetResources().GetMemoryMbRequest(),
		obj.GetResources().GetMemoryMbLimit(),
	}

	finalStr := "INSERT INTO deployments_containers (deployments_Id, idx, Image_Id, Image_Name_Registry, Image_Name_Remote, Image_Name_Tag, Image_Name_FullName, SecurityContext_Privileged, SecurityContext_DropCapabilities, SecurityContext_AddCapabilities, SecurityContext_ReadOnlyRootFilesystem, Resources_CpuCoresRequest, Resources_CpuCoresLimit, Resources_MemoryMbRequest, Resources_MemoryMbLimit) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) ON CONFLICT(deployments_Id, idx) DO UPDATE SET deployments_Id = EXCLUDED.deployments_Id, idx = EXCLUDED.idx, Image_Id = EXCLUDED.Image_Id, Image_Name_Registry = EXCLUDED.Image_Name_Registry, Image_Name_Remote = EXCLUDED.Image_Name_Remote, Image_Name_Tag = EXCLUDED.Image_Name_Tag, Image_Name_FullName = EXCLUDED.Image_Name_FullName, SecurityContext_Privileged = EXCLUDED.SecurityContext_Privileged, SecurityContext_DropCapabilities = EXCLUDED.SecurityContext_DropCapabilities, SecurityContext_AddCapabilities = EXCLUDED.SecurityContext_AddCapabilities, SecurityContext_ReadOnlyRootFilesystem = EXCLUDED.SecurityContext_ReadOnlyRootFilesystem, Resources_CpuCoresRequest = EXCLUDED.Resources_CpuCoresRequest, Resources_CpuCoresLimit = EXCLUDED.Resources_CpuCoresLimit, Resources_MemoryMbRequest = EXCLUDED.Resources_MemoryMbRequest, Resources_MemoryMbLimit = EXCLUDED.Resources_MemoryMbLimit"
	batch.Queue(finalStr, values...)

	var query string

	for childIdx, child := range obj.GetConfig().GetEnv() {
		if err := insertIntoDeploymentsContainersEnvs(ctx, batch, child, deploymentID, idx, childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_containers_envs where deployments_Id = $1 AND deployments_containers_idx = $2 AND idx >= $3"
	batch.Queue(query, pgutils.NilOrUUID(deploymentID), idx, len(obj.GetConfig().GetEnv()))
	for childIdx, child := range obj.GetVolumes() {
		if err := insertIntoDeploymentsContainersVolumes(ctx, batch, child, deploymentID, idx, childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_containers_volumes where deployments_Id = $1 AND deployments_containers_idx = $2 AND idx >= $3"
	batch.Queue(query, pgutils.NilOrUUID(deploymentID), idx, len(obj.GetVolumes()))
	for childIdx, child := range obj.GetSecrets() {
		if err := insertIntoDeploymentsContainersSecrets(ctx, batch, child, deploymentID, idx, childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_containers_secrets where deployments_Id = $1 AND deployments_containers_idx = $2 AND idx >= $3"
	batch.Queue(query, pgutils.NilOrUUID(deploymentID), idx, len(obj.GetSecrets()))
	return nil
}

func insertIntoDeploymentsContainersEnvs(_ context.Context, batch *pgx.Batch, obj *storage.ContainerConfig_EnvironmentConfig, deploymentID string, deploymentsContainersIdx int, idx int) error {
	if pgutils.NilOrUUID(deploymentID) == nil {
		utils.Should(errors.Errorf("deploymentID is not a valid uuid -- %q", deploymentID))
		return nil
	}

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(deploymentID),
		deploymentsContainersIdx,
		idx,
		obj.GetKey(),
		obj.GetValue(),
		obj.GetEnvVarSource(),
	}

	finalStr := "INSERT INTO deployments_containers_envs (deployments_Id, deployments_containers_idx, idx, Key, Value, EnvVarSource) VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT(deployments_Id, deployments_containers_idx, idx) DO UPDATE SET deployments_Id = EXCLUDED.deployments_Id, deployments_containers_idx = EXCLUDED.deployments_containers_idx, idx = EXCLUDED.idx, Key = EXCLUDED.Key, Value = EXCLUDED.Value, EnvVarSource = EXCLUDED.EnvVarSource"
	batch.Queue(finalStr, values...)

	values := []interface{}{
		"EBPF",
		"COLLECTION_METHOD",
		"KERNEL_MODULE",
	}

	updateCollectionMethodStr := "UPDATE deployments_containers_envs SET value = $1 WHERE key = $2 AND value = $3"
	batch.Queue(updateCollectionMethodStr, values...)

	return nil
}

func insertIntoDeploymentsContainersVolumes(_ context.Context, batch *pgx.Batch, obj *storage.Volume, deploymentID string, deploymentsContainersIdx int, idx int) error {
	if pgutils.NilOrUUID(deploymentID) == nil {
		utils.Should(errors.Errorf("deploymentID is not a valid uuid -- %q", deploymentID))
		return nil
	}

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(deploymentID),
		deploymentsContainersIdx,
		idx,
		obj.GetName(),
		obj.GetSource(),
		obj.GetDestination(),
		obj.GetReadOnly(),
		obj.GetType(),
	}

	finalStr := "INSERT INTO deployments_containers_volumes (deployments_Id, deployments_containers_idx, idx, Name, Source, Destination, ReadOnly, Type) VALUES($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT(deployments_Id, deployments_containers_idx, idx) DO UPDATE SET deployments_Id = EXCLUDED.deployments_Id, deployments_containers_idx = EXCLUDED.deployments_containers_idx, idx = EXCLUDED.idx, Name = EXCLUDED.Name, Source = EXCLUDED.Source, Destination = EXCLUDED.Destination, ReadOnly = EXCLUDED.ReadOnly, Type = EXCLUDED.Type"
	batch.Queue(finalStr, values...)

	return nil
}

func insertIntoDeploymentsContainersSecrets(_ context.Context, batch *pgx.Batch, obj *storage.EmbeddedSecret, deploymentID string, deploymentsContainersIdx int, idx int) error {
	if pgutils.NilOrUUID(deploymentID) == nil {
		utils.Should(errors.Errorf("deploymentID is not a valid uuid -- %q", deploymentID))
		return nil
	}

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(deploymentID),
		deploymentsContainersIdx,
		idx,
		obj.GetName(),
		obj.GetPath(),
	}

	finalStr := "INSERT INTO deployments_containers_secrets (deployments_Id, deployments_containers_idx, idx, Name, Path) VALUES($1, $2, $3, $4, $5) ON CONFLICT(deployments_Id, deployments_containers_idx, idx) DO UPDATE SET deployments_Id = EXCLUDED.deployments_Id, deployments_containers_idx = EXCLUDED.deployments_containers_idx, idx = EXCLUDED.idx, Name = EXCLUDED.Name, Path = EXCLUDED.Path"
	batch.Queue(finalStr, values...)

	return nil
}

func insertIntoDeploymentsPorts(ctx context.Context, batch *pgx.Batch, obj *storage.PortConfig, deploymentID string, idx int) error {
	if pgutils.NilOrUUID(deploymentID) == nil {
		utils.Should(errors.Errorf("deploymentID is not a valid uuid -- %q", deploymentID))
		return nil
	}

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(deploymentID),
		idx,
		obj.GetContainerPort(),
		obj.GetProtocol(),
		obj.GetExposure(),
	}

	finalStr := "INSERT INTO deployments_ports (deployments_Id, idx, ContainerPort, Protocol, Exposure) VALUES($1, $2, $3, $4, $5) ON CONFLICT(deployments_Id, idx) DO UPDATE SET deployments_Id = EXCLUDED.deployments_Id, idx = EXCLUDED.idx, ContainerPort = EXCLUDED.ContainerPort, Protocol = EXCLUDED.Protocol, Exposure = EXCLUDED.Exposure"
	batch.Queue(finalStr, values...)

	var query string

	for childIdx, child := range obj.GetExposureInfos() {
		if err := insertIntoDeploymentsPortsExposureInfos(ctx, batch, child, deploymentID, idx, childIdx); err != nil {
			return err
		}
	}

	query = "delete from deployments_ports_exposure_infos where deployments_Id = $1 AND deployments_ports_idx = $2 AND idx >= $3"
	batch.Queue(query, pgutils.NilOrUUID(deploymentID), idx, len(obj.GetExposureInfos()))
	return nil
}

func insertIntoDeploymentsPortsExposureInfos(_ context.Context, batch *pgx.Batch, obj *storage.PortConfig_ExposureInfo, deploymentID string, deploymentsPortsIdx int, idx int) error {
	if pgutils.NilOrUUID(deploymentID) == nil {
		utils.Should(errors.Errorf("deploymentID is not a valid uuid -- %q", deploymentID))
		return nil
	}

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(deploymentID),
		deploymentsPortsIdx,
		idx,
		obj.GetLevel(),
		obj.GetServiceName(),
		obj.GetServicePort(),
		obj.GetNodePort(),
		obj.GetExternalIps(),
		obj.GetExternalHostnames(),
	}

	finalStr := "INSERT INTO deployments_ports_exposure_infos (deployments_Id, deployments_ports_idx, idx, Level, ServiceName, ServicePort, NodePort, ExternalIps, ExternalHostnames) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT(deployments_Id, deployments_ports_idx, idx) DO UPDATE SET deployments_Id = EXCLUDED.deployments_Id, deployments_ports_idx = EXCLUDED.deployments_ports_idx, idx = EXCLUDED.idx, Level = EXCLUDED.Level, ServiceName = EXCLUDED.ServiceName, ServicePort = EXCLUDED.ServicePort, NodePort = EXCLUDED.NodePort, ExternalIps = EXCLUDED.ExternalIps, ExternalHostnames = EXCLUDED.ExternalHostnames"
	batch.Queue(finalStr, values...)

	return nil
}

func (s *storeImpl) copyFromDeployments(ctx context.Context, tx *postgres.Tx, objs ...*storage.Deployment) error {

	inputRows := [][]interface{}{}

	var err error

	// This is a copy so first we must delete the rows and re-add them
	// Which is essentially the desired behaviour of an upsert.
	var deletes []string

	copyCols := []string{

		"id",

		"name",

		"type",

		"namespace",

		"namespaceid",

		"orchestratorcomponent",

		"labels",

		"podlabels",

		"created",

		"clusterid",

		"clustername",

		"annotations",

		"priority",

		"imagepullsecrets",

		"serviceaccount",

		"serviceaccountpermissionlevel",

		"riskscore",

		"serialized",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj in the loop is not used as it only consists of the parent id and the idx.  Putting this here as a stop gap to simply use the object.  %s", obj)

		serialized, marshalErr := obj.Marshal()
		if marshalErr != nil {
			return marshalErr
		}

		if pgutils.NilOrUUID(obj.GetId()) == nil {
			log.Warnf("Id is not a valid uuid -- %q", obj.GetId())
			continue
		}

		inputRows = append(inputRows, []interface{}{

			pgutils.NilOrUUID(obj.GetId()),

			obj.GetName(),

			obj.GetType(),

			obj.GetNamespace(),

			pgutils.NilOrUUID(obj.GetNamespaceId()),

			obj.GetOrchestratorComponent(),

			obj.GetLabels(),

			obj.GetPodLabels(),

			pgutils.NilOrTime(obj.GetCreated()),

			pgutils.NilOrUUID(obj.GetClusterId()),

			obj.GetClusterName(),

			obj.GetAnnotations(),

			obj.GetPriority(),

			obj.GetImagePullSecrets(),

			obj.GetServiceAccount(),

			obj.GetServiceAccountPermissionLevel(),

			obj.GetRiskScore(),

			serialized,
		})

		// Add the id to be deleted.
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

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"deployments"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	for idx, obj := range objs {
		_ = idx // idx may or may not be used depending on how nested we are, so avoid compile-time errors.

		if err = s.copyFromDeploymentsContainers(ctx, tx, obj.GetId(), obj.GetContainers()...); err != nil {
			return err
		}
		if err = s.copyFromDeploymentsPorts(ctx, tx, obj.GetId(), obj.GetPorts()...); err != nil {
			return err
		}
	}

	return err
}

func (s *storeImpl) copyFromDeploymentsContainers(ctx context.Context, tx *postgres.Tx, deploymentID string, objs ...*storage.Container) error {

	inputRows := [][]interface{}{}

	var err error

	copyCols := []string{

		"deployments_id",

		"idx",

		"image_id",

		"image_name_registry",

		"image_name_remote",

		"image_name_tag",

		"image_name_fullname",

		"securitycontext_privileged",

		"securitycontext_dropcapabilities",

		"securitycontext_addcapabilities",

		"securitycontext_readonlyrootfilesystem",

		"resources_cpucoresrequest",

		"resources_cpucoreslimit",

		"resources_memorymbrequest",

		"resources_memorymblimit",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj in the loop is not used as it only consists of the parent id and the idx.  Putting this here as a stop gap to simply use the object.  %s", obj)

		if pgutils.NilOrUUID(deploymentID) == nil {
			log.Warnf("deploymentID is not a valid uuid -- %q", deploymentID)
			continue
		}

		inputRows = append(inputRows, []interface{}{

			pgutils.NilOrUUID(deploymentID),

			idx,

			obj.GetImage().GetId(),

			obj.GetImage().GetName().GetRegistry(),

			obj.GetImage().GetName().GetRemote(),

			obj.GetImage().GetName().GetTag(),

			obj.GetImage().GetName().GetFullName(),

			obj.GetSecurityContext().GetPrivileged(),

			obj.GetSecurityContext().GetDropCapabilities(),

			obj.GetSecurityContext().GetAddCapabilities(),

			obj.GetSecurityContext().GetReadOnlyRootFilesystem(),

			obj.GetResources().GetCpuCoresRequest(),

			obj.GetResources().GetCpuCoresLimit(),

			obj.GetResources().GetMemoryMbRequest(),

			obj.GetResources().GetMemoryMbLimit(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"deployments_containers"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	for idx, obj := range objs {
		_ = idx // idx may or may not be used depending on how nested we are, so avoid compile-time errors.

		if err = s.copyFromDeploymentsContainersEnvs(ctx, tx, deploymentID, idx, obj.GetConfig().GetEnv()...); err != nil {
			return err
		}
		if err = s.copyFromDeploymentsContainersVolumes(ctx, tx, deploymentID, idx, obj.GetVolumes()...); err != nil {
			return err
		}
		if err = s.copyFromDeploymentsContainersSecrets(ctx, tx, deploymentID, idx, obj.GetSecrets()...); err != nil {
			return err
		}
	}

	return err
}

func (s *storeImpl) copyFromDeploymentsContainersEnvs(ctx context.Context, tx *postgres.Tx, deploymentID string, deploymentsContainersIdx int, objs ...*storage.ContainerConfig_EnvironmentConfig) error {

	inputRows := [][]interface{}{}

	var err error

	copyCols := []string{

		"deployments_id",

		"deployments_containers_idx",

		"idx",

		"key",

		"value",

		"envvarsource",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj in the loop is not used as it only consists of the parent id and the idx.  Putting this here as a stop gap to simply use the object.  %s", obj)

		if pgutils.NilOrUUID(deploymentID) == nil {
			log.Warnf("deploymentID is not a valid uuid -- %q", deploymentID)
			continue
		}

		inputRows = append(inputRows, []interface{}{

			pgutils.NilOrUUID(deploymentID),

			deploymentsContainersIdx,

			idx,

			obj.GetKey(),

			obj.GetValue(),

			obj.GetEnvVarSource(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"deployments_containers_envs"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return err
}

func (s *storeImpl) copyFromDeploymentsContainersVolumes(ctx context.Context, tx *postgres.Tx, deploymentID string, deploymentsContainersIdx int, objs ...*storage.Volume) error {

	inputRows := [][]interface{}{}

	var err error

	copyCols := []string{

		"deployments_id",

		"deployments_containers_idx",

		"idx",

		"name",

		"source",

		"destination",

		"readonly",

		"type",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj in the loop is not used as it only consists of the parent id and the idx.  Putting this here as a stop gap to simply use the object.  %s", obj)

		if pgutils.NilOrUUID(deploymentID) == nil {
			log.Warnf("deploymentID is not a valid uuid -- %q", deploymentID)
			continue
		}

		inputRows = append(inputRows, []interface{}{

			pgutils.NilOrUUID(deploymentID),

			deploymentsContainersIdx,

			idx,

			obj.GetName(),

			obj.GetSource(),

			obj.GetDestination(),

			obj.GetReadOnly(),

			obj.GetType(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"deployments_containers_volumes"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return err
}

func (s *storeImpl) copyFromDeploymentsContainersSecrets(ctx context.Context, tx *postgres.Tx, deploymentID string, deploymentsContainersIdx int, objs ...*storage.EmbeddedSecret) error {

	inputRows := [][]interface{}{}

	var err error

	copyCols := []string{

		"deployments_id",

		"deployments_containers_idx",

		"idx",

		"name",

		"path",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj in the loop is not used as it only consists of the parent id and the idx.  Putting this here as a stop gap to simply use the object.  %s", obj)

		if pgutils.NilOrUUID(deploymentID) == nil {
			log.Warnf("deploymentID is not a valid uuid -- %q", deploymentID)
			continue
		}

		inputRows = append(inputRows, []interface{}{

			pgutils.NilOrUUID(deploymentID),

			deploymentsContainersIdx,

			idx,

			obj.GetName(),

			obj.GetPath(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"deployments_containers_secrets"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return err
}

func (s *storeImpl) copyFromDeploymentsPorts(ctx context.Context, tx *postgres.Tx, deploymentID string, objs ...*storage.PortConfig) error {

	inputRows := [][]interface{}{}

	var err error

	copyCols := []string{

		"deployments_id",

		"idx",

		"containerport",

		"protocol",

		"exposure",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj in the loop is not used as it only consists of the parent id and the idx.  Putting this here as a stop gap to simply use the object.  %s", obj)

		if pgutils.NilOrUUID(deploymentID) == nil {
			log.Warnf("deploymentID is not a valid uuid -- %q", deploymentID)
			continue
		}

		inputRows = append(inputRows, []interface{}{

			pgutils.NilOrUUID(deploymentID),

			idx,

			obj.GetContainerPort(),

			obj.GetProtocol(),

			obj.GetExposure(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"deployments_ports"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	for idx, obj := range objs {
		_ = idx // idx may or may not be used depending on how nested we are, so avoid compile-time errors.

		if err = s.copyFromDeploymentsPortsExposureInfos(ctx, tx, deploymentID, idx, obj.GetExposureInfos()...); err != nil {
			return err
		}
	}

	return err
}

func (s *storeImpl) copyFromDeploymentsPortsExposureInfos(ctx context.Context, tx *postgres.Tx, deploymentID string, deploymentsPortsIdx int, objs ...*storage.PortConfig_ExposureInfo) error {

	inputRows := [][]interface{}{}

	var err error

	copyCols := []string{

		"deployments_id",

		"deployments_ports_idx",

		"idx",

		"level",

		"servicename",

		"serviceport",

		"nodeport",

		"externalips",

		"externalhostnames",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj in the loop is not used as it only consists of the parent id and the idx.  Putting this here as a stop gap to simply use the object.  %s", obj)

		if pgutils.NilOrUUID(deploymentID) == nil {
			log.Warnf("deploymentID is not a valid uuid -- %q", deploymentID)
			continue
		}

		inputRows = append(inputRows, []interface{}{

			pgutils.NilOrUUID(deploymentID),

			deploymentsPortsIdx,

			idx,

			obj.GetLevel(),

			obj.GetServiceName(),

			obj.GetServicePort(),

			obj.GetNodePort(),

			obj.GetExternalIps(),

			obj.GetExternalHostnames(),
		})

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"deployments_ports_exposure_infos"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return err
}

func (s *storeImpl) copyFrom(ctx context.Context, objs ...*storage.Deployment) error {
	conn, release, err := s.acquireConn(ctx, ops.Get, "Deployment")
	if err != nil {
		return err
	}
	defer release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	if err := s.copyFromDeployments(ctx, tx, objs...); err != nil {
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

func (s *storeImpl) upsert(ctx context.Context, objs ...*storage.Deployment) error {
	conn, release, err := s.acquireConn(ctx, ops.Get, "Deployment")
	if err != nil {
		return err
	}
	defer release()

	for _, obj := range objs {
		batch := &pgx.Batch{}
		if err := insertIntoDeployments(ctx, batch, obj); err != nil {
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

func (s *storeImpl) Upsert(ctx context.Context, obj *storage.Deployment) error {

	return pgutils.Retry(func() error {
		return s.upsert(ctx, obj)
	})
}

func (s *storeImpl) UpsertMany(ctx context.Context, objs []*storage.Deployment) error {

	return pgutils.Retry(func() error {
		// Lock since copyFrom requires a delete first before being executed.  If multiple processes are updating
		// same subset of rows, both deletes could occur before the copyFrom resulting in unique constraint
		// violations
		s.mutex.Lock()
		defer s.mutex.Unlock()

		if len(objs) < batchAfter {
			return s.upsert(ctx, objs...)
		}
		return s.copyFrom(ctx, objs...)
	})
}

// Count returns the number of objects in the store
func (s *storeImpl) Count(ctx context.Context) (int, error) {

	var sacQueryFilter *v1.Query

	return pgSearch.RunCountRequestForSchema(ctx, schema, sacQueryFilter, s.db)
}

// Exists returns if the id exists in the store
func (s *storeImpl) Exists(ctx context.Context, id string) (bool, error) {

	var sacQueryFilter *v1.Query

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	count, err := pgSearch.RunCountRequestForSchema(ctx, schema, q, s.db)
	// With joins and multiple paths to the scoping resources, it can happen that the Count query for an object identifier
	// returns more than 1, despite the fact that the identifier is unique in the table.
	return count > 0, err
}

// Get returns the object, if it exists from the store
func (s *storeImpl) Get(ctx context.Context, id string) (*storage.Deployment, bool, error) {

	var sacQueryFilter *v1.Query

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	data, err := pgSearch.RunGetQueryForSchema[storage.Deployment](ctx, schema, q, s.db)
	if err != nil {
		return nil, false, pgutils.ErrNilIfNoRows(err)
	}

	return data, true, nil
}

func (s *storeImpl) acquireConn(ctx context.Context, _ ops.Op, _ string) (*postgres.Conn, func(), error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	return conn, conn.Release, nil
}

// Delete removes the specified ID from the store
func (s *storeImpl) Delete(ctx context.Context, id string) error {

	var sacQueryFilter *v1.Query

	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(id).ProtoQuery(),
	)

	return pgSearch.RunDeleteRequestForSchema(ctx, schema, q, s.db)
}

// DeleteByQuery removes the objects based on the passed query
func (s *storeImpl) DeleteByQuery(ctx context.Context, query *v1.Query) error {

	var sacQueryFilter *v1.Query

	q := search.ConjunctionQuery(
		sacQueryFilter,
		query,
	)

	return pgSearch.RunDeleteRequestForSchema(ctx, schema, q, s.db)
}

// GetIDs returns all the IDs for the store
func (s *storeImpl) GetIDs(ctx context.Context) ([]string, error) {
	var sacQueryFilter *v1.Query
	result, err := pgSearch.RunSearchRequestForSchema(ctx, schema, sacQueryFilter, s.db)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(result))
	for _, entry := range result {
		ids = append(ids, entry.ID)
	}

	return ids, nil
}

// GetMany returns the objects specified by the IDs or the index in the missing indices slice
func (s *storeImpl) GetMany(ctx context.Context, ids []string) ([]*storage.Deployment, []int, error) {

	if len(ids) == 0 {
		return nil, nil, nil
	}

	var sacQueryFilter *v1.Query
	q := search.ConjunctionQuery(
		sacQueryFilter,
		search.NewQueryBuilder().AddDocIDs(ids...).ProtoQuery(),
	)

	rows, err := pgSearch.RunGetManyQueryForSchema[storage.Deployment](ctx, schema, q, s.db)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			missingIndices := make([]int, 0, len(ids))
			for i := range ids {
				missingIndices = append(missingIndices, i)
			}
			return nil, missingIndices, nil
		}
		return nil, nil, err
	}
	resultsByID := make(map[string]*storage.Deployment, len(rows))
	for _, msg := range rows {
		resultsByID[msg.GetId()] = msg
	}
	missingIndices := make([]int, 0, len(ids)-len(resultsByID))
	// It is important that the elems are populated in the same order as the input ids
	// slice, since some calling code relies on that to maintain order.
	elems := make([]*storage.Deployment, 0, len(resultsByID))
	for i, id := range ids {
		if result, ok := resultsByID[id]; !ok {
			missingIndices = append(missingIndices, i)
		} else {
			elems = append(elems, result)
		}
	}
	return elems, missingIndices, nil
}

// GetByQuery returns the objects matching the query
func (s *storeImpl) GetByQuery(ctx context.Context, query *v1.Query) ([]*storage.Deployment, error) {

	var sacQueryFilter *v1.Query
	q := search.ConjunctionQuery(
		sacQueryFilter,
		query,
	)

	rows, err := pgSearch.RunGetManyQueryForSchema[storage.Deployment](ctx, schema, q, s.db)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return rows, nil
}

// Delete removes the specified IDs from the store
func (s *storeImpl) DeleteMany(ctx context.Context, ids []string) error {

	var sacQueryFilter *v1.Query

	// Batch the deletes
	localBatchSize := deleteBatchSize
	numRecordsToDelete := len(ids)
	for {
		if len(ids) == 0 {
			break
		}

		if len(ids) < localBatchSize {
			localBatchSize = len(ids)
		}

		idBatch := ids[:localBatchSize]
		q := search.ConjunctionQuery(
			sacQueryFilter,
			search.NewQueryBuilder().AddDocIDs(idBatch...).ProtoQuery(),
		)

		if err := pgSearch.RunDeleteRequestForSchema(ctx, schema, q, s.db); err != nil {
			err = errors.Wrapf(err, "unable to delete the records.  Successfully deleted %d out of %d", numRecordsToDelete-len(ids), numRecordsToDelete)
			log.Error(err)
			return err
		}

		// Move the slice forward to start the next batch
		ids = ids[localBatchSize:]
	}

	return nil
}

// Walk iterates over all of the objects in the store and applies the closure
func (s *storeImpl) Walk(ctx context.Context, fn func(obj *storage.Deployment) error) error {
	var sacQueryFilter *v1.Query
	fetcher, closer, err := pgSearch.RunCursorQueryForSchema[storage.Deployment](ctx, schema, sacQueryFilter, s.db)
	if err != nil {
		return err
	}
	defer closer()
	for {
		rows, err := fetcher(cursorBatchSize)
		if err != nil {
			return pgutils.ErrNilIfNoRows(err)
		}
		for _, data := range rows {
			if err := fn(data); err != nil {
				return err
			}
		}
		if len(rows) != cursorBatchSize {
			break
		}
	}
	return nil
}

//// Used for testing

func dropTableDeployments(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments CASCADE")
	dropTableDeploymentsContainers(ctx, db)
	dropTableDeploymentsPorts(ctx, db)

}

func dropTableDeploymentsContainers(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments_containers CASCADE")
	dropTableDeploymentsContainersEnvs(ctx, db)
	dropTableDeploymentsContainersVolumes(ctx, db)
	dropTableDeploymentsContainersSecrets(ctx, db)

}

func dropTableDeploymentsContainersEnvs(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments_containers_envs CASCADE")

}

func dropTableDeploymentsContainersVolumes(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments_containers_volumes CASCADE")

}

func dropTableDeploymentsContainersSecrets(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments_containers_secrets CASCADE")

}

func dropTableDeploymentsPorts(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments_ports CASCADE")
	dropTableDeploymentsPortsExposureInfos(ctx, db)

}

func dropTableDeploymentsPortsExposureInfos(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS deployments_ports_exposure_infos CASCADE")

}

// Destroy -- destroy table for test
func Destroy(ctx context.Context, db postgres.DB) {
	dropTableDeployments(ctx, db)
}
