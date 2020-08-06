package datastore

import (
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/alert/datastore/internal/commentsstore"
	"github.com/stackrox/rox/central/alert/datastore/internal/index"
	"github.com/stackrox/rox/central/alert/datastore/internal/search"
	"github.com/stackrox/rox/central/alert/datastore/internal/store/rocksdb"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/central/globalindex"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	once         sync.Once
	soleInstance DataStore
)

func initialize() {
	storage := rocksdb.NewFullStore(globaldb.GetRocksDB())
	commentsStorage := commentsstore.New(globaldb.GetGlobalDB())
	indexer := index.New(globalindex.GetAlertIndex())
	searcher := search.New(storage, indexer)
	var err error
	soleInstance, err = New(storage, commentsStorage, indexer, searcher)
	utils.Must(errors.Wrap(err, "unable to load datastore for alerts"))
}

// Singleton returns the sole instance of the DataStore service.
func Singleton() DataStore {
	once.Do(initialize)
	return soleInstance
}
