package counter

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/dackbox/graph"
	"github.com/stackrox/rox/pkg/dackbox/graph/mocks"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search/filtered"
	"github.com/stretchr/testify/suite"
)

var (
	prefix1       = []byte("pre1")
	prefix2       = []byte("pre2")
	prefix3       = []byte("pre3")
	clusterPrefix = []byte("cluster")

	id1 = []byte("id1")
	id2 = []byte("id2")
	id3 = []byte("id3")
	id4 = []byte("id4")
	id5 = []byte("id5")
	id6 = []byte("id6")
	id7 = []byte("id7")

	prefixedID1 = badgerhelper.GetBucketKey(prefix1, id1)
	prefixedID2 = badgerhelper.GetBucketKey(prefix2, id2)
	prefixedID3 = badgerhelper.GetBucketKey(prefix2, id3)
	prefixedID4 = badgerhelper.GetBucketKey(prefix3, id4)
	prefixedID5 = badgerhelper.GetBucketKey(prefix3, id5)
	prefixedID6 = badgerhelper.GetBucketKey(clusterPrefix, id6)
	prefixedID7 = badgerhelper.GetBucketKey(clusterPrefix, id7)

	// Fake hierarchy for test, use prefixed values since that is what will be stored in the graph.
	fromID1 = [][]byte{prefixedID2, prefixedID3}
	fromID2 = [][]byte{prefixedID4}
	fromID3 = [][]byte{prefixedID5}

	toID4 = [][]byte{prefixedID6}
	toID5 = [][]byte{prefixedID7}

	globalResource = permissions.ResourceMetadata{
		Resource: "resource",
		Scope:    permissions.GlobalScope,
	}
	clusterResource = permissions.ResourceMetadata{
		Resource: "resource",
		Scope:    permissions.ClusterScope,
	}
)

func TestDerivedFieldCounter(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(derivedFieldCounterTestSuite))
}

type derivedFieldCounterTestSuite struct {
	suite.Suite

	mockRGraph *mocks.MockRGraph

	mockCtrl *gomock.Controller
}

func (s *derivedFieldCounterTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockRGraph = mocks.NewMockRGraph(s.mockCtrl)
}

func (s *derivedFieldCounterTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

/*
id1 -> id2 -> id4
id1 -> id3 -> id5
*/
func (s *derivedFieldCounterTestSuite) TestCounterForward() {
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID1).Return(fromID1)
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID2).Return(fromID2)
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID3).Return(fromID3)

	ctx := sac.WithGlobalAccessScopeChecker(context.Background(), sac.AllowAllAccessScopeChecker())

	filter, err := filtered.NewSACFilter(
		filtered.WithResourceHelper(sac.ForResource(globalResource)),
	)
	s.NoError(err, "filter creation should have succeeded")

	counter := NewGraphBasedDerivedFieldCounter(fakeGraphProvider{mg: s.mockRGraph}, [][]byte{prefix1, prefix2, prefix3}, filter, true)
	count, _ := counter.Count(ctx, string(id1))
	s.Equal(map[string]int32{string(id1): int32(2)}, count)
}

/*
id1 -> id2 -> id4
id1 -> id3
*/
func (s *derivedFieldCounterTestSuite) TestCounterForwardWithPartialPath() {
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID1).Return(fromID1)
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID2).Return(fromID2)
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID3).Return([][]byte{})

	ctx := sac.WithGlobalAccessScopeChecker(context.Background(), sac.AllowAllAccessScopeChecker())

	filter, err := filtered.NewSACFilter(
		filtered.WithResourceHelper(sac.ForResource(globalResource)),
	)
	s.NoError(err, "filter creation should have succeeded")

	counter := NewGraphBasedDerivedFieldCounter(fakeGraphProvider{mg: s.mockRGraph}, [][]byte{prefix1, prefix2, prefix3}, filter, true)
	count, _ := counter.Count(ctx, string(id1))
	s.Equal(map[string]int32{string(id1): int32(1)}, count)
}

