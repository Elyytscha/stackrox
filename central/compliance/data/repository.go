package data

import (
	"context"
	"math"

	"github.com/stackrox/rox/central/compliance/framework"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/internalapi/compliance"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
)

var (
	log = logging.LoggerForModule()
)

type repository struct {
	cluster     *storage.Cluster
	nodes       map[string]*storage.Node
	deployments map[string]*storage.Deployment

	unresolvedAlerts      []*storage.ListAlert
	networkPolicies       map[string]*storage.NetworkPolicy
	networkGraph          *v1.NetworkGraph
	policies              map[string]*storage.Policy
	images                []*storage.ListImage
	imageIntegrations     []*storage.ImageIntegration
	registries            []framework.ImageMatcher
	scanners              []framework.ImageMatcher
	processIndicators     []*storage.ProcessIndicator
	networkFlows          []*storage.NetworkFlow
	notifiers             []*storage.Notifier
	roles                 []*storage.K8SRole
	bindings              []*storage.K8SRoleBinding
	cisDockerRunCheck     bool
	cisKubernetesRunCheck bool
	categoryToPolicies    map[string]set.StringSet // maps categories to policy set

	hostScrape map[string]*compliance.ComplianceReturn
}

func (r *repository) Cluster() *storage.Cluster {
	return r.cluster
}

func (r *repository) Nodes() map[string]*storage.Node {
	return r.nodes
}

func (r *repository) Deployments() map[string]*storage.Deployment {
	return r.deployments
}

func (r *repository) NetworkPolicies() map[string]*storage.NetworkPolicy {
	return r.networkPolicies
}

func (r *repository) NetworkGraph() *v1.NetworkGraph {
	return r.networkGraph
}

func (r *repository) Policies() map[string]*storage.Policy {
	return r.policies
}

func (r *repository) PolicyCategories() map[string]set.StringSet {
	return r.categoryToPolicies
}

func (r *repository) Images() []*storage.ListImage {
	return r.images
}

func (r *repository) ImageIntegrations() []*storage.ImageIntegration {
	return r.imageIntegrations
}

func (r *repository) ProcessIndicators() []*storage.ProcessIndicator {
	return r.processIndicators
}

func (r *repository) NetworkFlows() []*storage.NetworkFlow {
	return r.networkFlows
}

func (r *repository) Notifiers() []*storage.Notifier {
	return r.notifiers
}

func (r *repository) K8sRoles() []*storage.K8SRole {
	return r.roles
}

func (r *repository) K8sRoleBindings() []*storage.K8SRoleBinding {
	return r.bindings
}

func (r *repository) UnresolvedAlerts() []*storage.ListAlert {
	return r.unresolvedAlerts
}

func (r *repository) HostScraped(node *storage.Node) *compliance.ComplianceReturn {
	return r.hostScrape[node.GetName()]
}

func (r *repository) CISDockerTriggered() bool {
	return r.cisDockerRunCheck
}

func (r *repository) CISKubernetesTriggered() bool {
	return r.cisKubernetesRunCheck
}

func (r *repository) RegistryIntegrations() []framework.ImageMatcher {
	return r.registries
}

func (r *repository) ScannerIntegrations() []framework.ImageMatcher {
	return r.scanners
}

func newRepository(ctx context.Context, domain framework.ComplianceDomain, scrapeResults map[string]*compliance.ComplianceReturn, factory *factory) (*repository, error) {
	r := &repository{}
	if err := r.init(ctx, domain, scrapeResults, factory); err != nil {
		return nil, err
	}
	return r, nil
}

func nodesByID(nodes []*storage.Node) map[string]*storage.Node {
	result := make(map[string]*storage.Node, len(nodes))
	for _, node := range nodes {
		result[node.GetId()] = node
	}
	return result
}

func deploymentsByID(deployments []*storage.Deployment) map[string]*storage.Deployment {
	result := make(map[string]*storage.Deployment, len(deployments))
	for _, deployment := range deployments {
		result[deployment.GetId()] = deployment
	}
	return result
}

func networkPoliciesByID(policies []*storage.NetworkPolicy) map[string]*storage.NetworkPolicy {
	result := make(map[string]*storage.NetworkPolicy, len(policies))
	for _, policy := range policies {
		result[policy.GetId()] = policy
	}
	return result
}

func policiesByName(policies []*storage.Policy) map[string]*storage.Policy {
	result := make(map[string]*storage.Policy, len(policies))
	for _, policy := range policies {
		result[policy.GetName()] = policy
	}
	return result
}

func policyCategories(policies []*storage.Policy) map[string]set.StringSet {
	result := make(map[string]set.StringSet, len(policies))
	for _, policy := range policies {
		if policy.Disabled {
			continue
		}
		for _, category := range policy.Categories {
			policySet, ok := result[category]
			if !ok {
				policySet = set.NewStringSet()
			}
			policySet.Add(policy.Name)
			result[category] = policySet
		}
	}
	return result
}

