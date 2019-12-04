package badger

import (
	"os"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/stackrox/rox/central/networkflow/datastore/internal/store"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestFlowStore(t *testing.T) {
	if features.ManagedDB.Enabled() {
		t.Skip()
	}
	suite.Run(t, new(FlowStoreTestSuite))
}

type FlowStoreTestSuite struct {
	suite.Suite

	path   string
	db     *badger.DB
	tested store.FlowStore
}

func (suite *FlowStoreTestSuite) SetupSuite() {
	db, path, err := badgerhelper.NewTemp(suite.T().Name(), false)
	suite.Require().NoError(err)
	suite.db = db
	suite.path = path

	cs := NewClusterStore(suite.db)
	suite.tested, err = cs.CreateFlowStore("fakecluster")
	suite.Require().NoError(err)
}

func (suite *FlowStoreTestSuite) TearDownSuite() {
	_ = suite.db.Close()
	_ = os.RemoveAll(suite.path)
}

func (suite *FlowStoreTestSuite) TestStore() {
	t1 := time.Now().Add(-5 * time.Minute)
	t2 := time.Now()
	flows := []*storage.NetworkFlow{
		{
			Props: &storage.NetworkFlowProperties{
				SrcEntity:  &storage.NetworkEntityInfo{Type: storage.NetworkEntityInfo_DEPLOYMENT, Id: "someNode1"},
				DstEntity:  &storage.NetworkEntityInfo{Type: storage.NetworkEntityInfo_DEPLOYMENT, Id: "someNode2"},
				DstPort:    1,
				L4Protocol: storage.L4Protocol_L4_PROTOCOL_TCP,
			},
			LastSeenTimestamp: protoconv.ConvertTimeToTimestamp(t1),
		},
		{
			Props: &storage.NetworkFlowProperties{
				SrcEntity:  &storage.NetworkEntityInfo{Type: storage.NetworkEntityInfo_DEPLOYMENT, Id: "someOtherNode1"},
				DstEntity:  &storage.NetworkEntityInfo{Type: storage.NetworkEntityInfo_DEPLOYMENT, Id: "someOtherNode2"},
				DstPort:    2,
				L4Protocol: storage.L4Protocol_L4_PROTOCOL_TCP,
			},
			LastSeenTimestamp: protoconv.ConvertTimeToTimestamp(t2),
		},
		{
			Props: &storage.NetworkFlowProperties{
				SrcEntity:  &storage.NetworkEntityInfo{Type: storage.NetworkEntityInfo_DEPLOYMENT, Id: "yetAnotherNode1"},
				DstEntity:  &storage.NetworkEntityInfo{Type: storage.NetworkEntityInfo_DEPLOYMENT, Id: "yetAnotherNode2"},
				DstPort:    3,
				L4Protocol: storage.L4Protocol_L4_PROTOCOL_TCP,
			},
		},
	}
	var err error

	updateTS := timestamp.Now() - 1000000
	err = suite.tested.UpsertFlows(flows, updateTS)
	suite.NoError(err, "upsert should succeed on first insert")

	readFlows, readUpdateTS, err := suite.tested.GetAllFlows(nil)
	suite.Require().NoError(err)
	suite.ElementsMatch(readFlows, flows)
	suite.Equal(updateTS, timestamp.FromProtobuf(&readUpdateTS))

	readFlows, readUpdateTS, err = suite.tested.GetAllFlows(protoconv.ConvertTimeToTimestamp(t2))
	suite.Require().NoError(err)
	suite.ElementsMatch(readFlows, flows[1:])
	suite.Equal(updateTS, timestamp.FromProtobuf(&readUpdateTS))

	readFlows, readUpdateTS, err = suite.tested.GetAllFlows(protoconv.ConvertTimeToTimestamp(time.Now()))
	suite.Require().NoError(err)
	suite.ElementsMatch(readFlows, flows[2:])
	suite.Equal(updateTS, timestamp.FromProtobuf(&readUpdateTS))

	updateTS += 1337
	err = suite.tested.UpsertFlows(flows, updateTS)
	suite.NoError(err, "upsert should succeed on second insert")

	err = suite.tested.RemoveFlow(&storage.NetworkFlowProperties{
		SrcEntity:  flows[0].GetProps().GetSrcEntity(),
		DstEntity:  flows[0].GetProps().GetDstEntity(),
		DstPort:    flows[0].GetProps().GetDstPort(),
		L4Protocol: flows[0].GetProps().GetL4Protocol(),
	})
	suite.NoError(err, "remove should succeed when present")

	err = suite.tested.RemoveFlow(&storage.NetworkFlowProperties{
		SrcEntity:  flows[0].GetProps().GetSrcEntity(),
		DstEntity:  flows[0].GetProps().GetDstEntity(),
		DstPort:    flows[0].GetProps().GetDstPort(),
		L4Protocol: flows[0].GetProps().GetL4Protocol(),
	})
	suite.NoError(err, "remove should succeed when not present")

	var actualFlows []*storage.NetworkFlow
	actualFlows, readUpdateTS, err = suite.tested.GetAllFlows(nil)
	suite.NoError(err)
	suite.ElementsMatch(actualFlows, flows[1:])
	suite.Equal(updateTS, timestamp.FromProtobuf(&readUpdateTS))

	updateTS += 42
	err = suite.tested.UpsertFlows(flows, updateTS)
	suite.NoError(err, "upsert should succeed")

	actualFlows, readUpdateTS, err = suite.tested.GetAllFlows(nil)
	suite.NoError(err)
	suite.ElementsMatch(actualFlows, flows)
	suite.Equal(updateTS, timestamp.FromProtobuf(&readUpdateTS))

	node1Flows, readUpdateTS, err := suite.tested.GetMatchingFlows(func(props *storage.NetworkFlowProperties) bool {
		if props.GetDstEntity().GetType() == storage.NetworkEntityInfo_DEPLOYMENT && props.GetDstEntity().GetId() == "someNode1" {
			return true
		}
		if props.GetSrcEntity().GetType() == storage.NetworkEntityInfo_DEPLOYMENT && props.GetSrcEntity().GetId() == "someNode1" {
			return true
		}
		return false
	}, nil)
	suite.NoError(err)
	suite.ElementsMatch(node1Flows, flows[:1])
	suite.Equal(updateTS, timestamp.FromProtobuf(&readUpdateTS))
}

