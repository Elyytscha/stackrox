package lifecycle

import (
	"context"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	deploymentDatastore "github.com/stackrox/rox/central/deployment/datastore"
	"github.com/stackrox/rox/central/detection/alertmanager"
	"github.com/stackrox/rox/central/detection/deploytime"
	"github.com/stackrox/rox/central/detection/runtime"
	"github.com/stackrox/rox/central/enrichment"
	imageDatastore "github.com/stackrox/rox/central/image/datastore"
	processIndicatorDatastore "github.com/stackrox/rox/central/processindicator/datastore"
	"github.com/stackrox/rox/central/processwhitelist"
	whitelistDataStore "github.com/stackrox/rox/central/processwhitelist/datastore"
	"github.com/stackrox/rox/central/reprocessor"
	riskManager "github.com/stackrox/rox/central/risk/manager"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/central/sensor/service/common"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/images/enricher"
	"github.com/stackrox/rox/pkg/policies"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/sync"
	"golang.org/x/time/rate"
)

var (
	lifecycleMgrCtx = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.Alert, resources.Deployment, resources.Image, resources.Indicator, resources.Policy, resources.ProcessWhitelist)))
)

type indicatorWithInjector struct {
	indicator           *storage.ProcessIndicator
	msgToSensorInjector common.MessageInjector
}

type managerImpl struct {
	reprocessor        reprocessor.Loop
	enricher           enrichment.Enricher
	riskManager        riskManager.Manager
	runtimeDetector    runtime.Detector
	deploytimeDetector deploytime.Detector
	alertManager       alertmanager.AlertManager

	deploymentDataStore deploymentDatastore.DataStore
	processesDataStore  processIndicatorDatastore.DataStore
	whitelists          whitelistDataStore.DataStore
	imageDataStore      imageDatastore.DataStore

	queuedIndicators map[string]indicatorWithInjector

	queueLock           sync.Mutex
	flushProcessingLock concurrency.TransparentMutex

	limiter *rate.Limiter
	ticker  *time.Ticker
}

func (m *managerImpl) copyAndResetIndicatorQueue() map[string]indicatorWithInjector {
	m.queueLock.Lock()
	defer m.queueLock.Unlock()
	if len(m.queuedIndicators) == 0 {
		return nil
	}
	copiedMap := m.queuedIndicators
	m.queuedIndicators = make(map[string]indicatorWithInjector)

	return copiedMap
}

func (m *managerImpl) flushQueuePeriodically() {
	defer m.ticker.Stop()
	for range m.ticker.C {
		m.flushIndicatorQueue()
	}
}

func (m *managerImpl) flushIndicatorQueue() {
	// This is a potentially long-running operation, and we don't want to have a pile of goroutines queueing up on
	// this lock.
	if !m.flushProcessingLock.MaybeLock() {
		return
	}
	defer m.flushProcessingLock.Unlock()

	copiedQueue := m.copyAndResetIndicatorQueue()
	if len(copiedQueue) == 0 {
		return
	}

	// Map copiedQueue to slice
	indicatorSlice := make([]*storage.ProcessIndicator, 0, len(copiedQueue))
	for _, i := range copiedQueue {
		indicatorSlice = append(indicatorSlice, i.indicator)
	}

	// Index the process indicators in batch
	if err := m.processesDataStore.AddProcessIndicators(lifecycleMgrCtx, indicatorSlice...); err != nil {
		log.Errorf("Error adding process indicators: %v", err)
	}

	deploymentIDs := uniqueDeploymentIDs(copiedQueue)
	newAlerts, err := m.runtimeDetector.AlertsForDeployments(deploymentIDs...)
	if err != nil {
		log.Errorf("Failed to compute runtime alerts: %s", err)
		return
	}

	deploymentToMatchingIndicators := make(map[string][]*storage.ProcessIndicator)
	for _, ind := range indicatorSlice {
		userWhitelist, _, err := m.checkWhitelist(ind)
		if err != nil {
			log.Errorf("error checking whitelist for indicator: %v - %+v", err, ind)
			continue
		}
		if userWhitelist {
			deploymentToMatchingIndicators[ind.GetDeploymentId()] = append(deploymentToMatchingIndicators[ind.GetDeploymentId()], ind)
		}
	}
	if len(deploymentToMatchingIndicators) > 0 {
		// Compute whitelist alerts here
		whitelistAlerts, err := m.getWhitelistAlerts(deploymentToMatchingIndicators)
		if err != nil {
			log.Errorf("failed to get whitelist alerts: %v", err)
			return
		}
		newAlerts = append(newAlerts, whitelistAlerts...)
	}

	modified, err := m.alertManager.AlertAndNotify(lifecycleMgrCtx, newAlerts, alertmanager.WithLifecycleStage(storage.LifecycleStage_RUNTIME), alertmanager.WithDeploymentIDs(deploymentIDs...))
	if err != nil {
		log.Errorf("Couldn't alert and notify: %s", err)
	} else if modified {
		defer m.reprocessor.ReprocessRiskForDeployments(deploymentIDs...)
	}

	// Create enforcement actions for the new alerts and send them with the stored injectors.
	m.generateAndSendEnforcements(newAlerts, copiedQueue)
}

