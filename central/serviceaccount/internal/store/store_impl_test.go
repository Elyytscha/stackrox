package store

import (
	"testing"

	bolt "github.com/etcd-io/bbolt"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/bolthelper"
	"github.com/stackrox/rox/pkg/storecache"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/suite"
)

func TestServiceAccountStore(t *testing.T) {
	suite.Run(t, new(ServiceAccountStoreTestSuite))
}

type ServiceAccountStoreTestSuite struct {
	suite.Suite

	db *bolt.DB

	store Store
}

func (suite *ServiceAccountStoreTestSuite) SetupSuite() {
	db, err := bolthelper.NewTemp(suite.T().Name() + ".db")
	if err != nil {
		suite.FailNow("Failed to make BoltDB", err.Error())
	}

	suite.db = db
	suite.store, err = New(db, storecache.NewMapBackedCache())
	suite.Require().NoError(err)
}

func (suite *ServiceAccountStoreTestSuite) TearDownSuite() {
	testutils.TearDownDB(suite.db)
}

func (suite *ServiceAccountStoreTestSuite) TestServiceAccounts() {
	var serviceAccounts = []*storage.ServiceAccount{
		{
			Id: "serviceAccount1",
		},
		{
			Id: "serviceAccount2",
		},
	}

	for _, sa := range serviceAccounts {
		err := suite.store.UpsertServiceAccount(sa)
		suite.NoError(err)
	}

	// Get all service accounts
	retrievedServiceAccounts, err := suite.store.ListServiceAccounts()
	suite.Nil(err)
	suite.ElementsMatch(serviceAccounts, retrievedServiceAccounts)

	for _, s := range serviceAccounts {
		sa, exists, err := suite.store.GetServiceAccount(s.GetId())
		suite.NoError(err)
		suite.True(exists)
		suite.Equal(s, sa)
	}
}
