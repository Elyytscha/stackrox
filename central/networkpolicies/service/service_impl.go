package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	clusterDataStore "github.com/stackrox/rox/central/cluster/datastore"
	deploymentDataStore "github.com/stackrox/rox/central/deployment/datastore"
	networkEntityDS "github.com/stackrox/rox/central/networkflow/datastore/entities"
	npDS "github.com/stackrox/rox/central/networkpolicies/datastore"
	"github.com/stackrox/rox/central/networkpolicies/generator"
	"github.com/stackrox/rox/central/networkpolicies/graph"
	notifierDataStore "github.com/stackrox/rox/central/notifier/datastore"
	"github.com/stackrox/rox/central/notifiers"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/central/sensor/service/connection"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/grpc/authn"
	"github.com/stackrox/rox/pkg/grpc/authz"
	"github.com/stackrox/rox/pkg/grpc/authz/perrpc"
	"github.com/stackrox/rox/pkg/grpc/authz/user"
	"github.com/stackrox/rox/pkg/k8sutil"
	"github.com/stackrox/rox/pkg/namespaces"
	networkPolicyConversion "github.com/stackrox/rox/pkg/protoconv/networkpolicy"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

var (
	authorizer = perrpc.FromMap(map[authz.Authorizer][]string{
		user.With(permissions.View(resources.NetworkPolicy)): {
			"/v1.NetworkPolicyService/GetNetworkPolicy",
			"/v1.NetworkPolicyService/GetNetworkPolicies",
			"/v1.NetworkPolicyService/SimulateNetworkGraph",
			"/v1.NetworkPolicyService/GetNetworkGraph",
			"/v1.NetworkPolicyService/GetNetworkGraphEpoch",
			"/v1.NetworkPolicyService/GetUndoModification",
		},
		user.With(permissions.Modify(resources.NetworkPolicy)): {
			"/v1.NetworkPolicyService/ApplyNetworkPolicy",
		},
		user.With(permissions.Modify(resources.Notifier)): {
			"/v1.NetworkPolicyService/SendNetworkPolicyYAML",
		},
		user.With(permissions.View(resources.NetworkPolicy), permissions.View(resources.NetworkGraph)): {
			"/v1.NetworkPolicyService/GenerateNetworkPolicies",
		},
	})
)

// serviceImpl provides APIs for alerts.
type serviceImpl struct {
	sensorConnMgr   connection.Manager
	clusterStore    clusterDataStore.DataStore
	deployments     deploymentDataStore.DataStore
	externalSrcs    networkEntityDS.EntityDataStore
	networkPolicies npDS.DataStore
	notifierStore   notifierDataStore.DataStore
	graphEvaluator  graph.Evaluator

	policyGenerator generator.Generator
}

// RegisterServiceServer registers this service with the given gRPC Server.
func (s *serviceImpl) RegisterServiceServer(grpcServer *grpc.Server) {
	v1.RegisterNetworkPolicyServiceServer(grpcServer, s)
}

// RegisterServiceHandler registers this service with the given gRPC Gateway endpoint.
func (s *serviceImpl) RegisterServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return v1.RegisterNetworkPolicyServiceHandler(ctx, mux, conn)
}

// AuthFuncOverride specifies the auth criteria for this API.
func (s *serviceImpl) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, authorizer.Authorized(ctx, fullMethodName)
}

func populateYAML(np *storage.NetworkPolicy) {
	k8sNetworkPolicy := networkPolicyConversion.RoxNetworkPolicyWrap{NetworkPolicy: np}.ToKubernetesNetworkPolicy()
	encoder := json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil)

	stringBuilder := &strings.Builder{}
	err := encoder.Encode(k8sNetworkPolicy, stringBuilder)
	if err != nil {
		np.Yaml = fmt.Sprintf("Could not render Network Policy YAML: %s", err)
		return
	}
	np.Yaml = stringBuilder.String()
}

