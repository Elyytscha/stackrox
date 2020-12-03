package listener

import (
	"github.com/openshift/client-go/apps/informers/externalversions"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/kubernetes"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/sensor/common/clusterentities"
	"github.com/stackrox/rox/sensor/common/config"
	"github.com/stackrox/rox/sensor/common/detector"
	"github.com/stackrox/rox/sensor/common/processfilter"
	"github.com/stackrox/rox/sensor/kubernetes/listener/resources"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func handleAllEvents(sif, resyncingSif informers.SharedInformerFactory, osf externalversions.SharedInformerFactory, output chan<- *central.MsgFromSensor,
	stopSignal *concurrency.Signal, config config.Handler, detector detector.Detector) {
	// We want creates to be treated as updates while existing objects are loaded.
	var treatCreatesAsUpdates concurrency.Flag
	treatCreatesAsUpdates.Set(true)

	var syncedRBAC concurrency.Flag

	// Create the dispatcher registry, which provides dispatchers to all of the handlers.
	podInformer := resyncingSif.Core().V1().Pods()
	dispatchers := resources.NewDispatcherRegistry(&syncedRBAC, podInformer.Lister(), clusterentities.StoreInstance(), processfilter.Singleton(), config, detector)

	namespaceInformer := sif.Core().V1().Namespaces().Informer()
	secretInformer := sif.Core().V1().Secrets().Informer()
	saInformer := sif.Core().V1().ServiceAccounts().Informer()

	roleInformer := resyncingSif.Rbac().V1().Roles().Informer()
	clusterRoleInformer := resyncingSif.Rbac().V1().ClusterRoles().Informer()
	roleBindingInformer := resyncingSif.Rbac().V1().RoleBindings().Informer()
	clusterRoleBindingInformer := resyncingSif.Rbac().V1().ClusterRoleBindings().Informer()

	// prePodWaitGroup
	prePodWaitGroup := &concurrency.WaitGroup{}

	// we will single-thread event processing using this lock
	var eventLock sync.Mutex

	// Informers that need to be synced initially
	handle(namespaceInformer, dispatchers.ForNamespaces(), output, nil, prePodWaitGroup, stopSignal, &eventLock)
	handle(secretInformer, dispatchers.ForSecrets(), output, nil, prePodWaitGroup, stopSignal, &eventLock)
	handle(saInformer, dispatchers.ForServiceAccounts(), output, nil, prePodWaitGroup, stopSignal, &eventLock)

	// RBAC dispatchers handles multiple sets of data
	handle(roleInformer, dispatchers.ForRBAC(), output, nil, prePodWaitGroup, stopSignal, &eventLock)
	handle(clusterRoleInformer, dispatchers.ForRBAC(), output, nil, prePodWaitGroup, stopSignal, &eventLock)
	handle(roleBindingInformer, dispatchers.ForRBAC(), output, nil, prePodWaitGroup, stopSignal, &eventLock)
	handle(clusterRoleBindingInformer, dispatchers.ForRBAC(), output, nil, prePodWaitGroup, stopSignal, &eventLock)

	sif.Start(stopSignal.Done())
	resyncingSif.Start(stopSignal.Done())

	if !concurrency.WaitInContext(prePodWaitGroup, stopSignal) {
		return
	}
	syncedRBAC.Set(true)
	log.Info("Successfully synced namespaces, secrets, service accounts, roles and role bindings")

	// Wait for the pod informer to sync before processing other types.
	// This is required because the PodLister is used to populate the image ids of deployments.
	// However, do not ACTUALLY handle pod events yet -- those need to wait for deployments to be
	// synced, since we need to enrich pods with the deployment ids, and for that we need the entire
	// hierarchy to be populated.
	if !cache.WaitForCacheSync(stopSignal.Done(), podInformer.Informer().HasSynced) {
		return
	}
	log.Info("Successfully synced k8s pod cache")

	preTopLevelDeploymentWaitGroup := &concurrency.WaitGroup{}

	// Non-deployment types.
	handle(sif.Networking().V1().NetworkPolicies().Informer(), dispatchers.ForNetworkPolicies(), output, nil, preTopLevelDeploymentWaitGroup, stopSignal, &eventLock)
	handle(sif.Core().V1().Nodes().Informer(), dispatchers.ForNodes(), output, nil, preTopLevelDeploymentWaitGroup, stopSignal, &eventLock)
	handle(sif.Core().V1().Services().Informer(), dispatchers.ForServices(), output, nil, preTopLevelDeploymentWaitGroup, stopSignal, &eventLock)

	// Deployment subtypes (this ensures that the hierarchy maps are generated correctly)
	handle(resyncingSif.Batch().V1().Jobs().Informer(), dispatchers.ForJobs(), output, &treatCreatesAsUpdates, preTopLevelDeploymentWaitGroup, stopSignal, &eventLock)
	handle(resyncingSif.Apps().V1().ReplicaSets().Informer(), dispatchers.ForDeployments(kubernetes.ReplicaSet), output, &treatCreatesAsUpdates, preTopLevelDeploymentWaitGroup, stopSignal, &eventLock)
	handle(resyncingSif.Core().V1().ReplicationControllers().Informer(), dispatchers.ForDeployments(kubernetes.ReplicationController), output, &treatCreatesAsUpdates, preTopLevelDeploymentWaitGroup, stopSignal, &eventLock)

	sif.Start(stopSignal.Done())
	resyncingSif.Start(stopSignal.Done())

	if !concurrency.WaitInContext(preTopLevelDeploymentWaitGroup, stopSignal) {
		return
	}

	log.Info("Successfully synced network policies, nodes, services, jobs, replica sets, and replication controllers")

	wg := &concurrency.WaitGroup{}

	// Deployment types.
	handle(resyncingSif.Apps().V1().DaemonSets().Informer(), dispatchers.ForDeployments(kubernetes.DaemonSet), output, &treatCreatesAsUpdates, wg, stopSignal, &eventLock)
	handle(resyncingSif.Apps().V1().Deployments().Informer(), dispatchers.ForDeployments(kubernetes.Deployment), output, &treatCreatesAsUpdates, wg, stopSignal, &eventLock)
	handle(resyncingSif.Apps().V1().StatefulSets().Informer(), dispatchers.ForDeployments(kubernetes.StatefulSet), output, &treatCreatesAsUpdates, wg, stopSignal, &eventLock)
	handle(resyncingSif.Batch().V1beta1().CronJobs().Informer(), dispatchers.ForDeployments(kubernetes.CronJob), output, &treatCreatesAsUpdates, wg, stopSignal, &eventLock)

	if osf != nil {
		handle(osf.Apps().V1().DeploymentConfigs().Informer(), dispatchers.ForDeployments(kubernetes.DeploymentConfig), output, &treatCreatesAsUpdates, wg, stopSignal, &eventLock)
	}

	// SharedInformerFactories can have Start called multiple times which will start the rest of the handlers
	sif.Start(stopSignal.Done())
	resyncingSif.Start(stopSignal.Done())
	if osf != nil {
		osf.Start(stopSignal.Done())
	}

	// WaitForCacheSync synchronization is broken for SharedIndexInformers due to internal addCh/pendingNotifications
	// copy.  We have implemented our own sync in order to work around this.
	if !concurrency.WaitInContext(wg, stopSignal) {
		return
	}

	log.Info("Successfully synced daemonsets, deployments, stateful sets and cronjobs")

	// Finally, run the pod informer, and process pod events.
	podWaitGroup := &concurrency.WaitGroup{}
	handle(podInformer.Informer(), dispatchers.ForDeployments(kubernetes.Pod), output, &treatCreatesAsUpdates, podWaitGroup, stopSignal, &eventLock)
	sif.Start(stopSignal.Done())

	if !concurrency.WaitInContext(podWaitGroup, stopSignal) {
		return
	}

	log.Info("Successfully synced pods")

	// Set the flag that all objects present at start up have been consumed.
	treatCreatesAsUpdates.Set(false)

	output <- &central.MsgFromSensor{
		Msg: &central.MsgFromSensor_Event{
			Event: &central.SensorEvent{
				Resource: &central.SensorEvent_Synced{
					Synced: &central.SensorEvent_ResourcesSynced{},
				},
			},
		},
	}
}

// Helper function that creates and adds a handler to an informer.
//////////////////////////////////////////////////////////////////
func handle(informer cache.SharedIndexInformer, dispatcher resources.Dispatcher, output chan<- *central.MsgFromSensor, treatCreatesAsUpdates *concurrency.Flag, wg *concurrency.WaitGroup, stopSignal *concurrency.Signal, eventLock *sync.Mutex) {
	handlerImpl := &resourceEventHandlerImpl{
		eventLock:             eventLock,
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
