package datastore

import (
	"context"
	"fmt"

	processIndicatorStore "github.com/stackrox/rox/central/processindicator/datastore"
	"github.com/stackrox/rox/central/processlisteningonport/store/postgres"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/uuid"
)

type datastoreImpl struct {
	storage        postgres.Store
	indicatorStore processIndicatorStore.DataStore
}

var (
	indicatorSAC = sac.ForResource(resources.Indicator)
	log          = logging.LoggerForModule()
)

func newDatastoreImpl(
	storage postgres.Store,
	indicatorDataStore processIndicatorStore.DataStore,
) *datastoreImpl {
	return &datastoreImpl{
		storage:        storage,
		indicatorStore: indicatorDataStore,
	}
}

func (ds *datastoreImpl) AddProcessListeningOnPort(
	ctx context.Context,
	portProcesses ...*storage.ProcessListeningOnPort,
) error {

	var (
		lookups       []*v1.Query
		indicatorsMap = map[string]*storage.ProcessIndicator{}
	)

	if ok, err := indicatorSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return sac.ErrResourceAccessDenied
	}

	// First query all needed process indicators in one go
	for _, val := range portProcesses {
		lookups = append(lookups,
			search.NewQueryBuilder().
				AddExactMatches(search.ContainerName, val.Process.ContainerName).
				AddExactMatches(search.PodID, val.Process.PodId).
				AddExactMatches(search.ProcessName, val.Process.ProcessName).
				AddExactMatches(search.ProcessArguments, val.Process.ProcessArgs).
				AddExactMatches(search.ProcessExecPath, val.Process.ProcessExecFilePath).
				ProtoQuery())
	}

	indicatorsQuery := search.DisjunctionQuery(lookups...)
	log.Debugf("Sending query: %s", indicatorsQuery.String())
	indicators, err := ds.indicatorStore.SearchRawProcessIndicators(ctx, indicatorsQuery)
	if err != nil {
		return err
	}

	for _, val := range indicators {
		key := fmt.Sprintf("%s %s %s %s %s",
			val.GetContainerName(),
			val.GetPodId(),
			val.GetSignal().GetName(),
			val.GetSignal().GetArgs(),
			val.GetSignal().GetExecFilePath(),
		)
		log.Infof("indicators key= %s\n", key)

		// A bit of paranoia is always good
		if old, ok := indicatorsMap[key]; ok {
			log.Warnf("An indicator %s is already present, overwrite with %s",
				old.GetId(), val.GetId())
		}

		indicatorsMap[key] = val
	}

	plopObjects := make([]*storage.ProcessListeningOnPortStorage, len(portProcesses))
	for i, val := range portProcesses {
		indicatorId := ""

		key := fmt.Sprintf("%s %s %s %s %s",
			val.Process.ContainerName,
			val.Process.PodId,
			val.Process.ProcessName,
			val.Process.ProcessArgs,
			val.Process.ProcessExecFilePath,
		)

		log.Infof("portProcess key= %s\n", key)

		if indicator, ok := indicatorsMap[key]; ok {
			indicatorId = indicator.GetId()
			log.Infof("Got indicator %s: %+v", indicatorId, indicator)
			//log.Debugf("Got indicator %s: %+v", indicatorId, indicator)
		} else {
			log.Warnf("Found no matching indicators for %s", key)
		}

		plopObjects[i] = &storage.ProcessListeningOnPortStorage{
			// XXX, ResignatingFacepalm: Use regular GENERATE ALWAYS AS
			// IDENTITY, which would require changes in store generator
			Id:                 uuid.NewV4().String(),
			Port:               val.Port,
			Protocol:           val.Protocol,
			ProcessIndicatorId: indicatorId,
			CloseTimestamp:     val.CloseTimestamp,
		}
	}

	// Not save actual PLOP objects
	err = ds.storage.UpsertMany(ctx, plopObjects)
	if err != nil {
			log.Warnf("Unable to upsert")
		return err
	}

	return nil
}

func (ds *datastoreImpl) GetProcessListeningOnPortForDeployment(
	ctx context.Context,
	namespace string,
	deploymentId string,
) (
	[]*storage.ProcessListeningOnPort, error,
) {

	portProcessList, err := ds.storage.GetProcessListeningOnPortForDeployment(ctx, namespace, deploymentId)

	if err != nil {
		return nil, err
	}

	return portProcessList, nil
}

func (ds *datastoreImpl) GetProcessListeningOnPortForNamespace(
	ctx context.Context,
	namespace string,
) (
	[]*v1.ProcessListeningOnPortWithDeploymentId, error,
) {

	portProcessList, err := ds.storage.GetProcessListeningOnPortForNamespace(ctx, namespace)

	if err != nil {
		return nil, err
	}

	return portProcessList, nil
}
