package resolvers

//go:generate go run ./gen

import (
	"context"
	"fmt"
	"reflect"

	violationsDatastore "github.com/stackrox/rox/central/alert/datastore"
	"github.com/stackrox/rox/central/apitoken/backend"
	clusterDatastore "github.com/stackrox/rox/central/cluster/datastore"
	clusterCVEEdgeDataStore "github.com/stackrox/rox/central/clustercveedge/datastore"
	"github.com/stackrox/rox/central/compliance/aggregation"
	complianceDS "github.com/stackrox/rox/central/compliance/datastore"
	complianceManager "github.com/stackrox/rox/central/compliance/manager"
	"github.com/stackrox/rox/central/compliance/manager/service"
	complianceService "github.com/stackrox/rox/central/compliance/service"
	complianceStandards "github.com/stackrox/rox/central/compliance/standards"
	componentCVEEdgeDataStore "github.com/stackrox/rox/central/componentcveedge/datastore"
	cveDataStore "github.com/stackrox/rox/central/cve/datastore"
	"github.com/stackrox/rox/central/cve/fetcher"
	cveMatcher "github.com/stackrox/rox/central/cve/matcher"
	deploymentDatastore "github.com/stackrox/rox/central/deployment/datastore"
	groupDataStore "github.com/stackrox/rox/central/group/datastore"
	imageDatastore "github.com/stackrox/rox/central/image/datastore"
	imageComponentDataStore "github.com/stackrox/rox/central/imagecomponent/datastore"
	imageComponentEdgeDataStore "github.com/stackrox/rox/central/imagecomponentedge/datastore"
	imageCVEEdgeDataStore "github.com/stackrox/rox/central/imagecveedge/datastore"
	namespaceDataStore "github.com/stackrox/rox/central/namespace/datastore"
	nfDS "github.com/stackrox/rox/central/networkgraph/flow/datastore"
	npDS "github.com/stackrox/rox/central/networkpolicies/datastore"
	nodeDataStore "github.com/stackrox/rox/central/node/globaldatastore"
	notifierDataStore "github.com/stackrox/rox/central/notifier/datastore"
	podDatastore "github.com/stackrox/rox/central/pod/datastore"
	policyDatastore "github.com/stackrox/rox/central/policy/datastore"
	whitelistStore "github.com/stackrox/rox/central/processbaseline/datastore"
	processIndicatorStore "github.com/stackrox/rox/central/processindicator/datastore"
	k8sroleStore "github.com/stackrox/rox/central/rbac/k8srole/datastore"
	k8srolebindingStore "github.com/stackrox/rox/central/rbac/k8srolebinding/datastore"
	riskDataStore "github.com/stackrox/rox/central/risk/datastore"
	roleDataStore "github.com/stackrox/rox/central/role/datastore"
	"github.com/stackrox/rox/central/role/resources"
	secretDataStore "github.com/stackrox/rox/central/secret/datastore"
	serviceAccountDataStore "github.com/stackrox/rox/central/serviceaccount/datastore"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/grpc/authz"
	"github.com/stackrox/rox/pkg/grpc/authz/user"
)

// Resolver is the root GraphQL resolver
type Resolver struct {
	ComplianceAggregator        aggregation.Aggregator
	APITokenBackend             backend.Backend
	ClusterDataStore            clusterDatastore.DataStore
	ComplianceDataStore         complianceDS.DataStore
	ComplianceStandardStore     complianceStandards.Repository
	ComplianceService           v1.ComplianceServiceServer
	ComplianceManagementService v1.ComplianceManagementServiceServer
	ComplianceManager           complianceManager.ComplianceManager
	clusterCVEEdgeDataStore     clusterCVEEdgeDataStore.DataStore
	ComponentCVEEdgeDataStore   componentCVEEdgeDataStore.DataStore
	CVEDataStore                cveDataStore.DataStore
	DeploymentDataStore         deploymentDatastore.DataStore
	PodDataStore                podDatastore.DataStore
	ImageDataStore              imageDatastore.DataStore
	ImageComponentDataStore     imageComponentDataStore.DataStore
	ImageComponentEdgeDataStore imageComponentEdgeDataStore.DataStore
	ImageCVEEdgeDataStore       imageCVEEdgeDataStore.DataStore
	GroupDataStore              groupDataStore.DataStore
	NamespaceDataStore          namespaceDataStore.DataStore
	NetworkFlowDataStore        nfDS.ClusterDataStore
	NetworkPoliciesStore        npDS.DataStore
	NodeGlobalDataStore         nodeDataStore.GlobalDataStore
	NotifierStore               notifierDataStore.DataStore
	PolicyDataStore             policyDatastore.DataStore
	ProcessIndicatorStore       processIndicatorStore.DataStore
	K8sRoleStore                k8sroleStore.DataStore
	K8sRoleBindingStore         k8srolebindingStore.DataStore
	RoleDataStore               roleDataStore.DataStore
	RiskDataStore               riskDataStore.DataStore
	SecretsDataStore            secretDataStore.DataStore
	ServiceAccountsDataStore    serviceAccountDataStore.DataStore
	ViolationsDataStore         violationsDatastore.DataStore
	WhiteListDataStore          whitelistStore.DataStore
	k8sIstioCVEManager          fetcher.K8sIstioCVEManager
	cveMatcher                  *cveMatcher.CVEMatcher
}

