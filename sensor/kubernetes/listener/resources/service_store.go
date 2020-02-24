package resources

import (
	"github.com/stackrox/rox/pkg/sync"
	v1 "k8s.io/api/core/v1"
	k8sLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
)

// serviceStore stores service objects (by namespace and UID)
type serviceStore struct {
	services         map[string]map[types.UID]*serviceWrap
	serviceLock      sync.RWMutex
	nodePortServices map[types.UID]*serviceWrap
}

// newServiceStore creates and returns a new service store.
func newServiceStore() *serviceStore {
	return &serviceStore{
		services:         make(map[string]map[types.UID]*serviceWrap),
		nodePortServices: make(map[types.UID]*serviceWrap),
	}
}

func (ss *serviceStore) addOrUpdateService(svc *serviceWrap) {
	ss.serviceLock.Lock()
	defer ss.serviceLock.Unlock()

	nsMap := ss.services[svc.Namespace]
	if nsMap == nil {
		nsMap = make(map[types.UID]*serviceWrap)
		ss.services[svc.Namespace] = nsMap
	}
	nsMap[svc.UID] = svc
	if svc.Spec.Type == v1.ServiceTypeNodePort || svc.Spec.Type == v1.ServiceTypeLoadBalancer {
		ss.nodePortServices[svc.UID] = svc
	} else {
		delete(ss.nodePortServices, svc.UID)
	}
}

func (ss *serviceStore) removeService(svc *v1.Service) {
	ss.serviceLock.Lock()
	defer ss.serviceLock.Unlock()

	nsMap := ss.services[svc.Namespace]
	if nsMap == nil {
		return
	}
	delete(nsMap, svc.UID)
	delete(ss.nodePortServices, svc.UID)
}

// OnNamespaceDeleted reacts to a namespace deletion, deleting all services in that namespace from the store.
func (ss *serviceStore) OnNamespaceDeleted(ns string) {
	ss.serviceLock.Lock()
	defer ss.serviceLock.Unlock()

	delete(ss.services, ns)
}

func (ss *serviceStore) getMatchingServices(namespace string, labels map[string]string) (matching []*serviceWrap) {
	labelSet := k8sLabels.Set(labels)
	ss.serviceLock.RLock()
	defer ss.serviceLock.RUnlock()
	for _, entry := range ss.services[namespace] {
		if entry.selector.Matches(labelSet) {
			matching = append(matching, entry)
		}
	}
	return
}

func (ss *serviceStore) getService(namespace string, uid types.UID) *serviceWrap {
	ss.serviceLock.RLock()
	defer ss.serviceLock.RUnlock()
	return ss.services[namespace][uid]
}