func (m *managerImpl) addToQueue(indicator *storage.ProcessIndicator, injector common.MessageInjector) {
	m.queueLock.Lock()
	defer m.queueLock.Unlock()

	m.queuedIndicators[indicator.GetId()] = indicatorWithInjector{
		indicator:           indicator,
		msgToSensorInjector: injector,
	}
}

func (m *managerImpl) getWhitelistAlerts(deploymentsToIndicators map[string][]*storage.ProcessIndicator) ([]*storage.Alert, error) {
	whitelistExecutor := newWhitelistExecutor(lifecycleMgrCtx, m.deploymentDataStore, deploymentsToIndicators)
	if err := m.runtimeDetector.PolicySet().ForEach(whitelistExecutor); err != nil {
		return nil, err
	}
	return whitelistExecutor.alerts, nil
}

const (
	whitelistLockingGracePeriod = 5 * time.Second
)

func (m *managerImpl) checkWhitelist(indicator *storage.ProcessIndicator) (userWhitelist bool, roxWhitelist bool, err error) {
	// Always reprocess risk for the deployment, since that's needed to update its process-related information.
	defer m.reprocessor.ReprocessRiskForDeployments(indicator.GetDeploymentId())

	key := &storage.ProcessWhitelistKey{
		DeploymentId:  indicator.DeploymentId,
		ContainerName: indicator.ContainerName,
		ClusterId:     indicator.GetClusterId(),
		Namespace:     indicator.GetNamespace(),
	}

	// TODO joseph what to do if whitelist doesn't exist?  Always create for now?
	whitelist, err := m.whitelists.GetProcessWhitelist(lifecycleMgrCtx, key)
	if err != nil {
		return
	}

	insertableElement := &storage.WhitelistItem{Item: &storage.WhitelistItem_ProcessName{ProcessName: processwhitelist.WhitelistItemFromProcess(indicator)}}
	if whitelist == nil {
		_, err = m.whitelists.UpsertProcessWhitelist(lifecycleMgrCtx, key, []*storage.WhitelistItem{insertableElement}, true)
		if err == nil {
			// This updates the risk for deployments after the whitelist gets locked.
			// This isn't super pretty, but otherwise, our whitelistStatus for the deployment can become stale
			// since we don't explicit lock a whitelist -- we just set a locked time which is in the future.
			go func() {
				time.Sleep(env.WhitelistGenerationDuration.DurationSetting() + whitelistLockingGracePeriod)
				m.reprocessor.ReprocessRiskForDeployments(indicator.GetDeploymentId())
			}()
		}
		return
	}

	for _, element := range whitelist.GetElements() {
		if element.GetElement().GetProcessName() == insertableElement.GetProcessName() {
			return
		}
	}
	userWhitelist = processwhitelist.IsUserLocked(whitelist)
	roxWhitelist = processwhitelist.IsRoxLocked(whitelist)
	if userWhitelist || roxWhitelist {
		return
	}
	_, err = m.whitelists.UpdateProcessWhitelistElements(lifecycleMgrCtx, key, []*storage.WhitelistItem{insertableElement}, nil, true)
	if err != nil {
		return
	}
	return
}

