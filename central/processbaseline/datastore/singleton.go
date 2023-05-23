package datastore

import (
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/processbaseline/index"
	"github.com/stackrox/rox/central/processbaseline/search"
	"github.com/stackrox/rox/central/processbaseline/store"
	pgStore "github.com/stackrox/rox/central/processbaseline/store/postgres"
	"github.com/stackrox/rox/central/processbaselineresults/datastore"
	indicatorStore "github.com/stackrox/rox/central/processindicator/datastore"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	once sync.Once

	ad DataStore

	log = logging.LoggerForModule()
)

func initialize() {
	var storage store.Store
	var indexer index.Indexer
	var err error
	storage, err = pgStore.NewWithCache(pgStore.New(globaldb.GetPostgres()))
	if err != nil {
		log.Fatal("failed to open process baseline store")
	}
	indexer = pgStore.NewIndexer(globaldb.GetPostgres())

	searcher, err := search.New(storage, indexer)
	if err != nil {
		panic("unable to load search index for process baseline")
	}

	ad = New(storage, indexer, searcher, datastore.Singleton(), indicatorStore.Singleton())
}

// Singleton provides the interface for non-service external interaction.
func Singleton() DataStore {
	once.Do(initialize)
	return ad
}
