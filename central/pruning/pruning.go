package pruning

import (
	"context"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	alertDatastore "github.com/stackrox/rox/central/alert/datastore"
	clusterDatastore "github.com/stackrox/rox/central/cluster/datastore"
	configDatastore "github.com/stackrox/rox/central/config/datastore"
	deploymentDatastore "github.com/stackrox/rox/central/deployment/datastore"
	imageDatastore "github.com/stackrox/rox/central/image/datastore"
	imageComponentDatastore "github.com/stackrox/rox/central/imagecomponent/datastore"
	networkFlowDatastore "github.com/stackrox/rox/central/networkgraph/flow/datastore"
	podDatastore "github.com/stackrox/rox/central/pod/datastore"
	processDatastore "github.com/stackrox/rox/central/processindicator/datastore"
	processWhitelistDatastore "github.com/stackrox/rox/central/processwhitelist/datastore"
	riskDataStore "github.com/stackrox/rox/central/risk/datastore"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/protoutils"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/utils"
)

const (
	pruneInterval       = 1 * time.Hour
	orphanWindow        = 30 * time.Minute
	whitelistBatchLimit = 10000
)

var (
	log        = logging.LoggerForModule()
	pruningCtx = sac.WithAllAccess(context.Background())
)

// GarbageCollector implements a generic garbage collection mechanism
type GarbageCollector interface {
	Start()
	Stop()
}

func newGarbageCollector(alerts alertDatastore.DataStore,
	images imageDatastore.DataStore,
	clusters clusterDatastore.DataStore,
	deployments deploymentDatastore.DataStore,
	pods podDatastore.DataStore,
	processes processDatastore.DataStore,
	processwhitelist processWhitelistDatastore.DataStore,
	networkflows networkFlowDatastore.ClusterDataStore,
	config configDatastore.DataStore,
	imageComponents imageComponentDatastore.DataStore,
	risks riskDataStore.DataStore) GarbageCollector {
	return &garbageCollectorImpl{
		alerts:           alerts,
		clusters:         clusters,
		images:           images,
		imageComponents:  imageComponents,
		deployments:      deployments,
		pods:             pods,
		processes:        processes,
		processwhitelist: processwhitelist,
		networkflows:     networkflows,
		config:           config,
		risks:            risks,
		stopSig:          concurrency.NewSignal(),
		stoppedSig:       concurrency.NewSignal(),
	}
}

type garbageCollectorImpl struct {
	alerts           alertDatastore.DataStore
	clusters         clusterDatastore.DataStore
	images           imageDatastore.DataStore
	imageComponents  imageComponentDatastore.DataStore
	deployments      deploymentDatastore.DataStore
	pods             podDatastore.DataStore
	processes        processDatastore.DataStore
	processwhitelist processWhitelistDatastore.DataStore
	networkflows     networkFlowDatastore.ClusterDataStore
	config           configDatastore.DataStore
	risks            riskDataStore.DataStore
	stopSig          concurrency.Signal
	stoppedSig       concurrency.Signal
}

func (g *garbageCollectorImpl) Start() {
	go g.runGC()
}

func (g *garbageCollectorImpl) pruneBasedOnConfig() {
	config, err := g.config.GetConfig(pruningCtx)
	if err != nil {
		log.Error(err)
		return
	}
	if config == nil {
		log.Error("UNEXPECTED: Got nil config")
		return
	}
	log.Info("[Pruning] Starting a garbage collection cycle")
	pvtConfig := config.GetPrivateConfig()
	// Run collection initially then run on a ticker
	g.collectImages(pvtConfig)
	g.collectAlerts(pvtConfig)
	g.removeOrphanedResources()
	g.removeOrphanedRisks()
	log.Info("[Pruning] Finished garbage collection cycle")
}

func (g *garbageCollectorImpl) runGC() {
	g.pruneBasedOnConfig()

	t := time.NewTicker(pruneInterval)
	for {
		select {
		case <-t.C:
			g.pruneBasedOnConfig()
		case <-g.stopSig.Done():
			g.stoppedSig.Signal()
			return
		}
	}
}

func (g *garbageCollectorImpl) removeOrphanedResources() {
	deploymentIDs, err := g.deployments.GetDeploymentIDs()
	if err != nil {
		log.Error(errors.Wrap(err, "unable to fetch deployment IDs in pruning"))
		return
	}
	deploymentSet := set.NewFrozenStringSet(deploymentIDs...)
	podIDs, err := g.pods.GetPodIDs()
	if err != nil {
		log.Error(errors.Wrap(err, "unable to fetch pod IDs in pruning"))
		return
	}
	g.removeOrphanedProcesses(deploymentSet, set.NewFrozenStringSet(podIDs...))
	g.removeOrphanedProcessWhitelists(deploymentSet)
	g.markOrphanedAlertsAsResolved(deploymentSet)
	g.removeOrphanedNetworkFlows(deploymentSet)
}