func (suite *FlowStoreTestSuite) TestRemoveAllMatching() {
	t1 := time.Now().Add(-5 * time.Minute)
	t2 := time.Now()
	flows := []*storage.NetworkFlow{
		{
			Props: &storage.NetworkFlowProperties{
				SrcEntity:  &storage.NetworkEntityInfo{Type: storage.NetworkEntityInfo_DEPLOYMENT, Id: "someNode1"},
				DstEntity:  &storage.NetworkEntityInfo{Type: storage.NetworkEntityInfo_DEPLOYMENT, Id: "someNode2"},
				DstPort:    1,
				L4Protocol: storage.L4Protocol_L4_PROTOCOL_TCP,
			},
			LastSeenTimestamp: protoconv.ConvertTimeToTimestamp(t1),
		},
		{
			Props: &storage.NetworkFlowProperties{
				SrcEntity:  &storage.NetworkEntityInfo{Type: storage.NetworkEntityInfo_DEPLOYMENT, Id: "someOtherNode1"},
				DstEntity:  &storage.NetworkEntityInfo{Type: storage.NetworkEntityInfo_DEPLOYMENT, Id: "someOtherNode2"},
				DstPort:    2,
				L4Protocol: storage.L4Protocol_L4_PROTOCOL_TCP,
			},
			LastSeenTimestamp: protoconv.ConvertTimeToTimestamp(t2),
		},
		{
			Props: &storage.NetworkFlowProperties{
				SrcEntity:  &storage.NetworkEntityInfo{Type: storage.NetworkEntityInfo_DEPLOYMENT, Id: "yetAnotherNode1"},
				DstEntity:  &storage.NetworkEntityInfo{Type: storage.NetworkEntityInfo_DEPLOYMENT, Id: "yetAnotherNode2"},
				DstPort:    3,
				L4Protocol: storage.L4Protocol_L4_PROTOCOL_TCP,
			},
		},
	}
	updateTS := timestamp.Now() - 1000000
	err := suite.tested.UpsertFlows(flows, updateTS)
	suite.NoError(err)

	// Match none delete none
	err = suite.tested.RemoveMatchingFlows(func(props *storage.NetworkFlowProperties) bool {
		return false
	}, nil)
	suite.NoError(err)

	currFlows, _, err := suite.tested.GetAllFlows(nil)
	suite.NoError(err)
	suite.ElementsMatch(flows, currFlows)

	// Match dst port 1
	err = suite.tested.RemoveMatchingFlows(func(props *storage.NetworkFlowProperties) bool {
		return props.DstPort == 1
	}, nil)
	suite.NoError(err)

	currFlows, _, err = suite.tested.GetAllFlows(nil)
	suite.NoError(err)
	suite.ElementsMatch(flows[1:], currFlows)

	err = suite.tested.RemoveMatchingFlows(nil, func(flow *storage.NetworkFlow) bool {
		return flow.LastSeenTimestamp.Compare(protoconv.ConvertTimeToTimestamp(t2)) == 0
	})
	suite.NoError(err)

	currFlows, _, err = suite.tested.GetAllFlows(nil)
	suite.NoError(err)
	suite.ElementsMatch(flows[2:], currFlows)
}

func TestGetDeploymentIDsFromKey(t *testing.T) {
	s := &flowStoreImpl{}
	id := s.getID(&storage.NetworkFlowProperties{
		SrcEntity: &storage.NetworkEntityInfo{
			Type: storage.NetworkEntityInfo_DEPLOYMENT,
			Id:   "id1",
		},
		DstEntity: &storage.NetworkEntityInfo{
			Type: storage.NetworkEntityInfo_INTERNET,
			Id:   "id2",
		},
		DstPort:    8080,
		L4Protocol: storage.L4Protocol_L4_PROTOCOL_TCP,
	})

	id1, id2 := s.getDeploymentIDsFromKey(id)
	assert.Equal(t, []byte("id1"), id1)
	assert.Equal(t, []byte("id2"), id2)
}