func (s *serviceImpl) GetNetworkPolicy(ctx context.Context, request *v1.ResourceByID) (*storage.NetworkPolicy, error) {
	networkPolicy, exists, err := s.networkPolicies.GetNetworkPolicy(ctx, request.GetId())
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "network policy with id '%s' does not exist", request.GetId())
	}
	populateYAML(networkPolicy)
	return networkPolicy, nil
}

func (s *serviceImpl) GetNetworkPolicies(ctx context.Context, request *v1.GetNetworkPoliciesRequest) (*v1.NetworkPoliciesResponse, error) {
	// Check the cluster information.
	if err := s.clusterExists(ctx, request.GetClusterId()); err != nil {
		return nil, err
	}

	// Get the policies in the cluster
	networkPolicies, err := s.networkPolicies.GetNetworkPolicies(ctx, request.GetClusterId(), "")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// If there is a deployment query, filter the policies that apply to the deployments that match the query.
	if request.GetDeploymentQuery() != "" {
		// Get the deployments we want to check connectivity between.
		deployments, err := s.getDeployments(ctx, request.GetClusterId(), request.GetDeploymentQuery())
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		extSrcs, err := s.getExternalSrcs(ctx, request.GetClusterId())
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		networkPolicies = s.graphEvaluator.GetAppliedPolicies(deployments, extSrcs, networkPolicies)
	}

	// Fill in YAML fields where they are not set.
	for _, np := range networkPolicies {
		np.Yaml, err = networkPolicyConversion.RoxNetworkPolicyWrap{NetworkPolicy: np}.ToYaml()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	// Get the policies that apply to the fetched deployments.
	return &v1.NetworkPoliciesResponse{
		NetworkPolicies: networkPolicies,
	}, nil
}

func (s *serviceImpl) GetNetworkGraph(ctx context.Context, request *v1.GetNetworkGraphRequest) (*v1.NetworkGraph, error) {
	if !features.NetworkGraphPorts.Enabled() && request.GetIncludePorts() {
		return nil, status.Error(codes.Unimplemented, "support for ports in network policy graph is not enabled")
	}

	// Check that the cluster exists. If not there is nothing to we can process.
	if err := s.clusterExists(ctx, request.GetClusterId()); err != nil {
		return nil, err
	}

	// Gather all of the network policies that apply to the cluster and add the addition we are testing if applicable.
	networkPolicies, err := s.networkPolicies.GetNetworkPolicies(ctx, request.GetClusterId(), "")
	if err != nil {
		return nil, err
	}

	// Get the deployments we want to check connectivity between.
	deployments, err := s.getDeployments(ctx, request.GetClusterId(), request.GetQuery())
	if err != nil {
		return nil, err
	}

	extSrcs, err := s.getExternalSrcs(ctx, request.GetClusterId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Generate the graph.
	return s.graphEvaluator.GetGraph(request.GetClusterId(), deployments, extSrcs, networkPolicies, request.GetIncludePorts()), nil
}

func (s *serviceImpl) GetNetworkGraphEpoch(_ context.Context, req *v1.GetNetworkGraphEpochRequest) (*v1.NetworkGraphEpoch, error) {
	return &v1.NetworkGraphEpoch{
		Epoch: s.graphEvaluator.Epoch(req.GetClusterId()),
	}, nil
}

func (s *serviceImpl) ApplyNetworkPolicy(ctx context.Context, request *v1.ApplyNetworkPolicyYamlRequest) (*v1.Empty, error) {
	if request.GetModification().GetApplyYaml() == "" && len(request.GetModification().GetToDelete()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Modification must have contents")
	}

	clusterID := request.GetClusterId()
	conn := s.sensorConnMgr.GetConnection(clusterID)
	if conn == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "no active connection to cluster %q", clusterID)
	}

	undoMod, err := conn.NetworkPolicies().ApplyNetworkPolicies(ctx, request.GetModification())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not apply network policy modification: %v", err)
	}

	var user string
	identity := authn.IdentityFromContext(ctx)
	if identity != nil {
		user = identity.FriendlyName()
		if ap := identity.ExternalAuthProvider(); ap != nil {
			user += fmt.Sprintf(" [%s]", ap.Name())
		}
	}
	undoRecord := &storage.NetworkPolicyApplicationUndoRecord{
		User:                 user,
		ApplyTimestamp:       types.TimestampNow(),
		OriginalModification: request.GetModification(),
		UndoModification:     undoMod,
	}

	err = s.networkPolicies.UpsertUndoRecord(ctx, clusterID, undoRecord)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "network policy was applied, but undo record could not be stored: %v", err)
	}
	return &v1.Empty{}, nil
}