func (g *garbageCollectorImpl) removeOrphanedProcesses(deploymentIDs, podIDs set.FrozenStringSet) {
	var processesToPrune []string
	now := types.TimestampNow()
	err := g.processes.WalkAll(pruningCtx, func(pi *storage.ProcessIndicator) error {
		if pi.GetPodUid() != "" && podIDs.Contains(pi.GetPodUid()) {
			return nil
		}
		if pi.GetPodUid() == "" && deploymentIDs.Contains(pi.GetDeploymentId()) {
			return nil
		}
		if protoutils.Sub(now, pi.GetSignal().GetTime()) < orphanWindow {
			return nil
		}
		processesToPrune = append(processesToPrune, pi.GetId())
		return nil
	})
	if err != nil {
		log.Error(errors.Wrap(err, "unable to walk processes and mark for pruning"))
		return
	}
	log.Infof("[Process pruning] Found %d orphaned processes", len(processesToPrune))
	if err := g.processes.RemoveProcessIndicators(pruningCtx, processesToPrune); err != nil {
		log.Error(errors.Wrap(err, "error removing process indicators"))
	}
}

func (g *garbageCollectorImpl) removeOrphanedProcessWhitelists(deployments set.FrozenStringSet) {
	var whitelistBatchOffset, prunedProcessWhitelists int32
	for {
		allQuery := &v1.Query{
			Pagination: &v1.QueryPagination{
				Offset: whitelistBatchOffset,
				Limit:  whitelistBatchLimit,
			},
		}

		res, err := g.processwhitelist.Search(pruningCtx, allQuery)
		if err != nil {
			log.Error(errors.Wrap(err, "error searching process baselines"))
			return
		}

		whitelistBatchOffset += whitelistBatchLimit
		var whitelistKeysToPrune []*storage.ProcessWhitelistKey
		for _, whitelist := range res {
			whitelistKey, err := processWhitelistDatastore.IDToKey(whitelist.ID)
			if err != nil {
				log.Error(errors.Wrapf(err, "Invalid id %s", whitelist.ID))
				continue
			}

			if !deployments.Contains(whitelistKey.GetDeploymentId()) {
				whitelistKeysToPrune = append(whitelistKeysToPrune, whitelistKey)
			}
		}

		now := types.TimestampNow()
		for _, whitelistKey := range whitelistKeysToPrune {
			whitelist, exists, err := g.processwhitelist.GetProcessWhitelist(pruningCtx, whitelistKey)
			if err != nil {
				log.Error(errors.Wrapf(err, "unable to fetch process baseline for key %v", whitelistKey))
				continue
			}

			if !exists || protoutils.Sub(now, whitelist.GetCreated()) < orphanWindow {
				continue
			}

			if err = g.processwhitelist.RemoveProcessWhitelist(pruningCtx, whitelistKey); err != nil {
				log.Error(errors.Wrapf(err, "unable to remove process baseline: %v", whitelistKey))
				continue
			}

			prunedProcessWhitelists++
		}

		if len(res) < whitelistBatchLimit {
			break
		}
	}

	log.Infof("[Process baseline pruning] Removed %d process baselines", prunedProcessWhitelists)
}

func (g *garbageCollectorImpl) markOrphanedAlertsAsResolved(deployments set.FrozenStringSet) {
	var alertsToResolve []string
	now := types.TimestampNow()
	err := g.alerts.WalkAll(pruningCtx, func(alert *storage.ListAlert) error {
		// We should only remove orphaned deploy time alerts as they are not cleaned up by retention policies
		// This will only happen when there is data inconsistency
		if alert.GetLifecycleStage() != storage.LifecycleStage_DEPLOY {
			return nil
		}
		if alert.GetState() != storage.ViolationState_ACTIVE {
			return nil
		}
		if deployments.Contains(alert.GetDeployment().GetId()) {
			return nil
		}
		if protoutils.Sub(now, alert.GetTime()) < orphanWindow {
			return nil
		}
		alertsToResolve = append(alertsToResolve, alert.GetId())
		return nil
	})
	if err != nil {
		log.Error(errors.Wrap(err, "unable to walk alerts and mark for pruning"))
		return
	}
	log.Infof("[Alert pruning] Found %d orphaned alerts", len(alertsToResolve))
	for _, a := range alertsToResolve {
		if err := g.alerts.MarkAlertStale(pruningCtx, a); err != nil {
			log.Error(errors.Wrapf(err, "error marking alert %q as stale", a))
		}
	}
}