func (m *managerImpl) IndicatorAdded(indicator *storage.ProcessIndicator, injector common.MessageInjector) error {
	if indicator.GetId() == "" {
		return fmt.Errorf("invalid indicator received: %s, id was empty", proto.MarshalTextString(indicator))
	}

	m.addToQueue(indicator, injector)

	if m.limiter.Allow() {
		go m.flushIndicatorQueue()
	}
	return nil
}

func (m *managerImpl) DeploymentUpdated(deployment *storage.Deployment) (string, string, storage.EnforcementAction, error) {
	// Attempt to enrich the image before detection.
	images, updatedIndices, err := m.enricher.EnrichDeployment(enricher.EnrichmentContext{NoExternalMetadata: false}, deployment)
	if err != nil {
		log.Errorf("Error enriching deployment %s: %s", deployment.GetName(), err)
	}
	if len(updatedIndices) > 0 {
		for _, idx := range updatedIndices {
			img := images[idx]
			if err := m.imageDataStore.UpsertImage(lifecycleMgrCtx, img); err != nil {
				log.Errorf("Error persisting image %s: %s", img.GetName().GetFullName(), err)
			}
		}
	}

	// Update risk after processing and save the deployment.
	// There is no need to save the deployment in this function as it will be saved post reprocessing risk
	defer m.riskManager.ReprocessDeploymentRiskWithImages(deployment, images)

	presentAlerts, err := m.deploytimeDetector.Detect(deploytime.DetectionContext{}, deployment, images)
	if err != nil {
		return "", "", storage.EnforcementAction_UNSET_ENFORCEMENT, errors.Wrap(err, "fetching deploy time alerts")
	}

	if _, err := m.alertManager.AlertAndNotify(lifecycleMgrCtx, presentAlerts,
		alertmanager.WithLifecycleStage(storage.LifecycleStage_DEPLOY), alertmanager.WithDeploymentIDs(deployment.GetId())); err != nil {
		return "", "", storage.EnforcementAction_UNSET_ENFORCEMENT, err
	}

	// Generate enforcement actions based on the currently generated alerts.
	alertToEnforce, policyName, enforcementAction := determineEnforcement(presentAlerts)
	return alertToEnforce, policyName, enforcementAction, nil
}

func (m *managerImpl) UpsertPolicy(policy *storage.Policy) error {
	// Asynchronously update all deployments' risk after processing.
	defer m.reprocessor.ReprocessRisk()

	var presentAlerts []*storage.Alert

	// Add policy to set.
	if policies.AppliesAtDeployTime(policy) {
		if err := m.deploytimeDetector.PolicySet().UpsertPolicy(policy); err != nil {
			return errors.Wrapf(err, "adding policy %s to deploy time detector", policy.GetName())
		}
		deployTimeAlerts, err := m.deploytimeDetector.AlertsForPolicy(policy.GetId())
		if err != nil {
			return errors.Wrapf(err, "error generating deploy-time alerts for policy %s", policy.GetName())
		}
		presentAlerts = append(presentAlerts, deployTimeAlerts...)
	} else {
		err := m.deploytimeDetector.PolicySet().RemovePolicy(policy.GetId())
		if err != nil {
			return errors.Wrapf(err, "removing policy %s from deploy time detector", policy.GetName())
		}
	}

	if policies.AppliesAtRunTime(policy) {
		if err := m.runtimeDetector.PolicySet().UpsertPolicy(policy); err != nil {
			return errors.Wrapf(err, "adding policy %s to runtime detector", policy.GetName())
		}
		runTimeAlerts, err := m.runtimeDetector.AlertsForPolicy(policy.GetId())
		if err != nil {
			return errors.Wrapf(err, "error generating runtime alerts for policy %s", policy.GetName())
		}
		presentAlerts = append(presentAlerts, runTimeAlerts...)
	} else {
		err := m.runtimeDetector.PolicySet().RemovePolicy(policy.GetId())
		if err != nil {
			return errors.Wrapf(err, "removing policy %s from runtime detector", policy.GetName())
		}
	}

	// Perform notifications and update DB.
	_, err := m.alertManager.AlertAndNotify(lifecycleMgrCtx, presentAlerts, alertmanager.WithPolicyID(policy.GetId()))
	return err
}