func (s *serviceImpl) SimulateNetworkGraph(ctx context.Context, request *v1.SimulateNetworkGraphRequest) (*v1.SimulateNetworkGraphResponse, error) {
	if !features.NetworkGraphPorts.Enabled() && request.GetIncludePorts() {
		return nil, status.Error(codes.Unimplemented, "support for ports in network policy simulation is not enabled")
	}

	// Check that the cluster exists. If not there is nothing to we can process.
	if err := s.clusterExists(ctx, request.GetClusterId()); err != nil {
		return nil, err
	}

	// Gather all of the network policies that apply to the cluster and add the addition we are testing if applicable.
	networkPoliciesInSimulation, err := s.getNetworkPoliciesInSimulation(ctx, request.GetClusterId(), request.GetModification())
	if err != nil {
		return nil, err
	}

	//Confirm that network policies in restricted namespaces are not changed
	err = validateNoForbiddenModification(networkPoliciesInSimulation)
	if err != nil {
		return nil, err
	}

	// Get the deployments we want to check connectivity between.
	deployments, err := s.getDeployments(ctx, request.GetClusterId(), request.GetQuery())
	if err != nil {
		return nil, err
	}

	// Generate the base graph.
	newPolicies := make([]*storage.NetworkPolicy, 0, len(networkPoliciesInSimulation))
	oldPolicies := make([]*storage.NetworkPolicy, 0, len(networkPoliciesInSimulation))
	var hasChanges bool
	for _, policyInSim := range networkPoliciesInSimulation {
		switch policyInSim.GetStatus() {
		case v1.NetworkPolicyInSimulation_UNCHANGED:
			oldPolicies = append(oldPolicies, policyInSim.GetPolicy())
			newPolicies = append(newPolicies, policyInSim.GetPolicy())
		case v1.NetworkPolicyInSimulation_ADDED:
			newPolicies = append(newPolicies, policyInSim.GetPolicy())
			hasChanges = true
		case v1.NetworkPolicyInSimulation_MODIFIED:
			oldPolicies = append(oldPolicies, policyInSim.GetOldPolicy())
			newPolicies = append(newPolicies, policyInSim.GetPolicy())
			hasChanges = true
		case v1.NetworkPolicyInSimulation_DELETED:
			oldPolicies = append(oldPolicies, policyInSim.GetOldPolicy())
			hasChanges = true
		default:
			return nil, status.Errorf(codes.Internal, "unhandled policy status %v", policyInSim.GetStatus())
		}
	}

	extSrcs, err := s.getExternalSrcs(ctx, request.GetClusterId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	newGraph := s.graphEvaluator.GetGraph(request.GetClusterId(), deployments, extSrcs, newPolicies, request.GetIncludePorts())
	result := &v1.SimulateNetworkGraphResponse{
		SimulatedGraph: newGraph,
		Policies:       networkPoliciesInSimulation,
	}
	if !hasChanges {
		// no need to compute diff - no new policies
		return result, nil
	}

	oldGraph := s.graphEvaluator.GetGraph(request.GetClusterId(), deployments, extSrcs, oldPolicies, request.GetIncludePorts())
	removedEdges, addedEdges, err := graph.ComputeDiff(oldGraph, newGraph)
	if err != nil {
		return nil, errors.Wrap(err, "could not compute a network graph diff")
	}

	result.Removed, result.Added = removedEdges, addedEdges
	return result, nil
}

func (s *serviceImpl) SendNetworkPolicyYAML(ctx context.Context, request *v1.SendNetworkPolicyYamlRequest) (*v1.Empty, error) {
	if request.GetClusterId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Cluster ID must be specified")
	}
	if len(request.GetNotifierIds()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Notifier IDs must be specified")
	}
	if request.GetModification().GetApplyYaml() == "" && len(request.GetModification().GetToDelete()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Modification must have contents")
	}

	clusterName, exists, err := s.clusterStore.GetClusterName(ctx, request.GetClusterId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve cluster: %s", err.Error())
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Cluster '%s' not found", request.GetClusterId())
	}

	errorList := errorhelpers.NewErrorList("unable to use all requested notifiers")
	for _, notifierID := range request.GetNotifierIds() {
		notifierProto, exists, err := s.notifierStore.GetNotifier(ctx, notifierID)
		if err != nil {
			errorList.AddError(err)
			continue
		}
		if !exists {
			errorList.AddStringf("notifier with id:%s not found", notifierID)
			continue
		}

		notifier, err := notifiers.CreateNotifier(notifierProto)
		if err != nil {
			errorList.AddStringf("error creating notifier with id:%s (%s) and type %s: %v", notifierProto.GetId(), notifierProto.GetName(), notifierProto.GetType(), err)
			continue
		}
		netpolNotifier, ok := notifier.(notifiers.NetworkPolicyNotifier)
		if !ok {
			errorList.AddStringf("notifier %s cannot notify on network policies", notifierProto.GetName())
			continue
		}

		err = netpolNotifier.NetworkPolicyYAMLNotify(ctx, request.GetModification().GetApplyYaml(), clusterName)
		if err != nil {
			errorList.AddStringf("error sending yaml notification to %s: %v", notifierProto.GetName(), err)
		}
	}

	err = errorList.ToError()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &v1.Empty{}, nil
}

