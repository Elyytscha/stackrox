package reprocessor

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stackrox/rox/central/deployment/datastore/mocks"
	connectionMocks "github.com/stackrox/rox/central/sensor/service/connection/mocks"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stretchr/testify/suite"
)

func TestLoop(t *testing.T) {
	suite.Run(t, new(loopTestSuite))
}

type loopTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller

	mockManager    *connectionMocks.MockManager
	mockDeployment *mocks.MockDataStore
}

func (suite *loopTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockManager = connectionMocks.NewMockManager(suite.mockCtrl)
	suite.mockDeployment = mocks.NewMockDataStore(suite.mockCtrl)
}

func (suite *loopTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *loopTestSuite) expectCalls(times int, allowMore bool) {
	timesSpec := (*gomock.Call).Times
	if allowMore {
		timesSpec = (*gomock.Call).MinTimes
	}
	query := search.NewQueryBuilder().AddStringsHighlighted(search.ClusterID, search.WildcardString).ProtoQuery()
	timesSpec(suite.mockDeployment.EXPECT().SearchDeployments(getDeploymentsContext, query), times).Return(nil, nil)
}

func (suite *loopTestSuite) TestTimerDoesNotTick() {
	loop := NewLoop(suite.mockManager, suite.mockDeployment)
	loop.Start()
	loop.Stop()
	suite.mockManager.EXPECT().GetActiveConnections().MaxTimes(0)
}

func (suite *loopTestSuite) TestTimerTicksOnce() {
	duration := 1 * time.Second // Need this to be long enough that the enrichAndDetectTicker won't get called twice during the test.
	loop := newLoopWithDuration(suite.mockManager, suite.mockDeployment, duration, duration, duration)
	suite.expectCalls(1, false)
	loop.Start()
	time.Sleep(duration + 10*time.Millisecond)
	loop.Stop()
}

func (suite *loopTestSuite) TestTimerTicksTwice() {
	duration := 100 * time.Millisecond
	loop := newLoopWithDuration(suite.mockManager, suite.mockDeployment, duration, duration, duration)
	suite.expectCalls(2, true)
	loop.Start()
	time.Sleep((2 * duration) + (10 * time.Millisecond))
	loop.Stop()
}

func (suite *loopTestSuite) TestShortCircuitOnce() {
	loop := NewLoop(suite.mockManager, suite.mockDeployment)
	suite.expectCalls(1, false)
	loop.Start()
	go loop.ShortCircuit()
	// Sleep for a little bit of time to allow the mock calls to go through, since they happen asynchronously.
	time.Sleep(500 * time.Millisecond)
	loop.Stop()
}

func (suite *loopTestSuite) TestShortCircuitTwice() {
	loop := NewLoop(suite.mockManager, suite.mockDeployment)
	suite.expectCalls(2, false)
	loop.Start()
	go loop.ShortCircuit()
	go loop.ShortCircuit()
	// Sleep for a little bit of time to allow the mock calls to go through, since they happen asynchronously.
	time.Sleep(500 * time.Millisecond)
	loop.Stop()
}

func (suite *loopTestSuite) TestStopWorks() {
	loop := NewLoop(suite.mockManager, suite.mockDeployment)
	suite.expectCalls(1, false)
	loop.Start()
	go loop.ShortCircuit()
	time.Sleep(500 * time.Millisecond)
	loop.Stop()
	time.Sleep(100 * time.Millisecond)
	go loop.ShortCircuit()
	// Sleep for a little bit of time to allow the mock calls to go through, since they happen asynchronously.
	time.Sleep(500 * time.Millisecond)
}
