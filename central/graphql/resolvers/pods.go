package resolvers

import (
	"context"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/graph-gophers/graphql-go"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/generated/storage"
	pkgMetrics "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/utils"
)

func init() {
	const resolverName = "Pod"
	schema := getBuilder()
	utils.Must(
		schema.AddQuery("pod(id: ID): Pod"),
		schema.AddQuery("pods(query: String, pagination: Pagination): [Pod!]!"),
		schema.AddQuery("podCount(query: String): Int!"),
		schema.AddExtraResolver(resolverName, "containerCount: Int!"),
		schema.AddExtraResolver(resolverName, "policyViolationEvents: [PolicyViolationEvent!]!"),
		schema.AddExtraResolver(resolverName, "processActivityEvents: [ProcessActivityEvent!]!"),
		schema.AddExtraResolver(resolverName, "containerRestartEvents: [ContainerRestartEvent!]!"),
		schema.AddExtraResolver(resolverName, "containerTerminationEvents: [ContainerTerminationEvent!]!"),
		schema.AddExtraResolver(resolverName, "events: [DeploymentEvent!]!"),
	)
}

// Pod returns a GraphQL resolver for a given id.
func (resolver *Resolver) Pod(ctx context.Context, args struct{ *graphql.ID }) (*podResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "Pod")
	if err := readDeployments(ctx); err != nil {
		return nil, err
	}
	return resolver.wrapPod(resolver.PodDataStore.GetPod(ctx, string(*args.ID)))
}

// Pods returns GraphQL resolvers for all pods.
func (resolver *Resolver) Pods(ctx context.Context, args PaginatedQuery) ([]*podResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "Pods")
	if err := readDeployments(ctx); err != nil {
		return nil, err
	}
	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return nil, err
	}
	return resolver.wrapPods(resolver.PodDataStore.SearchRawPods(ctx, q))
}

// PodCount returns count of all pods across deployments
func (resolver *Resolver) PodCount(ctx context.Context, args RawQuery) (int32, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Root, "PodCount")
	if err := readDeployments(ctx); err != nil {
		return 0, err
	}
	q, err := args.AsV1QueryOrEmpty()
	if err != nil {
		return 0, err
	}
	results, err := resolver.PodDataStore.Search(ctx, q)
	if err != nil {
		return 0, err
	}
	return int32(len(results)), nil
}

// ContainerCount returns the number of active containers.
// Active is defined by being present in the pod spec.
func (resolver *podResolver) ContainerCount() int32 {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Pods, "ContainerCount")

	containerNames := set.NewStringSet()
	for _, instance := range resolver.data.GetLiveInstances() {
		containerNames.Add(instance.GetContainerName())
	}

	// It is possible that a container is "active", but there are currently no "live" instances.
	for _, instanceList := range resolver.data.GetTerminatedInstances() {
		instances := instanceList.GetInstances()
		if len(instances) == 0 {
			continue
		}

		containerNames.Add(instances[0].GetContainerName())
	}

	return int32(containerNames.Cardinality())
}

// PolicyViolationEvents returns all policy violations associated with this pod.
func (resolver *podResolver) PolicyViolationEvents(ctx context.Context) ([]*PolicyViolationEventResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Pods, "PolicyViolationEvents")

	// We search by PodID (name) to filter out processes involving other pods.
	// PodID is guaranteed to be unique within a deployment during the pod's
	// lifetime and all process indicators should have this field.
	// Not all process indicators will have PodUID, so we cannot filter based on that.
	q := search.NewConjunctionQuery(
		search.NewQueryBuilder().AddExactMatches(search.DeploymentID, resolver.data.GetDeploymentId()).ProtoQuery(),
		search.NewQueryBuilder().AddExactMatches(search.ViolationState, storage.ViolationState_ACTIVE.String()).ProtoQuery(),
		search.NewQueryBuilder().AddExactMatches(search.LifecycleStage, storage.LifecycleStage_RUNTIME.String()).ProtoQuery(),
		search.NewQueryBuilder().AddExactMatches(search.PodID, resolver.data.GetName()).ProtoQuery(),
	)

	return resolver.root.getPolicyViolationEvents(ctx, q)
}

// ProcessActivityEvents returns all the process activities associated with this pod.
func (resolver *podResolver) ProcessActivityEvents(ctx context.Context) ([]*ProcessActivityEventResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Pods, "ProcessActivityEvents")

	// It is possible that not all process indicators have PodUID populated. For now, it is safer to not use it.
	// PodID (name) is unique within a deployment.
	query := search.NewConjunctionQuery(
		search.NewQueryBuilder().AddExactMatches(search.DeploymentID, resolver.data.GetDeploymentId()).ProtoQuery(),
		search.NewQueryBuilder().AddExactMatches(search.PodID, resolver.data.GetName()).ProtoQuery(),
	)

	return resolver.root.getProcessActivityEvents(ctx, query)
}

