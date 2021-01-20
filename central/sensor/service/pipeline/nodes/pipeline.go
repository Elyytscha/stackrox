package nodes

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	clusterDataStore "github.com/stackrox/rox/central/cluster/datastore"
	"github.com/stackrox/rox/central/enrichment"
	countMetrics "github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/central/node/globaldatastore"
	"github.com/stackrox/rox/central/node/store"
	"github.com/stackrox/rox/central/sensor/service/common"
	"github.com/stackrox/rox/central/sensor/service/pipeline"
	"github.com/stackrox/rox/central/sensor/service/pipeline/reconciliation"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/nodes/enricher"
	"github.com/stackrox/rox/pkg/search"
)

var (
	log = logging.LoggerForModule()
)

// Template design pattern. We define control flow here and defer logic to subclasses.
//////////////////////////////////////////////////////////////////////////////////////

// GetPipeline returns an instantiation of this particular pipeline
func GetPipeline() pipeline.Fragment {
	return NewPipeline(clusterDataStore.Singleton(), globaldatastore.Singleton(), enrichment.NodeEnricherSingleton())
}

// NewPipeline returns a new instance of Pipeline.
func NewPipeline(clusters clusterDataStore.DataStore, nodes globaldatastore.GlobalDataStore, enricher enricher.NodeEnricher) pipeline.Fragment {
	return &pipelineImpl{
		clusterStore: clusters,
		nodeStore:    nodes,
		enricher:     enricher,
	}
}

type pipelineImpl struct {
	clusterStore clusterDataStore.DataStore
	nodeStore    globaldatastore.GlobalDataStore
	enricher     enricher.NodeEnricher
}

func (p *pipelineImpl) Reconcile(ctx context.Context, clusterID string, storeMap *reconciliation.StoreMap) error {
	query := search.NewQueryBuilder().AddExactMatches(search.ClusterID, clusterID).ProtoQuery()
	results, err := p.nodeStore.Search(ctx, query)
	if err != nil {
		return err
	}

	clusterStore, err := p.nodeStore.GetClusterNodeStore(ctx, clusterID, true)
	if err != nil {
		return errors.Wrap(err, "getting cluster-local node store")
	}

	store := storeMap.Get((*central.SensorEvent_Node)(nil))
	return reconciliation.Perform(store, search.ResultsToIDSet(results), "nodes", func(id string) error {
		return p.processRemove(clusterStore, &storage.Node{Id: id})
	})
}

func (p *pipelineImpl) Match(msg *central.MsgFromSensor) bool {
	return msg.GetEvent().GetNode() != nil
}

func (p *pipelineImpl) processRemove(store store.Store, n *storage.Node) error {
	return store.RemoveNode(n.GetId())
}

// Run runs the pipeline template on the input and returns the output.
func (p *pipelineImpl) Run(ctx context.Context, clusterID string, msg *central.MsgFromSensor, _ common.MessageInjector) error {
	defer countMetrics.IncrementResourceProcessedCounter(pipeline.ActionToOperation(msg.GetEvent().GetAction()), metrics.Node)

	event := msg.GetEvent()

	store, err := p.nodeStore.GetClusterNodeStore(ctx, clusterID, true)
	if err != nil {
		return errors.Wrap(err, "getting cluster-local node store")
	}

	node := event.GetNode()
	if node == nil {
		return fmt.Errorf("unexpected resource type %T for cluster status", event.Resource)
	}

	if event.GetAction() == central.ResourceAction_REMOVE_RESOURCE {
		return p.processRemove(store, node)
	}

	node = node.Clone()
	node.ClusterId = clusterID
	clusterName, ok, err := p.clusterStore.GetClusterName(ctx, clusterID)
	if err == nil && ok {
		node.ClusterName = clusterName
	}

	if features.HostScanning.Enabled() {
		enrichmentCtx := enricher.EnrichmentContext{
			FetchOpt: enricher.UseCachesIfPossible,
		}
		if err := p.enricher.EnrichNode(enrichmentCtx, node); err != nil {
			log.Warn(err.Error())
		}
	}
	return store.UpsertNode(node)
}

func (p *pipelineImpl) OnFinish(_ string) {}
