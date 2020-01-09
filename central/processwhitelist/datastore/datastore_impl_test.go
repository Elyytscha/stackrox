package datastore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stackrox/rox/central/globalindex"
	"github.com/stackrox/rox/central/processwhitelist/index"
	whitelistSearch "github.com/stackrox/rox/central/processwhitelist/search"
	"github.com/stackrox/rox/central/processwhitelist/store"
	"github.com/stackrox/rox/central/processwhitelistresults/datastore/mocks"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/bolthelper"
	"github.com/stackrox/rox/pkg/fixtures"
	"github.com/stackrox/rox/pkg/sac"
	pkgSearch "github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/uuid"
	"github.com/stretchr/testify/suite"
)

func TestProcessWhitelistDatastore(t *testing.T) {
	suite.Run(t, new(ProcessWhitelistDataStoreTestSuite))
}

type ProcessWhitelistDataStoreTestSuite struct {
	suite.Suite
	requestContext context.Context
	datastore      DataStore
	storage        store.Store
	indexer        index.Indexer
	searcher       whitelistSearch.Searcher

	whitelistResultsStore *mocks.MockDataStore

	mockCtrl *gomock.Controller
}

func (suite *ProcessWhitelistDataStoreTestSuite) SetupTest() {
	suite.requestContext = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.ProcessWhitelist),
		),
	)

	db, err := bolthelper.NewTemp(testutils.DBFileName(suite))

	suite.NoError(err)
	suite.storage, err = store.New(db)
	suite.NoError(err)

	tmpIndex, err := globalindex.TempInitializeIndices("")
	suite.NoError(err)
	suite.indexer = index.New(tmpIndex)

	suite.searcher, err = whitelistSearch.New(suite.storage, suite.indexer)
	suite.NoError(err)

	suite.mockCtrl = gomock.NewController(suite.T())

	suite.whitelistResultsStore = mocks.NewMockDataStore(suite.mockCtrl)
	suite.datastore = New(suite.storage, suite.indexer, suite.searcher, suite.whitelistResultsStore)
}

func (suite *ProcessWhitelistDataStoreTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *ProcessWhitelistDataStoreTestSuite) mustSerializeKey(key *storage.ProcessWhitelistKey) string {
	serialized, err := keyToID(key)
	suite.Require().NoError(err)
	return serialized
}

func (suite *ProcessWhitelistDataStoreTestSuite) createAndStoreWhitelist(key *storage.ProcessWhitelistKey) *storage.ProcessWhitelist {
	whitelist := fixtures.GetProcessWhitelist()
	whitelist.Key = key
	suite.NotNil(whitelist)
	id, err := suite.datastore.AddProcessWhitelist(suite.requestContext, whitelist)
	suite.NoError(err)
	suite.NotNil(id)
	suite.NotNil(whitelist.Created)
	suite.Equal(whitelist.Created, whitelist.LastUpdate)
	suite.True(whitelist.StackRoxLockedTimestamp.Compare(whitelist.Created) >= 0)

	suite.Equal(suite.mustSerializeKey(key), id)
	suite.Equal(id, whitelist.Id)
	return whitelist
}

func (suite *ProcessWhitelistDataStoreTestSuite) createAndStoreWhitelists(keys ...*storage.ProcessWhitelistKey) []*storage.ProcessWhitelist {
	whitelists := make([]*storage.ProcessWhitelist, len(keys))
	for i, key := range keys {
		whitelists[i] = suite.createAndStoreWhitelist(key)
	}
	return whitelists
}

func (suite *ProcessWhitelistDataStoreTestSuite) createAndStoreWhitelistWithRandomKey() *storage.ProcessWhitelist {
	return suite.createAndStoreWhitelist(&storage.ProcessWhitelistKey{
		DeploymentId:  uuid.NewV4().String(),
		ContainerName: uuid.NewV4().String(),
		ClusterId:     uuid.NewV4().String(),
		Namespace:     uuid.NewV4().String(),
	})
}

func (suite *ProcessWhitelistDataStoreTestSuite) doGet(key *storage.ProcessWhitelistKey, exists bool, equals *storage.ProcessWhitelist) *storage.ProcessWhitelist {
	whitelist, err := suite.datastore.GetProcessWhitelist(suite.requestContext, key)
	suite.NoError(err)
	if exists {
		suite.NotNil(whitelist)
		if equals != nil {
			suite.Equal(equals, whitelist)
		}
	} else {
		suite.Nil(whitelist)
	}
	return whitelist
}