func (s *serviceImpl) GenerateNetworkPolicies(ctx context.Context, req *v1.GenerateNetworkPoliciesRequest) (*v1.GenerateNetworkPoliciesResponse, error) {
	// Default to `none` delete existing mode.
	if req.DeleteExisting == v1.GenerateNetworkPoliciesRequest_UNKNOWN {
		req.DeleteExisting = v1.GenerateNetworkPoliciesRequest_NONE
	}

	if req.GetClusterId() == "" {
		return nil, status.Error(codes.InvalidArgument, "Cluster ID must be specified")
	}

	generated, toDelete, err := s.policyGenerator.Generate(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error generating network policies: %v", err)
	}

	var applyYAML string
	for _, generatedPolicy := range generated {
		yaml, err := networkPolicyConversion.RoxNetworkPolicyWrap{NetworkPolicy: generatedPolicy}.ToYaml()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error converting generated network policy to YAML: %v", err)
		}
		if applyYAML != "" {
			applyYAML += "\n---\n"
		}
		applyYAML += yaml
	}

	mod := &storage.NetworkPolicyModification{
		ApplyYaml: applyYAML,
		ToDelete:  toDelete,
	}

	return &v1.GenerateNetworkPoliciesResponse{
		Modification: mod,
	}, nil
}

func (s *serviceImpl) GetUndoModification(ctx context.Context, req *v1.GetUndoModificationRequest) (*v1.GetUndoModificationResponse, error) {
	undoRecord, exists, err := s.networkPolicies.GetUndoRecord(ctx, req.GetClusterId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not query undo store: %v", err)
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "no undo record stored for cluster %q", req.GetClusterId())
	}
	return &v1.GetUndoModificationResponse{
		UndoRecord: undoRecord,
	}, nil
}

func (s *serviceImpl) getDeployments(ctx context.Context, clusterID, query string) (deployments []*storage.Deployment, err error) {
	clusterQuery := search.NewQueryBuilder().AddExactMatches(search.ClusterID, clusterID).ProtoQuery()

	q := clusterQuery
	if query != "" {
		q, err = search.ParseQuery(query)
		if err != nil {
			return
		}
		q = search.ConjunctionQuery(q, clusterQuery)
	}

	deployments, err = s.deployments.SearchRawDeployments(ctx, q)
	return
}