/*
id1 -> id2 -> id4
id1 -> id3 -> id5 (no access)
*/
func (s *derivedFieldCounterTestSuite) TestCounterForwardWithSACFilter() {
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID1).Return(fromID1)
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID2).Return(fromID2)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID4).Return(toID4)
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID3).Return(fromID3)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID5).Return(toID5)

	graphProvider := fakeGraphProvider{mg: s.mockRGraph}

	ctx := sac.WithGlobalAccessScopeChecker(context.Background(), sac.AllowFixedScopes(
		sac.AccessModeScopeKeys(storage.Access_READ_ACCESS),
		sac.ResourceScopeKeys(clusterResource),
		sac.ClusterScopeKeys("id6")))

	filter, err := filtered.NewSACFilter(
		filtered.WithResourceHelper(sac.ForResource(clusterResource)),
		filtered.WithGraphProvider(graphProvider),
		filtered.WithClusterPath(prefix3, clusterPrefix),
	)
	s.NoError(err, "filter creation should have succeeded")

	counter := NewGraphBasedDerivedFieldCounter(graphProvider, [][]byte{prefix1, prefix2, prefix3}, filter, true)
	count, _ := counter.Count(ctx, string(id1))
	s.Equal(map[string]int32{string(id1): int32(1)}, count)
}

/*
id1 -> id2 -> id5
id1 -> id3 -> id5
*/
func (s *derivedFieldCounterTestSuite) TestCounterForwardRepeated() {
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID1).Return(fromID1)
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID2).Return([][]byte{prefixedID5})
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID3).Return(fromID3)

	graphProvider := fakeGraphProvider{mg: s.mockRGraph}
	ctx := sac.WithGlobalAccessScopeChecker(context.Background(), sac.AllowAllAccessScopeChecker())
	filter, err := filtered.NewSACFilter(
		filtered.WithResourceHelper(sac.ForResource(globalResource)),
	)
	s.NoError(err, "filter creation should have succeeded")

	counter := NewGraphBasedDerivedFieldCounter(graphProvider, [][]byte{prefix1, prefix2, prefix3}, filter, true)
	count, _ := counter.Count(ctx, string(id1))
	s.Equal(map[string]int32{string(id1): int32(1)}, count)
}

/*
id1 -> id2 -> id4
id1 -> id3 -> id5
id1 -> id3 -> id6
*/
func (s *derivedFieldCounterTestSuite) TestCounterForwardOneToMany() {
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID1).Return(fromID1)
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID2).Return(fromID2)
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID3).Return([][]byte{prefixedID5, badgerhelper.GetBucketKey(prefix3, id6)})

	graphProvider := fakeGraphProvider{mg: s.mockRGraph}
	ctx := sac.WithGlobalAccessScopeChecker(context.Background(), sac.AllowAllAccessScopeChecker())
	filter, err := filtered.NewSACFilter(
		filtered.WithResourceHelper(sac.ForResource(globalResource)),
	)
	s.NoError(err, "filter creation should have succeeded")

	counter := NewGraphBasedDerivedFieldCounter(graphProvider, [][]byte{prefix1, prefix2, prefix3}, filter, true)
	count, _ := counter.Count(ctx, string(id1))
	s.Equal(map[string]int32{string(id1): int32(3)}, count)
}

/*
id1 -> id2 -> id4
id1 -> id3 -> id5
id1 -> id3 -> id6 (diff prefix)
*/
func (s *derivedFieldCounterTestSuite) TestCounterForwardWithDiffPrefix() {
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID1).Return(fromID1)
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID2).Return(fromID2)
	s.mockRGraph.EXPECT().GetRefsFrom(prefixedID3).Return([][]byte{prefixedID5, prefixedID6})

	graphProvider := fakeGraphProvider{mg: s.mockRGraph}
	ctx := sac.WithGlobalAccessScopeChecker(context.Background(), sac.AllowAllAccessScopeChecker())
	filter, err := filtered.NewSACFilter(
		filtered.WithResourceHelper(sac.ForResource(globalResource)),
	)
	s.NoError(err, "filter creation should have succeeded")

	counter := NewGraphBasedDerivedFieldCounter(graphProvider, [][]byte{prefix1, prefix2, prefix3}, filter, true)
	count, _ := counter.Count(ctx, string(id1))
	s.Equal(map[string]int32{string(id1): int32(2)}, count)
}

type fakeGraphProvider struct {
	mg *mocks.MockRGraph
}

func (fgp fakeGraphProvider) NewGraphView() graph.DiscardableRGraph {
	return graph.NewDiscardableGraph(fgp.mg, func() {})
}
