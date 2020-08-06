package dackbox

import (
	"testing"

	"github.com/stackrox/rox/central/componentcveedge/store"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/dackbox"
	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/testutils/rocksdbtest"
	"github.com/stretchr/testify/suite"
)

func TestEdgeStore(t *testing.T) {
	suite.Run(t, new(EdgeStoreTestSuite))
}

type EdgeStoreTestSuite struct {
	suite.Suite

	db    *rocksdb.RocksDB
	dir   string
	dacky *dackbox.DackBox

	store store.Store
}

func (suite *EdgeStoreTestSuite) SetupSuite() {
	var err error
	suite.db, suite.dir, err = rocksdb.NewTemp("reference")
	if err != nil {
		suite.FailNowf("failed to create DB: %+v", err.Error())
	}
	suite.dacky, err = dackbox.NewRocksDBDackBox(suite.db, nil, []byte("graph"), []byte("dirty"), []byte("valid"))
	if err != nil {
		suite.FailNowf("failed to create counter: %+v", err.Error())
	}
	suite.store, err = New(suite.dacky)
	if err != nil {
		suite.FailNowf("failed to create counter: %+v", err.Error())
	}
}

func (suite *EdgeStoreTestSuite) TearDownSuite() {
	rocksdbtest.TearDownRocksDB(suite.db, suite.dir)
}

func (suite *EdgeStoreTestSuite) TestImages() {
	edges := []*storage.ComponentCVEEdge{
		{
			Id:        "CVE-2019-02-14",
			IsFixable: true,
		},
		{
			Id:        "CVE-2019-03-14",
			IsFixable: true,
		},
	}

	// Test Add
	for _, d := range edges {
		suite.NoError(suite.store.Upsert(d))
	}

	for _, d := range edges {
		got, exists, err := suite.store.Get(d.GetId())
		suite.NoError(err)
		suite.True(exists)
		suite.Equal(got, d)
	}

	// Test Update
	for _, d := range edges {
		d.IsFixable = false
	}

	suite.NoError(suite.store.Upsert(edges...))

	for _, d := range edges {
		got, exists, err := suite.store.Get(d.GetId())
		suite.NoError(err)
		suite.True(exists)
		suite.Equal(got, d)
	}

	// Test Count
	count, err := suite.store.Count()
	suite.NoError(err)
	suite.Equal(len(edges), count)
}