func (suite *ProcessWhitelistDataStoreTestSuite) testUpdate(key *storage.ProcessWhitelistKey, addProcesses []string, removeProcesses []string, auto bool, expectedResults set.StringSet) *storage.ProcessWhitelist {
	updated, err := suite.datastore.UpdateProcessWhitelistElements(suite.requestContext, key, fixtures.MakeWhitelistItems(addProcesses...), fixtures.MakeWhitelistItems(removeProcesses...), auto)
	suite.NoError(err)
	suite.NotNil(updated)
	suite.True(updated.GetLastUpdate().Compare(updated.GetCreated()) > 0)
	suite.NotNil(updated.Elements)
	suite.Equal(expectedResults.Cardinality(), len(updated.Elements))
	actualResults := set.NewStringSet()
	for _, process := range updated.Elements {
		actualResults.Add(process.GetElement().GetProcessName())
	}
	suite.Equal(expectedResults, actualResults)
	return updated
}

func (suite *ProcessWhitelistDataStoreTestSuite) TestGetById() {
	suite.doGet(&storage.ProcessWhitelistKey{DeploymentId: "FAKE", ContainerName: "whatever", ClusterId: "whatever", Namespace: "whatever"}, false, nil)

	key := &storage.ProcessWhitelistKey{
		DeploymentId:  "blah",
		ContainerName: "container",
		ClusterId:     "cluster1",
		Namespace:     "namespace",
	}
	whitelist := suite.createAndStoreWhitelist(key)
	suite.doGet(key, true, whitelist)
}

func (suite *ProcessWhitelistDataStoreTestSuite) TestRemoveProcessWhitelist() {
	whitelist := suite.createAndStoreWhitelistWithRandomKey()
	key := whitelist.GetKey()
	suite.doGet(whitelist.GetKey(), true, whitelist)
	suite.whitelistResultsStore.EXPECT().DeleteWhitelistResults(suite.requestContext, key.GetDeploymentId()).Return(nil)
	err := suite.datastore.RemoveProcessWhitelist(suite.requestContext, key)
	suite.NoError(err)
	suite.doGet(key, false, nil)
}

func (suite *ProcessWhitelistDataStoreTestSuite) TestLockAndUnlockWhitelist() {
	whitelist := suite.createAndStoreWhitelistWithRandomKey()
	key := whitelist.GetKey()
	suite.Nil(whitelist.GetUserLockedTimestamp())
	updatedWhitelist, err := suite.datastore.UserLockProcessWhitelist(suite.requestContext, key, true)
	suite.NoError(err)
	suite.NotNil(updatedWhitelist.GetUserLockedTimestamp())
	suite.doGet(key, true, updatedWhitelist)
	suite.True(updatedWhitelist.GetLastUpdate().Compare(updatedWhitelist.GetCreated()) > 0)

	updatedWhitelist, err = suite.datastore.UserLockProcessWhitelist(suite.requestContext, key, false)
	suite.NoError(err)
	suite.Nil(updatedWhitelist.GetUserLockedTimestamp())
	suite.doGet(key, true, updatedWhitelist)
	suite.True(updatedWhitelist.GetLastUpdate().Compare(updatedWhitelist.GetCreated()) > 0)
}

func (suite *ProcessWhitelistDataStoreTestSuite) TestUpdateProcessWhitelist() {
	whitelist := fixtures.GetProcessWhitelistWithKey()
	whitelist.Elements = nil // Fixture gives a single process but we want to test updates
	suite.NotNil(whitelist)
	key := whitelist.GetKey()
	id, err := suite.datastore.AddProcessWhitelist(suite.requestContext, whitelist)
	suite.NoError(err)
	suite.NotNil(id)

	processName := []string{"Some process name"}
	processNameSet := set.NewStringSet(processName...)
	otherProcess := []string{"Some other process"}
	otherProcessSet := set.NewStringSet(otherProcess...)
	updated := suite.testUpdate(key, processName, nil, true, processNameSet)
	suite.True(updated.Elements[0].Auto)

	updated = suite.testUpdate(key, processName, nil, false, processNameSet)
	suite.False(updated.Elements[0].Auto)

	updated = suite.testUpdate(key, otherProcess, processName, true, otherProcessSet)
	suite.True(updated.Elements[0].Auto)

	multiAdd := []string{"a", "b", "c"}
	multiAddExpected := set.NewStringSet(multiAdd...)
	updated = suite.testUpdate(key, multiAdd, otherProcess, false, multiAddExpected)
	for _, process := range updated.Elements {
		suite.False(process.Auto)
	}
}

func (suite *ProcessWhitelistDataStoreTestSuite) TestUpsertProcessWhitelist() {
	key := fixtures.GetWhitelistKey()
	firstProcess := "Joseph Rules"
	newItem := []*storage.WhitelistItem{{Item: &storage.WhitelistItem_ProcessName{ProcessName: firstProcess}}}
	whitelist, err := suite.datastore.UpsertProcessWhitelist(suite.requestContext, key, newItem, true)
	suite.NoError(err)
	suite.Equal(1, len(whitelist.GetElements()))
	suite.Equal(firstProcess, whitelist.GetElements()[0].GetElement().GetProcessName())
	suite.Equal(key, whitelist.GetKey())
	suite.True(whitelist.GetLastUpdate().Compare(whitelist.GetCreated()) == 0)

	secondProcess := "Joseph is the Best"
	newItem = []*storage.WhitelistItem{{Item: &storage.WhitelistItem_ProcessName{ProcessName: secondProcess}}}
	whitelist, err = suite.datastore.UpsertProcessWhitelist(suite.requestContext, key, newItem, true)
	suite.NoError(err)
	suite.Equal(2, len(whitelist.GetElements()))
	processNames := make([]string, 0, 2)
	for _, element := range whitelist.GetElements() {
		processNames = append(processNames, element.GetElement().GetProcessName())
	}
	suite.ElementsMatch([]string{firstProcess, secondProcess}, processNames)
	suite.Equal(key, whitelist.GetKey())
	suite.True(whitelist.GetLastUpdate().Compare(whitelist.GetCreated()) > 0)
}