// New returns a Resolver wired into the relevant data stores
func New() *Resolver {
	resolver := &Resolver{
		ComplianceAggregator:        aggregation.Singleton(),
		APITokenBackend:             backend.Singleton(),
		ComplianceDataStore:         complianceDS.Singleton(),
		ComplianceStandardStore:     complianceStandards.RegistrySingleton(),
		ComplianceManagementService: service.Singleton(),
		ComplianceManager:           complianceManager.Singleton(),
		ComplianceService:           complianceService.Singleton(),
		ClusterDataStore:            clusterDatastore.Singleton(),
		clusterCVEEdgeDataStore:     clusterCVEEdgeDataStore.Singleton(),
		ComponentCVEEdgeDataStore:   componentCVEEdgeDataStore.Singleton(),
		CVEDataStore:                cveDataStore.Singleton(),
		DeploymentDataStore:         deploymentDatastore.Singleton(),
		PodDataStore:                podDatastore.Singleton(),
		ImageDataStore:              imageDatastore.Singleton(),
		ImageComponentDataStore:     imageComponentDataStore.Singleton(),
		ImageComponentEdgeDataStore: imageComponentEdgeDataStore.Singleton(),
		ImageCVEEdgeDataStore:       imageCVEEdgeDataStore.Singleton(),
		GroupDataStore:              groupDataStore.Singleton(),
		NamespaceDataStore:          namespaceDataStore.Singleton(),
		NetworkPoliciesStore:        npDS.Singleton(),
		NetworkFlowDataStore:        nfDS.Singleton(),
		NodeGlobalDataStore:         nodeDataStore.Singleton(),
		NotifierStore:               notifierDataStore.Singleton(),
		PolicyDataStore:             policyDatastore.Singleton(),
		ProcessIndicatorStore:       processIndicatorStore.Singleton(),
		K8sRoleStore:                k8sroleStore.Singleton(),
		K8sRoleBindingStore:         k8srolebindingStore.Singleton(),
		RiskDataStore:               riskDataStore.Singleton(),
		RoleDataStore:               roleDataStore.Singleton(),
		SecretsDataStore:            secretDataStore.Singleton(),
		ServiceAccountsDataStore:    serviceAccountDataStore.Singleton(),
		ViolationsDataStore:         violationsDatastore.Singleton(),
		WhiteListDataStore:          whitelistStore.Singleton(),
		k8sIstioCVEManager:          fetcher.SingletonManager(),
		cveMatcher:                  cveMatcher.Singleton(),
	}
	return resolver
}

//lint:file-ignore U1000 It's okay for some of the variables below to be unused.
var (
	readAlerts                 = readAuth(resources.Alert)
	readTokens                 = readAuth(resources.APIToken)
	readClusters               = readAuth(resources.Cluster)
	readCompliance             = readAuth(resources.Compliance)
	readComplianceRuns         = readAuth(resources.ComplianceRuns)
	readComplianceRunSchedule  = readAuth(resources.ComplianceRunSchedule)
	readDeployments            = readAuth(resources.Deployment)
	readGroups                 = readAuth(resources.Group)
	readImages                 = readAuth(resources.Image)
	readIndicators             = readAuth(resources.Indicator)
	readNamespaces             = readAuth(resources.Namespace)
	readNodes                  = readAuth(resources.Node)
	readNotifiers              = readAuth(resources.Notifier)
	readPolicies               = readAuth(resources.Policy)
	readK8sRoles               = readAuth(resources.K8sRole)
	readK8sRoleBindings        = readAuth(resources.K8sRoleBinding)
	readK8sSubjects            = readAuth(resources.K8sSubject)
	readRisks                  = readAuth(resources.Risk)
	readRoles                  = readAuth(resources.Role)
	readSecrets                = readAuth(resources.Secret)
	readServiceAccounts        = readAuth(resources.ServiceAccount)
	readWhitelists             = readAuth(resources.ProcessWhitelist)
	writeAlerts                = writeAuth(resources.Alert)
	writeCompliance            = writeAuth(resources.Compliance)
	writeComplianceRuns        = writeAuth(resources.ComplianceRuns)
	writeComplianceRunSchedule = writeAuth(resources.ComplianceRunSchedule)
	writeIndicators            = writeAuth(resources.Indicator)
)

type authorizerOverride struct{}

// SetAuthorizerOverride returns a context that will override the default permissions checking with custom
// logic. This is for testing only. It also feels pretty dangerous.
func SetAuthorizerOverride(ctx context.Context, authorizer authz.Authorizer) context.Context {
	return context.WithValue(ctx, authorizerOverride{}, authorizer)
}

func applyAuthorizer(authorizer authz.Authorizer) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		override := ctx.Value(authorizerOverride{})
		if override != nil {
			return override.(authz.Authorizer).Authorized(ctx, "graphql")
		}
		return authorizer.Authorized(ctx, "graphql")
	}
}

func readAuth(resource permissions.ResourceMetadata) func(ctx context.Context) error {
	return applyAuthorizer(user.With(permissions.View(resource)))
}

func writeAuth(resource permissions.ResourceMetadata) func(ctx context.Context) error {
	return applyAuthorizer(user.With(permissions.Modify(resource)))
}

func stringSlice(inputSlice interface{}) []string {
	r := reflect.ValueOf(inputSlice)
	output := make([]string, r.Len())
	for i := 0; i < r.Len(); i++ {
		output[i] = fmt.Sprint(r.Index(i).Interface())
	}
	return output
}
