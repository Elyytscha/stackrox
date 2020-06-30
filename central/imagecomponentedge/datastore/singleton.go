package datastore

import (
	globaldb "github.com/stackrox/rox/central/globaldb/dackbox"
	"github.com/stackrox/rox/central/globalindex"
	"github.com/stackrox/rox/central/imagecomponentedge/index"
	"github.com/stackrox/rox/central/imagecomponentedge/search"
	"github.com/stackrox/rox/central/imagecomponentedge/store/dackbox"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	once sync.Once

	ad DataStore
)

func initialize() {
	storage, err := dackbox.New(globaldb.GetGlobalDackBox())
	utils.Must(err)

	searcher := search.New(storage, index.New(globalindex.GetGlobalIndex()))

	ad, err = New(globaldb.GetGlobalDackBox(), storage, index.New(globalindex.GetGlobalIndex()), searcher)
	utils.Must(err)
}

// Singleton provides the interface for non-service external interaction.
func Singleton() DataStore {
	once.Do(initialize)
	return ad
}
