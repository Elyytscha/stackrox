package all

import (
	"context"

	"github.com/stackrox/rox/central/sensor/service/pipeline"
	"github.com/stackrox/rox/central/sensor/service/pipeline/alerts"
	"github.com/stackrox/rox/central/sensor/service/pipeline/clusterhealthupdate"
	"github.com/stackrox/rox/central/sensor/service/pipeline/clusterstatusupdate"
	"github.com/stackrox/rox/central/sensor/service/pipeline/deploymentevents"
	"github.com/stackrox/rox/central/sensor/service/pipeline/imageintegrations"
	"github.com/stackrox/rox/central/sensor/service/pipeline/namespaces"
	"github.com/stackrox/rox/central/sensor/service/pipeline/networkflowupdate"
	"github.com/stackrox/rox/central/sensor/service/pipeline/networkpolicies"
	"github.com/stackrox/rox/central/sensor/service/pipeline/nodes"
	"github.com/stackrox/rox/central/sensor/service/pipeline/podevents"
	"github.com/stackrox/rox/central/sensor/service/pipeline/processindicators"
	"github.com/stackrox/rox/central/sensor/service/pipeline/reprocessing"
	"github.com/stackrox/rox/central/sensor/service/pipeline/rolebindings"
	"github.com/stackrox/rox/central/sensor/service/pipeline/roles"
	"github.com/stackrox/rox/central/sensor/service/pipeline/secrets"
	"github.com/stackrox/rox/central/sensor/service/pipeline/serviceaccounts"
)

// NewFactory returns a new instance of a Factory that produces a pipeline handling all message types.
func NewFactory() pipeline.Factory {
	return &factoryImpl{}
}

type factoryImpl struct{}

// sendMessages grabs items from the queue, processes them, and sends them back to sensor.
func (s *factoryImpl) PipelineForCluster(ctx context.Context, clusterID string) (pipeline.ClusterPipeline, error) {
	flowUpdateFragment, err := networkflowupdate.Singleton().GetFragment(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	pipelines := []pipeline.Fragment{
		deploymentevents.GetPipeline(),
		podevents.GetPipeline(),
		processindicators.GetPipeline(),
		networkpolicies.GetPipeline(),
		namespaces.GetPipeline(),
		secrets.GetPipeline(),
		nodes.GetPipeline(),
		flowUpdateFragment,
		imageintegrations.GetPipeline(),
		clusterstatusupdate.GetPipeline(),
		clusterhealthupdate.GetPipeline(),
		serviceaccounts.GetPipeline(),
		roles.GetPipeline(),
		rolebindings.GetPipeline(),
		reprocessing.GetPipeline(),
		alerts.GetPipeline(),
	}

	return NewClusterPipeline(clusterID, pipelines...), nil
}
