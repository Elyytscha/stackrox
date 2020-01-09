package search

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mockIndex "github.com/stackrox/rox/central/processwhitelist/index/mocks"
	mockStore "github.com/stackrox/rox/central/processwhitelist/store/mocks"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/fixtures"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestProcessWhitelistSearch(t *testing.T) {
	suite.Run(t, new(ProcessWhitelistSearchTestSuite))
}

func getFakeSearchResults(num int) ([]search.Result, []*storage.ProcessWhitelist) {
	var dbResults []*storage.ProcessWhitelist
	var indexResults []search.Result
	for i := 0; i < num; i++ {
		whitelist := fixtures.GetProcessWhitelistWithID()
		dbResults = append(dbResults, whitelist)
		fakeResult := search.Result{ID: whitelist.Id}
		indexResults = append(indexResults, fakeResult)
	}

	return indexResults, dbResults
}

type ProcessWhitelistSearchTestSuite struct {
	suite.Suite

	controller *gomock.Controller
	indexer    *mockIndex.MockIndexer
	store      *mockStore.MockStore

	searcher    Searcher
	allowAllCtx context.Context
}

func (suite *ProcessWhitelistSearchTestSuite) SetupTest() {
	suite.allowAllCtx = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.ProcessWhitelist),
		))
	suite.controller = gomock.NewController(suite.T())
	suite.indexer = mockIndex.NewMockIndexer(suite.controller)
	suite.store = mockStore.NewMockStore(suite.controller)
	suite.store.EXPECT().WalkAll(gomock.Any()).Return(nil)
	searcher, err := New(suite.store, suite.indexer)

	suite.NoError(err)
	suite.searcher = searcher
}

func (suite *ProcessWhitelistSearchTestSuite) TearDownTest() {
	suite.controller.Finish()
}

func (suite *ProcessWhitelistSearchTestSuite) TestErrors() {
	q := search.EmptyQuery()
	someError := errors.New("this is a test error")
	suite.indexer.EXPECT().Search(q).Return(nil, someError)
	results, err := suite.searcher.SearchRawProcessWhitelists(suite.allowAllCtx, q)
	suite.Equal(someError, err)
	suite.Nil(results)

	indexResults, _ := getFakeSearchResults(1)
	suite.indexer.EXPECT().Search(q).Return(indexResults, nil)
	suite.store.EXPECT().GetWhitelists(search.ResultsToIDs(indexResults)).Return(nil, nil, someError)
	results, err = suite.searcher.SearchRawProcessWhitelists(suite.allowAllCtx, q)
	suite.Error(err)
	suite.Nil(results)
}

func (suite *ProcessWhitelistSearchTestSuite) TestSearchForAll() {
	q := search.EmptyQuery()
	var emptyList []search.Result
	suite.indexer.EXPECT().Search(q).Return(emptyList, nil)
	// It's an implementation detail whether this method is called, so allow but don't require it.
	suite.store.EXPECT().GetWhitelists(testutils.AssertionMatcher(assert.Empty)).MinTimes(0).MaxTimes(1)
	results, err := suite.searcher.SearchRawProcessWhitelists(suite.allowAllCtx, q)
	suite.NoError(err)
	suite.Empty(results)

	indexResults, dbResults := getFakeSearchResults(3)
	suite.indexer.EXPECT().Search(q).Return(indexResults, nil)
	suite.store.EXPECT().GetWhitelists(search.ResultsToIDs(indexResults)).Return(dbResults, nil, nil)
	results, err = suite.searcher.SearchRawProcessWhitelists(suite.allowAllCtx, q)
	suite.NoError(err)
	suite.Equal(dbResults, results)
}
