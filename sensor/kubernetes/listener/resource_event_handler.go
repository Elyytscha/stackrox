package listener

import (
	"github.com/openshift/client-go/apps/informers/externalversions"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/kubernetes"
	"github.com/stackrox/rox/sensor/common/clusterentities"
	"github.com/stackrox/rox/sensor/common/processfilter"
	"github.com/stackrox/rox/sensor/common/roxmetadata"
	"github.com/stackrox/rox/sensor/kubernetes/listener/resources"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func handleAllEvents(sif informers.SharedInformerFactory, osf externalversions.SharedInformerFactory, output chan<- *central.SensorEvent, stopSignal *concurrency.Signal) {
	// We want creates to be treated as updates while existing objects are loaded.
	var treatCreatesAsUpdates concurrency.Flag
	treatCreatesAsUpdates.Set(true)

	// Create the dispatcher registry, which provides dispatchers to all of the handlers.
	podInformer := sif.Core().V1().Pods()
	dispatchers := resources.NewDispatcherRegistry(podInformer.Lister(), clusterentities.StoreInstance(), roxmetadata.Singleton(), processfilter.Singleton())

	namespaceInformer := sif.Core().V1().Namespaces().Informer()
	secretInformer := sif.Core().V1().Secrets().Informer()
	saInformer := sif.Core().V1().ServiceAccounts().Informer()

	roleInformer := sif.Rbac().V1().Roles().Informer()
	clusterRoleInformer := sif.Rbac().V1().ClusterRoles().Informer()
	roleBindingInformer := sif.Rbac().V1().RoleBindings().Informer()
	clusterRoleBindingInformer := sif.Rbac().V1().ClusterRoleBindings().Informer()

	// prePodWaitGroup
	prePodWaitGroup := &concurrency.WaitGroup{}

	// Informers that need to be synced initially
	handle(namespaceInformer, dispatchers.ForNamespaces(), output, nil, prePodWaitGroup, stopSignal)
	handle(secretInformer, dispatchers.ForSecrets(), output, nil, prePodWaitGroup, stopSignal)
	handle(saInformer, dispatchers.ForServiceAccounts(), output, nil, prePodWaitGroup, stopSignal)

	// RBAC dispatchers handles multiple sets of data
	handle(roleInformer, dispatchers.ForRoles(), output, nil, prePodWaitGroup, stopSignal)
	handle(clusterRoleInformer, dispatchers.ForClusterRoles(), output, nil, prePodWaitGroup, stopSignal)
	handle(roleBindingInformer, dispatchers.ForRoleBindings(), output, nil, prePodWaitGroup, stopSignal)
	handle(clusterRoleBindingInformer, dispatchers.ForClusterRoleBindings(), output, nil, prePodWaitGroup, stopSignal)

	sif.Start(stopSignal.Done())

	// Run the namespace and rbac object informers first since other handlers rely on their outputs.
	informersToSync := []cache.SharedInformer{namespaceInformer, secretInformer,
		saInformer, roleInformer, clusterRoleInformer, roleBindingInformer, clusterRoleBindingInformer}
	syncFuncs := make([]cache.InformerSynced, len(informersToSync))
	for i, informer := range informersToSync {
		syncFuncs[i] = informer.HasSynced
	}
	cache.WaitForCacheSync(stopSignal.Done(), syncFuncs...)

	if !concurrency.WaitInContext(prePodWaitGroup, stopSignal) {
		return
	}

	// Run the pod informer second since other handlers rely on its output.
	podWaitGroup := &concurrency.WaitGroup{}
	handle(podInformer.Informer(), dispatchers.ForDeployments(kubernetes.Pod), output, &treatCreatesAsUpdates, podWaitGroup, stopSignal)
	sif.Start(stopSignal.Done())
	cache.WaitForCacheSync(stopSignal.Done(), podInformer.Informer().HasSynced)

	if !concurrency.WaitInContext(podWaitGroup, stopSignal) {
		return
	}

	wg := &concurrency.WaitGroup{}

	// Non-deployment types.
	handle(sif.Networking().V1().NetworkPolicies().Informer(), dispatchers.ForNetworkPolicies(), output, nil, wg, stopSignal)
	handle(sif.Core().V1().Nodes().Informer(), dispatchers.ForNodes(), output, nil, wg, stopSignal)
	handle(sif.Core().V1().Services().Informer(), dispatchers.ForServices(), output, nil, wg, stopSignal)

	// Deployment types.
	handle(sif.Extensions().V1beta1().DaemonSets().Informer(), dispatchers.ForDeployments(kubernetes.DaemonSet), output, &treatCreatesAsUpdates, wg, stopSignal)
	handle(sif.Extensions().V1beta1().Deployments().Informer(), dispatchers.ForDeployments(kubernetes.Deployment), output, &treatCreatesAsUpdates, wg, stopSignal)
	handle(sif.Extensions().V1beta1().ReplicaSets().Informer(), dispatchers.ForDeployments(kubernetes.ReplicaSet), output, &treatCreatesAsUpdates, wg, stopSignal)
	handle(sif.Core().V1().ReplicationControllers().Informer(), dispatchers.ForDeployments(kubernetes.ReplicationController), output, &treatCreatesAsUpdates, wg, stopSignal)
	handle(sif.Apps().V1beta1().StatefulSets().Informer(), dispatchers.ForDeployments(kubernetes.StatefulSet), output, &treatCreatesAsUpdates, wg, stopSignal)
	handle(sif.Batch().V1().Jobs().Informer(), dispatchers.ForDeployments(kubernetes.Job), output, &treatCreatesAsUpdates, wg, stopSignal)
	handle(sif.Batch().V1beta1().CronJobs().Informer(), dispatchers.ForDeployments(kubernetes.CronJob), output, &treatCreatesAsUpdates, wg, stopSignal)

	if osf != nil {
		handle(osf.Apps().V1().DeploymentConfigs().Informer(), dispatchers.ForDeployments(kubernetes.DeploymentConfig), output, &treatCreatesAsUpdates, wg, stopSignal)
	}

	// SharedInformerFactories can have Start called multiple times which will start the rest of the handlers
	sif.Start(stopSignal.Done())
	if osf != nil {
		osf.Start(stopSignal.Done())
	}

	// WaitForCacheSync synchronization is broken for SharedIndexInformers due to internal addCh/pendingNotifications
	// copy.  We have implemented our own sync in order to work around this.

	if !concurrency.WaitInContext(wg, stopSignal) {
		return
	}

	// Set the flag that all objects present at start up have been consumed.
	treatCreatesAsUpdates.Set(false)

	output <- &central.SensorEvent{
		Resource: &central.SensorEvent_Synced{
			Synced: &central.SensorEvent_ResourcesSynced{},
		},
	}
}

// Helper function that creates and adds a handler to an informer.
//////////////////////////////////////////////////////////////////
func handle(informer cache.SharedIndexInformer, dispatcher resources.Dispatcher, output chan<- *central.SensorEvent, treatCreatesAsUpdates *concurrency.Flag, wg *concurrency.WaitGroup, stopSignal *concurrency.Signal) {
	handlerImpl := &resourceEventHandlerImpl{
		dispatcher:            dispatcher,
		output:                output,
		treatCreatesAsUpdates: treatCreatesAsUpdates,

		hasSeenAllInitialIDsSignal: concurrency.NewSignal(),
		seenIDs:                    make(map[types.UID]struct{}),
		missingInitialIDs:          nil,
	}
	informer.AddEventHandler(handlerImpl)
	wg.Add(1)
	go func() {
		defer wg.Add(-1)
		if !cache.WaitForCacheSync(stopSignal.Done(), informer.HasSynced) {
			return
		}
		doneChannel := handlerImpl.PopulateInitialObjects(informer.GetIndexer().List())
		select {
		case <-stopSignal.Done():
		case <-doneChannel:
		}
	}()
}