func (s *serviceImpl) getExternalSrcs(ctx context.Context, clusterID string) ([]*storage.NetworkEntityInfo, error) {
	if !features.NetworkGraphExternalSrcs.Enabled() {
		return nil, nil
	}

	entities, err := s.externalSrcs.GetAllEntitiesForCluster(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	extSrcs := make([]*storage.NetworkEntityInfo, len(entities))
	for _, entity := range entities {
		extSrcs = append(extSrcs, entity.GetInfo())
	}
	return extSrcs, nil
}

func (s *serviceImpl) getNetworkPoliciesInSimulation(ctx context.Context, clusterID string, modification *storage.NetworkPolicyModification) ([]*v1.NetworkPolicyInSimulation, error) {
	additionalPolicies, err := compileValidateYaml(modification.GetApplyYaml())
	if err != nil {
		return nil, err
	}

	// Gather all of the network policies that apply to the cluster and add the addition we are testing if applicable.
	currentPolicies, err := s.networkPolicies.GetNetworkPolicies(ctx, clusterID, "")
	if err != nil {
		return nil, err
	}

	return applyPolicyModification(policyModification{
		ExistingPolicies: currentPolicies,
		NewPolicies:      additionalPolicies,
		ToDelete:         modification.GetToDelete(),
	})
}

type policyModification struct {
	ExistingPolicies []*storage.NetworkPolicy
	ToDelete         []*storage.NetworkPolicyReference
	NewPolicies      []*storage.NetworkPolicy
}

// applyPolicyModification returns the input slice of policies modified to use newPolicies.
// If oldPolicies contains a network policy with the same name and namespace as newPolicy, we consider newPolicy a
// replacement.
// If oldPolicies does not contain a network policy with a matching namespace and name, we consider it a new additional
// policy.
func applyPolicyModification(policies policyModification) (outputPolicies []*v1.NetworkPolicyInSimulation, err error) {
	outputPolicies = make([]*v1.NetworkPolicyInSimulation, 0, len(policies.NewPolicies)+len(policies.ExistingPolicies))
	policiesByRef := make(map[k8sutil.NSObjRef]*v1.NetworkPolicyInSimulation, len(policies.ExistingPolicies))
	for _, oldPolicy := range policies.ExistingPolicies {
		simPolicy := &v1.NetworkPolicyInSimulation{
			Policy: oldPolicy,
			Status: v1.NetworkPolicyInSimulation_UNCHANGED,
		}
		outputPolicies = append(outputPolicies, simPolicy)
		policiesByRef[k8sutil.RefOf(oldPolicy)] = simPolicy
	}

	// Delete policies that should be deleted
	for _, toDeleteRef := range policies.ToDelete {
		ref := k8sutil.RefOf(toDeleteRef)
		simPolicy := policiesByRef[ref]
		if simPolicy == nil {
			return nil, fmt.Errorf("policy %s in namespace %s marked for deletion does not exist", toDeleteRef.GetName(), toDeleteRef.GetNamespace())
		}

		if simPolicy.OldPolicy == nil {
			simPolicy.OldPolicy = simPolicy.Policy
		}
		simPolicy.Policy = nil
		simPolicy.Status = v1.NetworkPolicyInSimulation_DELETED
	}

	// Add new policies that have no matching old policies.
	for _, newPolicy := range policies.NewPolicies {
		oldPolicySim := policiesByRef[k8sutil.RefOf(newPolicy)]
		if oldPolicySim != nil {
			oldPolicySim.Status = v1.NetworkPolicyInSimulation_MODIFIED
			if oldPolicySim.OldPolicy == nil {
				oldPolicySim.OldPolicy = oldPolicySim.Policy
			}
			oldPolicySim.Policy = newPolicy
			continue
		}
		newPolicySim := &v1.NetworkPolicyInSimulation{
			Status: v1.NetworkPolicyInSimulation_ADDED,
			Policy: newPolicy,
		}
		outputPolicies = append(outputPolicies, newPolicySim)
	}

	// Fix IDs: For all modified policies, the ID of the new and old policies should be the same (that way the
	// diff does not get cluttered with just policy ID changes); for all new policies, we generate new, fresh UUIDs
	// that do not collide with any other IDs.
	// Rationale: IDs are (almost) meaningless - IDs from the simulation YAML will be changed by kubectl create/apply
	// anyway.
	for _, policy := range outputPolicies {
		if policy.GetStatus() == v1.NetworkPolicyInSimulation_MODIFIED {
			policy.Policy.Id = policy.GetOldPolicy().GetId()
		} else if policy.GetStatus() == v1.NetworkPolicyInSimulation_ADDED {
			policy.Policy.Id = uuid.NewV4().String()
		}
	}
	return
}

// compileValidateYaml compiles the YAML into a storage.NetworkPolicy, and verifies that a valid namespace exists.
func compileValidateYaml(simulationYaml string) ([]*storage.NetworkPolicy, error) {
	if simulationYaml == "" {
		return nil, nil
	}

	simulationYaml = strings.TrimPrefix(simulationYaml, "---\n")

	// Convert the YAMLs into rox network policy objects.
	policies, err := networkPolicyConversion.YamlWrap{Yaml: simulationYaml}.ToRoxNetworkPolicies()
	if err != nil {
		return nil, err
	}

	// Check that all resulting policies have namespaces.
	for _, policy := range policies {
		if policy.GetNamespace() == "" {
			return nil, errors.New("yamls tested against must apply to a namespace")
		}
	}

	return policies, nil
}

// validateNoForbiddenModification verifies whether network policy changes are not applied to 'stackrox' and 'kube-system' namespace
func validateNoForbiddenModification(networkPoliciesInSimulation []*v1.NetworkPolicyInSimulation) error {
	for _, policyInSim := range networkPoliciesInSimulation {
		policyNamespace := policyInSim.GetOldPolicy().GetNamespace()
		if policyNamespace == "" {
			policyNamespace = policyInSim.GetPolicy().GetNamespace()
		}

		if policyNamespace != namespaces.StackRox && policyNamespace != namespaces.KubeSystem {
			continue
		}

		if policyInSim.GetStatus() == v1.NetworkPolicyInSimulation_UNCHANGED {
			continue
		}

		policyName := policyInSim.GetPolicy().GetName()
		if policyInSim.GetStatus() != v1.NetworkPolicyInSimulation_MODIFIED {
			if policyInSim.GetStatus() != v1.NetworkPolicyInSimulation_ADDED {
				policyName = policyInSim.GetOldPolicy().GetName()
			}
			return errors.Errorf("%q cannot be applied since network policy change in '%q' namespace is forbidden", policyName, policyNamespace)
		}

		err := validateNoPolicyDiff(policyInSim.GetPolicy(), policyInSim.GetOldPolicy())
		if err != nil {
			return errors.Errorf("%q cannot be applied since network policy change in '%q' namespace is forbidden", policyName, policyNamespace)
		}
	}

	return nil
}

// validateNoPolicyDiff returns an error if the YAML of two network policies is different
func validateNoPolicyDiff(applyPolicy *storage.NetworkPolicy, currPolicy *storage.NetworkPolicy) error {
	if applyPolicy.GetYaml() != currPolicy.GetYaml() {
		return errors.New("network policies do not match")
	}

	return nil
}

func (s *serviceImpl) clusterExists(ctx context.Context, clusterID string) error {
	if clusterID == "" {
		return status.Error(codes.InvalidArgument, "cluster ID must be specified")
	}
	exists, err := s.clusterStore.Exists(ctx, clusterID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if !exists {
		return status.Errorf(codes.NotFound, "cluster with ID %q doesn't exist", clusterID)
	}
	return nil
}
