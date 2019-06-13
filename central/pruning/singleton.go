package pruning

import (
	alertDatastore "github.com/stackrox/rox/central/alert/datastore"
	configDatastore "github.com/stackrox/rox/central/config/datastore"
	deploymentDatastore "github.com/stackrox/rox/central/deployment/datastore"
	imagesDatastore "github.com/stackrox/rox/central/image/datastore"

	"github.com/stackrox/rox/pkg/sync"
)

var (
	once sync.Once
	gc   GarbageCollector
)

// Singleton returns the global instance of the garbage collection
func Singleton() GarbageCollector {
	once.Do(func() {
		gc = newGarbageCollector(alertDatastore.Singleton(), imagesDatastore.Singleton(), deploymentDatastore.Singleton(), configDatastore.Singleton())
	})
	return gc
}
