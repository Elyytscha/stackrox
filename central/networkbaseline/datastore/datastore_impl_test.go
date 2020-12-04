package datastore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stackrox/rox/central/networkbaseline/store"
	"github.com/stackrox/rox/central/networkbaseline/store/rocksdb"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/fixtures"
	pkgRocksDB "github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stretchr/testify/suite"
)

var (
	allAllowedCtx    = sac.WithAllAccess(context.Background())
	expectedBaseline = fixtures.GetNetworkBaseline()
)

func TestNetworkBaselineDatastoreSuite(t *testing.T) {
	suite.Run(t, new(NetworkBaselineDataStoreTestSuite))
}

type NetworkBaselineDataStoreTestSuite struct {
	suite.Suite

	datastore DataStore
	storage   store.Store

	mockCtrl *gomock.Controller
}

func (suite *NetworkBaselineDataStoreTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.mockCtrl = mockCtrl
	db, err := pkgRocksDB.NewTemp(suite.T().Name())
	if err != nil {
		suite.FailNowf("failed to create DB: %+v", err.Error())
	}
	suite.storage, err = rocksdb.New(db)
	if err != nil {
		suite.FailNowf("failed to create network baseline store: %+v", err.Error())
	}

	suite.datastore = newNetworkBaselineDataStore(suite.storage)
}

func (suite *NetworkBaselineDataStoreTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *NetworkBaselineDataStoreTestSuite) TestNoAccessAllowed() {
	// First create a baseline in datastore to make sure when we return false on get
	// we are indeed hitting permission issue
	suite.Nil(suite.datastore.UpsertNetworkBaseline(allAllowedCtx, expectedBaseline))

	ctx := sac.WithNoAccess(context.Background())

	_, ok, _ := suite.datastore.GetNetworkBaseline(ctx, expectedBaseline.GetDeploymentId())
	suite.False(ok)

	ok, _ = suite.datastore.Exists(ctx, expectedBaseline.GetDeploymentId())
	suite.False(ok)

	suite.Error(suite.datastore.UpsertNetworkBaseline(ctx, expectedBaseline), "permission denied")

	suite.Error(suite.datastore.DeleteNetworkBaseline(ctx, expectedBaseline.GetDeploymentId()), "permission denied")
	// BTW if we try to delete non-existent/already deleted baseline, it should just return nil
	suite.Nil(suite.datastore.DeleteNetworkBaseline(ctx, "non-existent deployment ID"))
}

func (suite *NetworkBaselineDataStoreTestSuite) TestNetworkBaselines() {
	// With all allowed access, we should be able to perform all ops on datastore
	// Create
	suite.Nil(suite.datastore.UpsertNetworkBaseline(allAllowedCtx, expectedBaseline))

	// Check exist and get
	ok, err := suite.datastore.Exists(allAllowedCtx, expectedBaseline.GetDeploymentId())
	suite.Nil(err)
	suite.True(ok)

	baseline, ok, err := suite.datastore.GetNetworkBaseline(allAllowedCtx, expectedBaseline.GetDeploymentId())
	suite.Nil(err)
	suite.True(ok)
	suite.Equal(expectedBaseline.GetClusterId(), baseline.GetClusterId())
	suite.Equal(expectedBaseline.GetNamespace(), baseline.GetNamespace())
	suite.Equal(expectedBaseline.GetLocked(), baseline.GetLocked())

	// Update
	originalBaselineLocked := expectedBaseline.GetLocked()
	expectedBaseline.Locked = !expectedBaseline.GetLocked()
	suite.Nil(suite.datastore.UpsertNetworkBaseline(allAllowedCtx, expectedBaseline))
	// Check update
	baseline, ok, err = suite.datastore.GetNetworkBaseline(allAllowedCtx, expectedBaseline.GetDeploymentId())
	suite.Nil(err)
	suite.True(ok)
	suite.NotEqual(originalBaselineLocked, baseline.GetLocked())
	suite.Equal(baseline.GetLocked(), expectedBaseline.GetLocked())

	// Delete
	suite.Nil(suite.datastore.DeleteNetworkBaseline(allAllowedCtx, expectedBaseline.GetDeploymentId()))

	// Verify deletion
	ok, err = suite.datastore.Exists(allAllowedCtx, expectedBaseline.GetDeploymentId())
	suite.Nil(err)
	suite.False(ok)
}

