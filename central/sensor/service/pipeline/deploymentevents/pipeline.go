package deploymentevents

import (
	"context"

	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	clusterDataStore "github.com/stackrox/rox/central/cluster/datastore"
	deploymentDataStore "github.com/stackrox/rox/central/deployment/datastore"
	"github.com/stackrox/rox/central/detection/lifecycle"
	imageDataStore "github.com/stackrox/rox/central/image/datastore"
	countMetrics "github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/central/networkpolicies/graph"
	"github.com/stackrox/rox/central/sensor/service/common"
	"github.com/stackrox/rox/central/sensor/service/pipeline"
	"github.com/stackrox/rox/central/sensor/service/pipeline/reconciliation"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/images/enricher"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/search"
)

var (
	log = logging.LoggerForModule()
)

// Template design pattern. We define control flow here and defer logic to subclasses.
//////////////////////////////////////////////////////////////////////////////////////

// GetPipeline returns an instantiation of this particular pipeline
func GetPipeline() pipeline.Fragment {
	return NewPipeline(clusterDataStore.Singleton(),
		deploymentDataStore.Singleton(),
		imageDataStore.Singleton(),
		lifecycle.SingletonManager(),
		graph.Singleton())
}

// NewPipeline returns a new instance of Pipeline.
func NewPipeline(clusters clusterDataStore.DataStore, deployments deploymentDataStore.DataStore,
	images imageDataStore.DataStore, manager lifecycle.Manager,
	graphEvaluator graph.Evaluator) pipeline.Fragment {
	return &pipelineImpl{
		validateInput:     newValidateInput(),
		clusterEnrichment: newClusterEnrichment(clusters),
		updateImages:      newUpdateImages(images),
		persistDeployment: newPersistDeployment(deployments),
		lifecycleManager:  manager,

		graphEvaluator: graphEvaluator,
		deployments:    deployments,
		clusters:       clusters,
	}
}

type pipelineImpl struct {
	// pipeline stages.
	validateInput     *validateInputImpl
	clusterEnrichment *clusterEnrichmentImpl
	updateImages      *updateImagesImpl
	persistDeployment *persistDeploymentImpl
	lifecycleManager  lifecycle.Manager

	deployments deploymentDataStore.DataStore
	clusters    clusterDataStore.DataStore

	graphEvaluator graph.Evaluator
}

func (s *pipelineImpl) Reconcile(ctx context.Context, clusterID string, storeMap *reconciliation.StoreMap) error {
	query := search.NewQueryBuilder().AddExactMatches(search.ClusterID, clusterID).ProtoQuery()
	results, err := s.deployments.Search(ctx, query)
	if err != nil {
		return err
	}

	store := storeMap.Get((*central.SensorEvent_Deployment)(nil))
	return reconciliation.Perform(store, search.ResultsToIDSet(results), "deployments", func(id string) error {
		return s.runRemovePipeline(ctx, central.ResourceAction_REMOVE_RESOURCE, &storage.Deployment{Id: id})
	})
}

func (s *pipelineImpl) Match(msg *central.MsgFromSensor) bool {
	return msg.GetEvent().GetDeployment() != nil
}

// Run runs the pipeline template on the input and returns the output.
func (s *pipelineImpl) Run(ctx context.Context, clusterID string, msg *central.MsgFromSensor, injector common.MessageInjector) error {
	defer countMetrics.IncrementResourceProcessedCounter(pipeline.ActionToOperation(msg.GetEvent().GetAction()), metrics.Deployment)

	event := msg.GetEvent()
	deployment := event.GetDeployment()
	deployment.ClusterId = clusterID

	var err error
	switch event.GetAction() {
	case central.ResourceAction_REMOVE_RESOURCE:
		err = s.runRemovePipeline(ctx, event.GetAction(), deployment)
	default:
		err = s.runGeneralPipeline(ctx, event.GetAction(), deployment, injector)
	}
	return err
}

// Run runs the pipeline template on the input and returns the output.
func (s *pipelineImpl) runRemovePipeline(ctx context.Context, action central.ResourceAction, deployment *storage.Deployment) error {
	// Validate the the deployment we receive has necessary fields set.
	if err := s.validateInput.do(deployment); err != nil {
		return err
	}

	// Add/Update/Remove the deployment from persistence depending on the deployment action.
	if err := s.persistDeployment.do(ctx, action, deployment, false); err != nil {
		return err
	}

	s.graphEvaluator.IncrementEpoch()

	// Process the deployment (enrichment, alert generation, enforcement action generation.)
	if err := s.lifecycleManager.DeploymentRemoved(deployment); err != nil {
		return err
	}
	return nil
}

func (s *pipelineImpl) dedupeBasedOnHash(ctx context.Context, action central.ResourceAction, newDeployment *storage.Deployment) (bool, error) {
	var err error
	newDeployment.Hash, err = hashstructure.Hash(newDeployment, &hashstructure.HashOptions{})
	if err != nil {
		return false, err
	}

	// Check if this deployment needs to be processed based on hash
	oldDeployment, exists, err := s.deployments.GetDeployment(ctx, newDeployment.GetId())
	if err != nil || !exists {
		return false, err
	}
	// If it already exists and the hash is the same, then just update the container instances of the old deployment and upsert
	if exists && oldDeployment.GetHash() == newDeployment.GetHash() {
		// Using the index of Container is save as this is ensured by the hash
		for i, c := range newDeployment.GetContainers() {
			oldDeployment.Containers[i].Instances = c.Instances
		}
		return true, s.persistDeployment.do(ctx, action, oldDeployment, false)
	}
	return false, nil
}

// Run runs the pipeline template on the input and returns the output.
func (s *pipelineImpl) runGeneralPipeline(ctx context.Context, action central.ResourceAction, deployment *storage.Deployment, injector common.MessageInjector) error {
	// Validate the the deployment we receive has necessary fields set.
	if err := s.validateInput.do(deployment); err != nil {
		return err
	}

	dedupe, err := s.dedupeBasedOnHash(ctx, action, deployment)
	if err != nil {
		err = errors.Wrapf(err, "Could not check deployment %q for deduping", deployment.GetName())
		log.Error(err)
		return err
	}
	if dedupe {
		return nil
	}

	// Fill in cluster information.
	if err := s.clusterEnrichment.do(ctx, deployment); err != nil {
		log.Errorf("Couldn't get cluster identity: %s", err)
		return err
	}

	// Add/Update/Remove the deployment from persistence depending on the deployment action.
	if err := s.persistDeployment.do(ctx, action, deployment, true); err != nil {
		return err
	}

	// Update the deployments images with the latest version from storage.
	s.updateImages.do(ctx, deployment)

	// Process the deployment (alert generation, enforcement action generation)
	// Only pass in the enforcement injector (which is a signal to the lifecycle manager
	// that it should inject enforcement) on creates.
	var injectorToPass common.MessageInjector
	if action == central.ResourceAction_CREATE_RESOURCE {
		injectorToPass = injector
	}
	if err := s.lifecycleManager.DeploymentUpdated(enricher.EnrichmentContext{UseNonBlockingCallsWherePossible: true}, deployment, injectorToPass); err != nil {
		return err
	}
	s.graphEvaluator.IncrementEpoch()
	return nil
}

func (s *pipelineImpl) OnFinish(_ string) {}