func isOrphanedDeployment(deployments set.FrozenStringSet, info *storage.NetworkEntityInfo) bool {
	return info.GetType() == storage.NetworkEntityInfo_DEPLOYMENT && !deployments.Contains(info.GetId())
}

func (g *garbageCollectorImpl) removeOrphanedNetworkFlows(deployments set.FrozenStringSet) {
	clusters, err := g.clusters.GetClusters(pruningCtx)
	if err != nil {
		utils.Should(errors.Wrap(err, "could not fetch clusters for orphaned network flows"))
		return
	}

	for _, c := range clusters {
		store, err := g.networkflows.GetFlowStore(pruningCtx, c.GetId())
		if err != nil {
			log.Errorf("error getting flow store for cluster %q: %v", c.GetId(), err)
			continue
		} else if store == nil {
			continue
		}
		now := types.TimestampNow()

		keyMatchFn := func(props *storage.NetworkFlowProperties) bool {
			return isOrphanedDeployment(deployments, props.GetSrcEntity()) ||
				isOrphanedDeployment(deployments, props.GetDstEntity())
		}
		valueMatchFn := func(flow *storage.NetworkFlow) bool {
			return flow.LastSeenTimestamp != nil && protoutils.Sub(now, flow.LastSeenTimestamp) > orphanWindow
		}
		err = store.RemoveMatchingFlows(pruningCtx, keyMatchFn, valueMatchFn)
		if err != nil {
			log.Errorf("error removing orphaned flows for cluster %q: %v", c.GetName(), err)
		}
	}
}

func (g *garbageCollectorImpl) collectImages(config *storage.PrivateConfig) {
	pruneImageAfterDays := config.GetImageRetentionDurationDays()
	qb := search.NewQueryBuilder().AddDays(search.LastUpdatedTime, int64(pruneImageAfterDays)).ProtoQuery()
	imageResults, err := g.images.Search(pruningCtx, qb)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("[Image pruning] Found %d image search results", len(imageResults))

	imagesToPrune := make([]string, 0, len(imageResults))
	for _, result := range imageResults {
		q1 := search.NewQueryBuilder().AddExactMatches(search.ImageSHA, result.ID).ProtoQuery()
		q2 := search.NewQueryBuilder().AddExactMatches(search.ContainerImageDigest, result.ID).ProtoQuery()
		deploymentResults, err := g.deployments.Search(pruningCtx, q1)
		if err != nil {
			log.Errorf("[Image pruning] searching deployments: %v", err)
			continue
		}

		podResults, err := g.pods.Search(pruningCtx, q2)
		if err != nil {
			log.Errorf("[Image pruning] searching pods: %v", err)
			continue
		}

		if len(deploymentResults) == 0 && len(podResults) == 0 {
			imagesToPrune = append(imagesToPrune, result.ID)
		}
	}
	if len(imagesToPrune) > 0 {
		log.Infof("[Image Pruning] Removing %d images", len(imagesToPrune))
		log.Debugf("[Image Pruning] Removing images %+v", imagesToPrune)
		if err := g.images.DeleteImages(pruningCtx, imagesToPrune...); err != nil {
			log.Error(err)
		}
	}
}

func getConfigValues(config *storage.PrivateConfig) (pruneResolvedDeployAfter, pruneAllRuntimeAfter, pruneDeletedRuntimeAfter int32) {
	alertRetention := config.GetAlertRetention()
	if val, ok := alertRetention.(*storage.PrivateConfig_DEPRECATEDAlertRetentionDurationDays); ok {
		global := val.DEPRECATEDAlertRetentionDurationDays
		return global, global, global

	} else if val, ok := alertRetention.(*storage.PrivateConfig_AlertConfig); ok {
		return val.AlertConfig.GetResolvedDeployRetentionDurationDays(),
			val.AlertConfig.GetAllRuntimeRetentionDurationDays(),
			val.AlertConfig.GetDeletedRuntimeRetentionDurationDays()
	}
	return 0, 0, 0
}