func (m *managerImpl) RecompilePolicy(policy *storage.Policy) error {
	// Asynchronously update all deployments' risk after processing.
	defer m.reprocessor.ReprocessRisk()

	var presentAlerts []*storage.Alert

	if policies.AppliesAtDeployTime(policy) {
		if err := m.deploytimeDetector.PolicySet().Recompile(policy.GetId()); err != nil {
			return errors.Wrapf(err, "adding policy %s to deploy time detector", policy.GetName())
		}
		deployTimeAlerts, err := m.deploytimeDetector.AlertsForPolicy(policy.GetId())
		if err != nil {
			return errors.Wrapf(err, "error generating deploy-time alerts for policy %s", policy.GetName())
		}
		presentAlerts = append(presentAlerts, deployTimeAlerts...)
	} else {
		err := m.deploytimeDetector.PolicySet().RemovePolicy(policy.GetId())
		if err != nil {
			return errors.Wrapf(err, "removing policy %s from deploy time detector", policy.GetName())
		}
	}

	if policies.AppliesAtRunTime(policy) {
		if err := m.runtimeDetector.PolicySet().Recompile(policy.GetId()); err != nil {
			return errors.Wrapf(err, "adding policy %s to runtime detector", policy.GetName())
		}
		runTimeAlerts, err := m.runtimeDetector.AlertsForPolicy(policy.GetId())
		if err != nil {
			return errors.Wrapf(err, "error generating runtime alerts for policy %s", policy.GetName())
		}
		presentAlerts = append(presentAlerts, runTimeAlerts...)
	} else {
		err := m.runtimeDetector.PolicySet().RemovePolicy(policy.GetId())
		if err != nil {
			return errors.Wrapf(err, "removing policy %s from runtime detector", policy.GetName())
		}
	}

	// Perform notifications and update DB.
	_, err := m.alertManager.AlertAndNotify(lifecycleMgrCtx, presentAlerts, alertmanager.WithPolicyID(policy.GetId()))
	return err
}

func (m *managerImpl) DeploymentRemoved(deployment *storage.Deployment) error {
	_, err := m.alertManager.AlertAndNotify(lifecycleMgrCtx, nil, alertmanager.WithDeploymentIDs(deployment.GetId()))
	return err
}

func (m *managerImpl) RemovePolicy(policyID string) error {
	// Asynchronously update all deployments' risk after processing.
	defer m.reprocessor.ReprocessRisk()

	if err := m.deploytimeDetector.PolicySet().RemovePolicy(policyID); err != nil {
		return err
	}
	if err := m.runtimeDetector.PolicySet().RemovePolicy(policyID); err != nil {
		return err
	}
	_, err := m.alertManager.AlertAndNotify(lifecycleMgrCtx, nil, alertmanager.WithPolicyID(policyID))
	return err
}

// determineEnforcement returns the alert and its enforcement action to use from the input list (if any have enforcement).
func determineEnforcement(alerts []*storage.Alert) (alertID, policyName string, action storage.EnforcementAction) {
	for _, alert := range alerts {
		if alert.GetEnforcement().GetAction() == storage.EnforcementAction_SCALE_TO_ZERO_ENFORCEMENT {
			return alert.GetId(), alert.GetPolicy().GetName(), storage.EnforcementAction_SCALE_TO_ZERO_ENFORCEMENT
		}

		if alert.GetEnforcement().GetAction() != storage.EnforcementAction_UNSET_ENFORCEMENT {
			alertID = alert.GetId()
			policyName = alert.GetPolicy().GetName()
			action = alert.GetEnforcement().GetAction()
		}
	}
	return
}

