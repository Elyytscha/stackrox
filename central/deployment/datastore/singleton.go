package datastore

import (
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/deployment/cache"
	"github.com/stackrox/rox/central/deployment/datastore/internal/processtagsstore"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/globaldb/dackbox"
	"github.com/stackrox/rox/central/globalindex"
	imageDatastore "github.com/stackrox/rox/central/image/datastore"
	nfDS "github.com/stackrox/rox/central/networkflow/datastore"
	"github.com/stackrox/rox/central/processindicator/filter"
	pwDS "github.com/stackrox/rox/central/processwhitelist/datastore"
	"github.com/stackrox/rox/central/ranking"
	riskDS "github.com/stackrox/rox/central/risk/datastore"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	once sync.Once

	ad DataStore

	log = logging.LoggerForModule()
)

func initialize() {
	var err error
	ad, err = New(dackbox.GetGlobalDackBox(),
		dackbox.GetKeyFence(),
		processtagsstore.New(globaldb.GetGlobalDB()),
		globalindex.GetGlobalIndex(),
		globalindex.GetProcessIndex(),
		imageDatastore.Singleton(),
		pwDS.Singleton(),
		nfDS.Singleton(),
		riskDS.Singleton(),
		cache.DeletedDeploymentCacheSingleton(),
		filter.Singleton(),
		ranking.ClusterRanker(),
		ranking.NamespaceRanker(),
		ranking.DeploymentRanker())
	utils.Must(errors.Wrap(err, "unable to load datastore for deployments"))
}

// Singleton provides the interface for non-service external interaction.
func Singleton() DataStore {
	once.Do(initialize)
	return ad
}