func makeItemList(elementList []*storage.WhitelistElement) []*storage.WhitelistItem {
	itemList := make([]*storage.WhitelistItem, len(elementList))
	for i, element := range elementList {
		itemList[i] = element.GetElement()
	}
	return itemList
}

func (suite *ProcessWhitelistDataStoreTestSuite) TestGraveyard() {
	whitelist := suite.createAndStoreWhitelistWithRandomKey()
	itemList := makeItemList(whitelist.GetElements())
	suite.NotEmpty(itemList)
	suite.Empty(whitelist.GetElementGraveyard())
	updatedWhitelist, err := suite.datastore.UpdateProcessWhitelistElements(suite.requestContext, whitelist.GetKey(), nil, itemList, true)
	// The elements should have been removed from the whitelist and put in the graveyard
	suite.NoError(err)
	suite.ElementsMatch(whitelist.GetElements(), updatedWhitelist.GetElementGraveyard())

	updatedWhitelist, err = suite.datastore.UpdateProcessWhitelistElements(suite.requestContext, whitelist.GetKey(), itemList, nil, true)
	suite.NoError(err)
	// The elements should NOT be added back on to the whitelist because they are in the graveyard and auto = true
	suite.Empty(updatedWhitelist.GetElements())
	suite.ElementsMatch(whitelist.GetElements(), updatedWhitelist.GetElementGraveyard())

	updatedWhitelist, err = suite.datastore.UpdateProcessWhitelistElements(suite.requestContext, whitelist.GetKey(), itemList, nil, false)
	suite.NoError(err)
	// The elements SHOULD be added back on to the whitelist because auto = false
	suite.Empty(updatedWhitelist.GetElementGraveyard())
	updatedItems := makeItemList(updatedWhitelist.GetElements())
	suite.ElementsMatch(itemList, updatedItems)
}

func (suite *ProcessWhitelistDataStoreTestSuite) doQuery(q *v1.Query, len int) {
	result, err := suite.datastore.SearchRawProcessWhitelists(suite.requestContext, q)
	suite.NoError(err)
	suite.Len(result, len)
}

func (suite *ProcessWhitelistDataStoreTestSuite) TestRemoveByDeployment() {
	dep1 := "1"
	key1 := &storage.ProcessWhitelistKey{DeploymentId: dep1, ContainerName: "1", ClusterId: "1", Namespace: "1"}
	key2 := &storage.ProcessWhitelistKey{DeploymentId: dep1, ContainerName: "2", ClusterId: "1", Namespace: "2"}
	key3 := &storage.ProcessWhitelistKey{DeploymentId: "2", ContainerName: "1", ClusterId: "1", Namespace: "3"}
	suite.createAndStoreWhitelists(key1, key2, key3)

	queryDep1 := pkgSearch.NewQueryBuilder().AddExactMatches(pkgSearch.DeploymentID, dep1).ProtoQuery()
	queryDep2 := pkgSearch.NewQueryBuilder().AddExactMatches(pkgSearch.DeploymentID, "2").ProtoQuery()
	suite.doQuery(queryDep1, 2)
	suite.doQuery(queryDep2, 1)
	suite.doGet(key1, true, nil)
	suite.doGet(key2, true, nil)
	suite.doGet(key3, true, nil)

	suite.whitelistResultsStore.EXPECT().DeleteWhitelistResults(suite.requestContext, dep1).Return(nil)
	err := suite.datastore.RemoveProcessWhitelistsByDeployment(suite.requestContext, dep1)
	suite.NoError(err)

	suite.doQuery(queryDep1, 0)
	suite.doQuery(queryDep2, 1)
	suite.doGet(key1, false, nil)
	suite.doGet(key2, false, nil)
	suite.doGet(key3, true, nil)
}

func (suite *ProcessWhitelistDataStoreTestSuite) TestIDToKeyConversion() {
	key := &storage.ProcessWhitelistKey{
		DeploymentId:  "blah",
		ContainerName: "container",
		ClusterId:     "cluster1",
		Namespace:     "namespace",
	}

	id, err := keyToID(key)
	suite.NoError(err)
	resKey, err := IDToKey(id)
	suite.NoError(err)
	suite.NotNil(resKey)
	suite.Equal(*key, *resKey)
}