func uniqueDeploymentIDs(indicatorsToInfo map[string]indicatorWithInjector) []string {
	m := set.NewStringSet()
	for _, infoWithInjector := range indicatorsToInfo {
		deploymentID := infoWithInjector.indicator.GetDeploymentId()
		if deploymentID == "" {
			continue
		}
		m.Add(deploymentID)
	}
	return m.AsSlice()
}

func (m *managerImpl) generateAndSendEnforcements(alerts []*storage.Alert, indicatorsToInfo map[string]indicatorWithInjector) {
	for _, alert := range alerts {
		// Skip alerts without runtime enforcement.
		if alert.GetEnforcement().GetAction() != storage.EnforcementAction_KILL_POD_ENFORCEMENT {
			continue
		}

		// If the alert has enforcement, we want to generate a list of enforcement and injector pairs.
		for _, singleIndicator := range alert.GetProcessViolation().GetProcesses() {
			if infoWithInjector, ok := indicatorsToInfo[singleIndicator.GetId()]; ok {
				// Get the deployment details to check that the deployment is still running.
				deployment, exists, err := m.deploymentDataStore.GetDeployment(lifecycleMgrCtx, alert.GetDeployment().GetId())
				if err != nil {
					log.Errorf("Couldn't enforce on deployment %s: failed to retrieve: %s", alert.GetDeployment().GetId(), err)
					continue
				}
				if !exists {
					log.Errorf("Couldn't enforce on deployment %s: not found in store", alert.GetDeployment().GetId())
					continue
				}

				// Generate the enforcement action.
				containerID := infoWithInjector.indicator.GetSignal().GetContainerId()
				enforcement := createEnforcementAction(alert, deployment, containerID)
				if enforcement == nil {
					log.Errorf("Couldn't enforce on container %s, not found in deployment %s/%s", containerID, deployment.GetNamespace(), deployment.GetName())
					continue
				}

				// Attempt to send the enforcement with the injector.
				err = infoWithInjector.msgToSensorInjector.InjectMessage(context.Background(), &central.MsgToSensor{
					Msg: &central.MsgToSensor_Enforcement{
						Enforcement: enforcement,
					},
				})
				if err != nil {
					log.Errorf("Failed to inject enforcement action %s: %v", proto.MarshalTextString(enforcement), err)
				}
			}
		}
	}
}

func createEnforcementAction(alert *storage.Alert, deployment *storage.Deployment, containerID string) *central.SensorEnforcement {
	// Find the container instance in the deployment, and generate an enforcement action for it.
	for _, container := range deployment.GetContainers() {
		for _, instance := range container.GetInstances() {
			// Skip instances without well formed IDs.
			if len(instance.GetInstanceId().GetId()) < 12 {
				continue
			}

			// If the id matches the container we are generating enforcement from, create an eforcement action for the instance.
			if containerID == instance.GetInstanceId().GetId()[:12] {
				return createEnforcementActionFromInstance(alert, deployment, instance)
			}
		}
	}
	return nil
}

func createEnforcementActionFromInstance(alert *storage.Alert, deployment *storage.Deployment, instance *storage.ContainerInstance) *central.SensorEnforcement {
	resource := &central.SensorEnforcement_ContainerInstance{
		ContainerInstance: &central.ContainerInstanceEnforcement{
			ContainerInstanceId: instance.GetInstanceId().GetId(),
			PodId:               instance.GetContainingPodId(),
			DeploymentEnforcement: &central.DeploymentEnforcement{
				DeploymentId:   deployment.GetId(),
				DeploymentName: deployment.GetName(),
				Namespace:      deployment.GetNamespace(),
				DeploymentType: deployment.GetType(),
				PolicyName:     alert.GetPolicy().GetName(),
			},
		},
	}
	return &central.SensorEnforcement{
		Enforcement: storage.EnforcementAction_KILL_POD_ENFORCEMENT,
		Resource:    resource,
	}
}