func (g *garbageCollectorImpl) collectAlerts(config *storage.PrivateConfig) {
	alertRetention := config.GetAlertRetention()
	if alertRetention == nil {
		log.Info("[Alert pruning] Alert pruning has been disabled.")
		return
	}

	pruneResolvedDeployAfter, pruneAllRuntimeAfter, pruneDeletedRuntimeAfter := getConfigValues(config)

	var queries []*v1.Query

	if pruneResolvedDeployAfter > 0 {
		q := search.NewQueryBuilder().
			AddStrings(search.LifecycleStage, storage.LifecycleStage_DEPLOY.String()).
			AddStrings(search.ViolationState, storage.ViolationState_RESOLVED.String()).
			AddDays(search.ViolationTime, int64(pruneResolvedDeployAfter)).
			ProtoQuery()
		queries = append(queries, q)
	}

	if pruneAllRuntimeAfter > 0 {
		q := search.NewQueryBuilder().
			AddStrings(search.LifecycleStage, storage.LifecycleStage_RUNTIME.String()).
			AddDays(search.ViolationTime, int64(pruneAllRuntimeAfter)).
			ProtoQuery()
		queries = append(queries, q)
	}

	if pruneDeletedRuntimeAfter > 0 && pruneAllRuntimeAfter != pruneDeletedRuntimeAfter {
		q := search.NewQueryBuilder().
			AddStrings(search.LifecycleStage, storage.LifecycleStage_RUNTIME.String()).
			AddDays(search.ViolationTime, int64(pruneDeletedRuntimeAfter)).
			AddBools(search.Inactive, true).
			ProtoQuery()
		queries = append(queries, q)
	}

	if len(queries) == 0 {
		log.Info("No alert retention configuration, skipping")
		return
	}

	alertResults, err := g.alerts.Search(pruningCtx,
		search.NewDisjunctionQuery(queries...))
	if err != nil {
		log.Error(err)
		return
	}

	alertsToPrune := search.ResultsToIDs(alertResults)

	if len(alertsToPrune) > 0 {
		log.Infof("[Alert pruning] Removing %d alerts", len(alertsToPrune))
		if err := g.alerts.DeleteAlerts(pruningCtx, alertsToPrune...); err != nil {
			log.Error(err)
		}
	}
}

func (g *garbageCollectorImpl) removeOrphanedRisks() {
	g.removeOrphanedDeploymentRisks()
	g.removeOrphanedImageRisks()
	g.removeOrphanedImageComponentRisks()
}

func (g *garbageCollectorImpl) removeOrphanedDeploymentRisks() {
	deploymentsWithRisk := g.getRisks(storage.RiskSubjectType_DEPLOYMENT)
	results, err := g.deployments.Search(pruningCtx, search.EmptyQuery())
	if err != nil {
		log.Error(err)
		return
	}

	prunable := deploymentsWithRisk.Difference(set.NewStringSet(search.ResultsToIDs(results)...)).AsSlice()
	log.Infof("[Risk pruning] Removing %d deployment risks", len(prunable))
	g.removeRisks(storage.RiskSubjectType_DEPLOYMENT, prunable...)
}

func (g *garbageCollectorImpl) removeOrphanedImageRisks() {
	imagesWithRisk := g.getRisks(storage.RiskSubjectType_IMAGE)
	results, err := g.images.Search(pruningCtx, search.EmptyQuery())
	if err != nil {
		log.Error(err)
		return
	}

	prunable := imagesWithRisk.Difference(set.NewStringSet(search.ResultsToIDs(results)...)).AsSlice()
	log.Infof("[Risk pruning] Removing %d image risks", len(prunable))
	g.removeRisks(storage.RiskSubjectType_IMAGE, prunable...)
}

func (g *garbageCollectorImpl) removeOrphanedImageComponentRisks() {
	componentsWithRisk := g.getRisks(storage.RiskSubjectType_IMAGE_COMPONENT)
	results, err := g.imageComponents.Search(pruningCtx, search.EmptyQuery())
	if err != nil {
		log.Error(err)
		return
	}

	prunable := componentsWithRisk.Difference(set.NewStringSet(search.ResultsToIDs(results)...)).AsSlice()
	log.Infof("[Risk pruning] Removing %d image component risks", len(prunable))
	g.removeRisks(storage.RiskSubjectType_IMAGE_COMPONENT, prunable...)
}

func (g *garbageCollectorImpl) getRisks(riskType storage.RiskSubjectType) set.StringSet {
	risks := set.NewStringSet()
	results, err := g.risks.Search(pruningCtx, search.NewQueryBuilder().AddExactMatches(search.RiskSubjectType, riskType.String()).ProtoQuery())
	if err != nil {
		log.Error(err)
		return risks
	}

	for _, result := range results {
		_, id, err := riskDataStore.GetIDParts(result.ID)
		if err != nil {
			log.Error(err)
			continue
		}
		risks.Add(id)
	}
	return risks
}

func (g *garbageCollectorImpl) removeRisks(riskType storage.RiskSubjectType, ids ...string) {
	for _, id := range ids {
		if err := g.risks.RemoveRisk(pruningCtx, id, riskType); err != nil {
			log.Error(err)
		}
	}
}

func (g *garbageCollectorImpl) Stop() {
	g.stopSig.Signal()
	<-g.stoppedSig.Done()
}
