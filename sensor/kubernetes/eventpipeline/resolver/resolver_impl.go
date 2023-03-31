package resolver

import (
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/sensor/common/store"
	"github.com/stackrox/rox/sensor/kubernetes/eventpipeline/component"
	"github.com/stackrox/rox/sensor/kubernetes/eventpipeline/tester"
)

var (
	log = logging.LoggerForModule()
)

type resolverImpl struct {
	outputQueue component.OutputQueue
	innerQueue  chan *component.ResourceEvent

	storeProvider store.Provider
}

// Start the resolverImpl component
func (r *resolverImpl) Start() error {
	go r.runResolver()
	return nil
}

// Stop the resolverImpl component
func (r *resolverImpl) Stop(_ error) {
	defer close(r.innerQueue)
}

// Send a ResourceEvent message to the inner queue
func (r *resolverImpl) Send(event *component.ResourceEvent) {
	r.innerQueue <- event
}

// runResolver reads messages from the inner queue and process the message
func (r *resolverImpl) runResolver() {
	for {
		msg, more := <-r.innerQueue
		if !more {
			return
		}
		r.processMessage(msg)
	}
}

// processMessage resolves the dependencies and forwards the message to the outputQueue
func (r *resolverImpl) processMessage(msg *component.ResourceEvent) {
	if msg.DeploymentReferences != nil {

		for _, deploymentReference := range msg.DeploymentReferences {
			if deploymentReference.Reference == nil {
				continue
			}

			referenceIds := deploymentReference.Reference(r.storeProvider.Deployments())

			if deploymentReference.ForceDetection && len(referenceIds) > 0 {
				// We append the referenceIds to the msg to be reprocessed
				msg.AddDeploymentForReprocessing(referenceIds...)
			}

			for _, id := range referenceIds {
				preBuiltDeployment := r.storeProvider.Deployments().Get(id)
				if preBuiltDeployment == nil {
					log.Warnf("Deployment with id %s not found", id)
					continue
				}
				// Skip resolving the deployment dependencies
				if deploymentReference.SkipResolving {
					d := r.storeProvider.Deployments().GetBuiltDeployment(id)
					if d == nil {
						log.Warnf("Deployment with id %s not found", id)
						continue
					}
					msg.AddDeploymentForDetection(component.DetectorMessage{Object: d, Action: deploymentReference.ParentResourceAction})
					continue
				}

				// Remove actions are done at the handler level. This is not ideal but for now it allows us to be able to fetch deployments from the store
				// in the resolver instead of sending a copy. We still manage OnDeploymentCreateOrUpdate here.
				r.storeProvider.EndpointManager().OnDeploymentCreateOrUpdateByID(id)

				permissionLevel := r.storeProvider.RBAC().GetPermissionLevelForDeployment(preBuiltDeployment)
				exposureInfo := r.storeProvider.Services().
					GetExposureInfos(preBuiltDeployment.GetNamespace(), preBuiltDeployment.GetPodLabels())

				d, err := r.storeProvider.Deployments().BuildDeploymentWithDependencies(id, store.Dependencies{
					PermissionLevel: permissionLevel,
					Exposures:       exposureInfo,
				})

				if err != nil {
					log.Warnf("Failed to build deployment dependency: %s", err)
					continue
				}
				msg.AddSensorEvent(toEvent(deploymentReference.ParentResourceAction, d, msg.DeploymentTiming)).
					AddDeploymentForDetection(component.DetectorMessage{Object: d, Action: deploymentReference.ParentResourceAction})
				msg.UpdateMsgToTester(toEvent(deploymentReference.ParentResourceAction, d, msg.DeploymentTiming))
			}
		}
	}
	if env.ResyncTester.BooleanSetting() {
		for _, toTester := range msg.TesterMsg {
			tester.GetEventPipelineTester().Send(toTester)
		}
	}

	r.outputQueue.Send(msg)
}

func toEvent(action central.ResourceAction, deployment *storage.Deployment, timing *central.Timing) *central.SensorEvent {
	return &central.SensorEvent{
		Id:     deployment.GetId(),
		Action: action,
		Timing: timing,
		Resource: &central.SensorEvent_Deployment{
			Deployment: deployment.Clone(),
		},
	}
}

var _ component.Resolver = (*resolverImpl)(nil)