// ContainerRestartEvents returns all the container restart events associated with this pod.
func (resolver *podResolver) ContainerRestartEvents() []*ContainerRestartEventResolver {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Pods, "ContainerRestartEvents")

	var events []*ContainerRestartEventResolver
	liveInstances := resolver.data.GetLiveInstances()
	liveInstancesByName := make(map[string]*storage.ContainerInstance, len(liveInstances))
	for _, liveInstance := range liveInstances {
		liveInstancesByName[liveInstance.GetContainerName()] = liveInstance
	}

	for _, instances := range resolver.data.GetTerminatedInstances() {
		terminatedInstances := instances.GetInstances()
		if len(terminatedInstances) == 0 {
			continue
		}

		// The first terminated instance could not have been created from a restart.
		for i := 1; i < len(terminatedInstances); i++ {
			timestamp, ok := convertTimestamp(terminatedInstances[i].GetContainerName(), "container instance", terminatedInstances[i].GetStarted())
			if !ok {
				continue
			}
			events = append(events, &ContainerRestartEventResolver{
				id:        graphql.ID(terminatedInstances[i].GetInstanceId().GetId()),
				name:      terminatedInstances[i].GetContainerName(),
				timestamp: timestamp,
			})
		}

		// A current live instance can be due to a restart.
		containerName := terminatedInstances[0].GetContainerName()
		if instance, exists := liveInstancesByName[containerName]; exists {
			timestamp, ok := convertTimestamp(instance.GetContainerName(), "container instance", instance.GetStarted())
			if !ok {
				continue
			}
			events = append(events, &ContainerRestartEventResolver{
				id:        graphql.ID(instance.GetInstanceId().GetId()),
				name:      instance.GetContainerName(),
				timestamp: timestamp,
			})
			delete(liveInstancesByName, containerName)
		}
	}
	return events
}

// ContainerTerminationEvents returns all the container termination events associated with this pod.
func (resolver *podResolver) ContainerTerminationEvents() []*ContainerTerminationEventResolver {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Pods, "ContainerTerminationEvents")

	var events []*ContainerTerminationEventResolver
	for _, instances := range resolver.data.GetTerminatedInstances() {
		for _, instance := range instances.GetInstances() {
			timestamp, ok := convertTimestamp(instance.GetContainerName(), "container instance", instance.GetFinished())
			if !ok {
				continue
			}
			events = append(events, &ContainerTerminationEventResolver{
				id:        graphql.ID(instance.GetInstanceId().GetId()),
				name:      instance.GetContainerName(),
				timestamp: timestamp,
				exitCode:  instance.GetExitCode(),
				reason:    instance.GetTerminationReason(),
			})
		}
	}

	return events
}

// Events returns all events associated with this pod.
func (resolver *podResolver) Events(ctx context.Context) ([]*DeploymentEventResolver, error) {
	defer metrics.SetGraphQLOperationDurationTime(time.Now(), pkgMetrics.Pods, "Events")

	var events []*DeploymentEventResolver

	policyViolations, err := resolver.PolicyViolationEvents(ctx)
	if err != nil {
		return nil, err
	}
	for _, policyViolation := range policyViolations {
		events = append(events, &DeploymentEventResolver{policyViolation})
	}

	processActivities, err := resolver.ProcessActivityEvents(ctx)
	if err != nil {
		return nil, err
	}
	for _, processActivity := range processActivities {
		events = append(events, &DeploymentEventResolver{processActivity})
	}

	containerRestarts := resolver.ContainerRestartEvents()
	for _, containerRestart := range containerRestarts {
		events = append(events, &DeploymentEventResolver{containerRestart})
	}

	containerTerminations := resolver.ContainerTerminationEvents()
	for _, containerTermination := range containerTerminations {
		events = append(events, &DeploymentEventResolver{containerTermination})
	}

	return events, nil
}

// convertTimestamp is a thin wrapper over types.TimestampFromProto.
// That function is used often and, if it errors, we log errors each time.
// This function just saves the need to write so much repeat code.
func convertTimestamp(name, resource string, t *types.Timestamp) (time.Time, bool) {
	timestamp, err := types.TimestampFromProto(t)
	if err != nil {
		log.Errorf("unable to convert timestamp for %s %s: %v", resource, name, err)
	}
	return timestamp, err == nil
}