func (suite *NetworkBaselineDataStoreTestSuite) TestSAC() {
	ctxWithWrongClusterReadAccess :=
		sac.WithGlobalAccessScopeChecker(
			context.Background(),
			sac.AllowFixedScopes(
				sac.AccessModeScopeKeys(storage.Access_READ_ACCESS),
				sac.ResourceScopeKeys(resources.NetworkBaseline),
				sac.ClusterScopeKeys("a-wrong-cluster")))

	ctxWithReadAccess :=
		sac.WithGlobalAccessScopeChecker(
			context.Background(),
			sac.AllowFixedScopes(
				sac.AccessModeScopeKeys(storage.Access_READ_ACCESS),
				sac.ResourceScopeKeys(resources.NetworkBaseline),
				sac.ClusterScopeKeys(expectedBaseline.GetClusterId()),
				sac.NamespaceScopeKeys(expectedBaseline.GetNamespace())))

	ctxWithWrongClusterWriteAccess :=
		sac.WithGlobalAccessScopeChecker(
			context.Background(),
			sac.AllowFixedScopes(
				sac.AccessModeScopeKeys(storage.Access_READ_WRITE_ACCESS),
				sac.ResourceScopeKeys(resources.NetworkBaseline),
				sac.ClusterScopeKeys("a-wrong-cluster")))

	ctxWithWriteAccess :=
		sac.WithGlobalAccessScopeChecker(
			context.Background(),
			sac.AllowFixedScopes(
				sac.AccessModeScopeKeys(storage.Access_READ_WRITE_ACCESS),
				sac.ResourceScopeKeys(resources.NetworkBaseline),
				sac.ClusterScopeKeys(expectedBaseline.GetClusterId()),
				sac.NamespaceScopeKeys(expectedBaseline.GetNamespace())))

	// Test Create
	{
		suite.Error(suite.datastore.UpsertNetworkBaseline(ctxWithWrongClusterReadAccess, expectedBaseline), "permission denied")
		suite.Error(suite.datastore.UpsertNetworkBaseline(ctxWithReadAccess, expectedBaseline), "permission denied")
		suite.Error(suite.datastore.UpsertNetworkBaseline(ctxWithWrongClusterWriteAccess, expectedBaseline), "permission denied")
		suite.Nil(suite.datastore.UpsertNetworkBaseline(ctxWithWriteAccess, expectedBaseline))
	}

	// Test Exists. This should be tested after create so that when permissions are correct, we should be able to
	// actually get something.
	{
		ok, _ := suite.datastore.Exists(ctxWithWrongClusterReadAccess, expectedBaseline.GetDeploymentId())
		suite.False(ok)
		ok, _ = suite.datastore.Exists(ctxWithReadAccess, expectedBaseline.GetDeploymentId())
		suite.True(ok)
		ok, _ = suite.datastore.Exists(ctxWithWrongClusterWriteAccess, expectedBaseline.GetDeploymentId())
		suite.False(ok)
		ok, _ = suite.datastore.Exists(ctxWithWriteAccess, expectedBaseline.GetDeploymentId())
		suite.False(ok)
	}

	// Test Get
	{
		_, ok, _ := suite.datastore.GetNetworkBaseline(ctxWithWrongClusterReadAccess, expectedBaseline.GetDeploymentId())
		suite.False(ok)
		_, ok, _ = suite.datastore.GetNetworkBaseline(ctxWithReadAccess, expectedBaseline.GetDeploymentId())
		suite.True(ok)
		_, ok, _ = suite.datastore.GetNetworkBaseline(ctxWithWrongClusterWriteAccess, expectedBaseline.GetDeploymentId())
		suite.False(ok)
		_, ok, _ = suite.datastore.GetNetworkBaseline(ctxWithWriteAccess, expectedBaseline.GetDeploymentId())
		suite.False(ok)
	}

	// Test Delete
	{
		suite.Error(suite.datastore.DeleteNetworkBaseline(ctxWithWrongClusterReadAccess, expectedBaseline.GetDeploymentId()), "permission denied")
		suite.Error(suite.datastore.DeleteNetworkBaseline(ctxWithReadAccess, expectedBaseline.GetDeploymentId()), "permission denied")
		suite.Error(suite.datastore.DeleteNetworkBaseline(ctxWithWrongClusterWriteAccess, expectedBaseline.GetDeploymentId()), "permission denied")
		suite.Nil(suite.datastore.DeleteNetworkBaseline(ctxWithWriteAccess, expectedBaseline.GetDeploymentId()))
	}
}
