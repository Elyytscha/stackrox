package service

import (
	"sync"

	alertDataStore "github.com/stackrox/rox/central/alert/datastore"
	clusterDataStore "github.com/stackrox/rox/central/cluster/datastore"
	deploymentDataStore "github.com/stackrox/rox/central/deployment/datastore"
	imageDataStore "github.com/stackrox/rox/central/image/datastore"
)

var (
	once sync.Once

	as Service
)

func initialize() {
	as = New(alertDataStore.Singleton(), clusterDataStore.Singleton(), deploymentDataStore.Singleton(), imageDataStore.Singleton())
}

// Singleton provides the instance of the Service interface to register.
func Singleton() Service {
	once.Do(initialize)
	return as
}