func expandFile(parent *compliance.File) map[string]*compliance.File {
	expanded := make(map[string]*compliance.File)
	for _, child := range parent.GetChildren() {
		childExpanded := expandFile(child)
		for k, v := range childExpanded {
			expanded[k] = v
		}
	}
	expanded[parent.GetPath()] = parent
	return expanded
}

func (r *repository) init(ctx context.Context, domain framework.ComplianceDomain, scrapeResults map[string]*compliance.ComplianceReturn, f *factory) error {
	r.cluster = domain.Cluster().Cluster()
	r.nodes = nodesByID(framework.Nodes(domain))

	deployments := framework.Deployments(domain)
	r.deployments = deploymentsByID(deployments)

	clusterID := r.cluster.GetId()

	clusterQuery := search.NewQueryBuilder().AddExactMatches(search.ClusterID, clusterID).ProtoQuery()
	infPagination := &v1.QueryPagination{
		Limit: math.MaxInt32,
	}
	clusterQuery.Pagination = infPagination

	networkPolicies, err := f.networkPoliciesStore.GetNetworkPolicies(ctx, clusterID, "")
	if err != nil {
		return err
	}
	r.networkPolicies = networkPoliciesByID(networkPolicies)

	r.networkGraph = f.networkGraphEvaluator.GetGraph(clusterID, deployments, networkPolicies)

	policies, err := f.policyStore.GetPolicies(ctx)
	if err != nil {
		return err
	}

	r.policies = policiesByName(policies)
	r.categoryToPolicies = policyCategories(policies)

	r.images, err = f.imageStore.SearchListImages(ctx, clusterQuery)
	if err != nil {
		return err
	}

	r.imageIntegrations, err = f.imageIntegrationStore.GetImageIntegrations(ctx,
		&v1.GetImageIntegrationsRequest{},
	)
	if err != nil {
		return err
	}

	for _, registryIntegration := range f.imageIntegrationsSet.RegistrySet().GetAll() {
		r.registries = append(r.registries, registryIntegration)
	}
	for _, scannerIntegration := range f.imageIntegrationsSet.ScannerSet().GetAll() {
		r.scanners = append(r.scanners, scannerIntegration)
	}

	r.processIndicators, err = f.processIndicatorStore.SearchRawProcessIndicators(ctx, clusterQuery)
	if err != nil {
		return err
	}

	flowStore := f.networkFlowDataStore.GetFlowStore(ctx, domain.Cluster().ID())
	r.networkFlows, _, err = flowStore.GetAllFlows(ctx, nil)
	if err != nil {
		return err
	}

	r.notifiers, err = f.notifierDataStore.GetNotifiers(ctx, &v1.GetNotifiersRequest{})
	if err != nil {
		return err
	}

	r.roles, err = f.roleDataStore.SearchRawRoles(ctx, clusterQuery)
	if err != nil {
		return err
	}

	r.bindings, err = f.bindingDataStore.SearchRawRoleBindings(ctx, clusterQuery)
	if err != nil {
		return err
	}

	alertQuery := search.ConjunctionQuery(
		clusterQuery,
		search.NewQueryBuilder().AddStrings(search.ViolationState, storage.ViolationState_ACTIVE.String(), storage.ViolationState_SNOOZED.String()).ProtoQuery(),
	)
	alertQuery.Pagination = infPagination
	r.unresolvedAlerts, err = f.alertStore.SearchListAlerts(ctx, alertQuery)
	if err != nil {
		return err
	}

	// Flatten the files so we can do direct lookups on the nested values
	for _, n := range scrapeResults {
		totalNodeFiles := make(map[string]*compliance.File)
		for path, file := range n.GetFiles() {
			expanded := expandFile(file)
			for k, v := range expanded {
				totalNodeFiles[k] = v
			}
			totalNodeFiles[path] = file
		}
		n.Files = totalNodeFiles
	}

	r.hostScrape = scrapeResults

	// check for latest compliance results to determine
	// if CIS benchmarks were ever run
	cisDockerStandardID, err := f.standardsRepo.GetCISDockerStandardID()
	if err != nil {
		return err
	}

	cisKubernetesStandardID, err := f.standardsRepo.GetCISKubernetesStandardID()
	if err != nil {
		return err
	}

	dockerCISRunResults, err := f.complianceStore.GetLatestRunResults(ctx, clusterID, cisDockerStandardID, 0)
	if err == nil && dockerCISRunResults.LastSuccessfulResults != nil {
		r.cisDockerRunCheck = true
	}

	kubeCISRunResults, err := f.complianceStore.GetLatestRunResults(ctx, clusterID, cisKubernetesStandardID, 0)
	if err == nil && kubeCISRunResults.LastSuccessfulResults != nil {
		r.cisKubernetesRunCheck = true
	}

	return nil
}
