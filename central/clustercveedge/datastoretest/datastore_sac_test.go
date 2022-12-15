package datastore

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stackrox/rox/central/clustercveedge/datastore"
	cveConverterV1 "github.com/stackrox/rox/central/cve/converter"
	"github.com/stackrox/rox/central/cve/converter/utils"
	dackboxTestUtils "github.com/stackrox/rox/central/dackbox/testutils"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/dackbox/edges"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/fixtures"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/sac/testconsts"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stretchr/testify/suite"
)

const (
	clusterOS = ""

	waitForIndexing     = true
	dontWaitForIndexing = false

	cluster1ToCVE1EdgeID = "cluster1ToCVE1EdgeID"
	cluster1ToCVE2EdgeID = "cluster1ToCVE2EdgeID"
	cluster2ToCVE2EdgeID = "cluster2ToCVE2EdgeID"
	cluster2ToCVE3EdgeID = "cluster2ToCVE3EdgeID"
)

var (
	allAccessCtx = sac.WithAllAccess(context.Background())
)

func TestClusterCVEEdgeDatastoreSAC(t *testing.T) {
	suite.Run(t, new(clusterCVEEdgeDatastoreSACSuite))
}

type clusterCVEEdgeDatastoreSACSuite struct {
	suite.Suite

	datastore        datastore.DataStore
	dackboxTestStore dackboxTestUtils.DackboxTestDataStore
}

func (s *clusterCVEEdgeDatastoreSACSuite) SetupSuite() {
	var err error
	s.dackboxTestStore, err = dackboxTestUtils.NewDackboxTestDataStore(s.T())
	s.Require().NoError(err)
	if env.PostgresDatastoreEnabled.BooleanSetting() {
		s.datastore, err = datastore.GetTestPostgresDataStore(s.T(), s.dackboxTestStore.GetPostgresPool())
		s.Require().NoError(err)
	} else {
		s.datastore, err = datastore.GetTestRocksBleveDataStore(
			s.T(),
			s.dackboxTestStore.GetRocksEngine(),
			s.dackboxTestStore.GetBleveIndex(),
			s.dackboxTestStore.GetDackbox(),
			s.dackboxTestStore.GetKeyFence(),
		)
		s.Require().NoError(err)
	}
}

func (s *clusterCVEEdgeDatastoreSACSuite) TearDownSuite() {
	s.Require().NoError(s.dackboxTestStore.Cleanup(s.T()))
}

func (s *clusterCVEEdgeDatastoreSACSuite) cleanImageToVulnerabilitiesGraph(waitForIndexing bool) {
	s.Require().NoError(s.dackboxTestStore.CleanClusterToVulnerabilitiesGraph(waitForIndexing))
}

func getCveID(vulnerability *storage.EmbeddedVulnerability, os string) string {
	if env.PostgresDatastoreEnabled.BooleanSetting() {
		return vulnerability.GetCve()
	}
	return utils.EmbeddedCVEToProtoCVE(os, vulnerability).GetId()
}

func getEdgeID(vulnerability *storage.EmbeddedVulnerability, clusterID, os string) string {
	if env.PostgresDatastoreEnabled.BooleanSetting() {
		return clusterID + "#" + vulnerability.GetCve()
	}
	cve := utils.EmbeddedCVEToProtoCVE(os, vulnerability)
	return edges.EdgeID{ParentID: clusterID, ChildID: cve.GetId()}.ToString()
}

var (
	embeddedCVE1 = fixtures.GetEmbeddedClusterCVE1234x0001()
	embeddedCVE2 = fixtures.GetEmbeddedClusterCVE4567x0002()
	embeddedCVE3 = fixtures.GetEmbeddedClusterCVE2345x0003()
)

type testCase struct {
	name         string
	ctx          context.Context
	visibleEdges map[string]bool
}

func getClusterCVEEdges(cluster1, cluster2 string) map[string]string {
	return map[string]string{
		cluster1ToCVE1EdgeID: getEdgeID(embeddedCVE1, cluster1, clusterOS),
		cluster1ToCVE2EdgeID: getEdgeID(embeddedCVE2, cluster1, clusterOS),
		cluster2ToCVE2EdgeID: getEdgeID(embeddedCVE2, cluster2, clusterOS),
		cluster2ToCVE3EdgeID: getEdgeID(embeddedCVE3, cluster2, clusterOS),
	}
}

func getClusterCVEEdgeReadTestCases(_ *testing.T, validCluster1 string, validCluster2 string) []testCase {
	return []testCase{
		{
			name: "Full read-write access has access to all data",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: true,
				cluster1ToCVE2EdgeID: true,
				cluster2ToCVE2EdgeID: true,
				cluster2ToCVE3EdgeID: true,
			},
		},
		{
			name: "Full read-only access has read access to all data",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: true,
				cluster1ToCVE2EdgeID: true,
				cluster2ToCVE2EdgeID: true,
				cluster2ToCVE3EdgeID: true,
			},
		},
		{
			name: "Full cluster access has access to all data for the cluster",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
					sac.ClusterScopeKeys(validCluster1),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: true,
				cluster1ToCVE2EdgeID: true,
				cluster2ToCVE2EdgeID: false,
				cluster2ToCVE3EdgeID: false,
			},
		},
		{
			name: "Partial cluster access has access to no data",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
					sac.ClusterScopeKeys(validCluster1),
					sac.NamespaceScopeKeys(testconsts.NamespaceA),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: false,
				cluster1ToCVE2EdgeID: false,
				cluster2ToCVE2EdgeID: false,
				cluster2ToCVE3EdgeID: false,
			},
		},
		{
			name: "Full access to other cluster has access to all data for that cluster",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
					sac.ClusterScopeKeys(validCluster2),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: false,
				cluster1ToCVE2EdgeID: false,
				cluster2ToCVE2EdgeID: true,
				cluster2ToCVE3EdgeID: true,
			},
		},
		{
			name: "Partial access to other cluster has access to no data",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
					sac.ClusterScopeKeys(validCluster2),
					sac.NamespaceScopeKeys(testconsts.NamespaceB),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: false,
				cluster1ToCVE2EdgeID: false,
				cluster2ToCVE2EdgeID: false,
				cluster2ToCVE3EdgeID: false,
			},
		},
		{
			name: "Full access to wrong cluster has access to no data",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
					sac.ClusterScopeKeys(testconsts.WrongCluster),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: false,
				cluster1ToCVE2EdgeID: false,
				cluster2ToCVE2EdgeID: false,
				cluster2ToCVE3EdgeID: false,
			},
		},
	}
}

func getClusterCVEEdgeWriteTestCases(_ *testing.T, validCluster1 string, validCluster2 string) []testCase {
	return []testCase{
		{
			name: "Full read-write access has access to all data",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: true,
				cluster1ToCVE2EdgeID: true,
				cluster2ToCVE2EdgeID: true,
				cluster2ToCVE3EdgeID: true,
			},
		},
		{
			name: "Full read-only access has read access to all data",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: false,
				cluster1ToCVE2EdgeID: false,
				cluster2ToCVE2EdgeID: false,
				cluster2ToCVE3EdgeID: false,
			},
		},
		{
			name: "Full cluster access has access to all data for the cluster",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
					sac.ClusterScopeKeys(validCluster1),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: false,
				cluster1ToCVE2EdgeID: false,
				cluster2ToCVE2EdgeID: false,
				cluster2ToCVE3EdgeID: false,
			},
		},
		{
			name: "Partial cluster access has access to no data",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
					sac.ClusterScopeKeys(validCluster1),
					sac.NamespaceScopeKeys(testconsts.NamespaceA),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: false,
				cluster1ToCVE2EdgeID: false,
				cluster2ToCVE2EdgeID: false,
				cluster2ToCVE3EdgeID: false,
			},
		},
		{
			name: "Full access to other cluster has access to all data for that cluster",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
					sac.ClusterScopeKeys(validCluster2),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: false,
				cluster1ToCVE2EdgeID: false,
				cluster2ToCVE2EdgeID: false,
				cluster2ToCVE3EdgeID: false,
			},
		},
		{
			name: "Partial access to other cluster has access to no data",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
					sac.ClusterScopeKeys(validCluster2),
					sac.NamespaceScopeKeys(testconsts.NamespaceB),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: false,
				cluster1ToCVE2EdgeID: false,
				cluster2ToCVE2EdgeID: false,
				cluster2ToCVE3EdgeID: false,
			},
		},
		{
			name: "Full access to wrong cluster has access to no data",
			ctx: sac.WithGlobalAccessScopeChecker(
				context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
					sac.ResourceScopeKeys(resources.Cluster),
					sac.ClusterScopeKeys(testconsts.WrongCluster),
				),
			),
			visibleEdges: map[string]bool{
				cluster1ToCVE1EdgeID: false,
				cluster1ToCVE2EdgeID: false,
				cluster2ToCVE2EdgeID: false,
				cluster2ToCVE3EdgeID: false,
			},
		},
	}
}

func (s *clusterCVEEdgeDatastoreSACSuite) checkClusterCVEEdgePresence(id string, expectedPresent bool) {
	edge, found, err := s.datastore.Get(allAccessCtx, id)
	s.NoError(err)
	if expectedPresent {
		s.True(found)
		s.Equal(id, edge.GetId())
	} else {
		s.False(found)
		s.Nil(edge)
	}
}

func (s *clusterCVEEdgeDatastoreSACSuite) waitForIndexing() {
	if env.PostgresDatastoreEnabled.BooleanSetting() {
		return
	}
	deleteDoneSignal := concurrency.NewSignal()
	s.dackboxTestStore.GetIndexQ().PushSignal(&deleteDoneSignal)
	<-deleteDoneSignal.Done()
}

func (s *clusterCVEEdgeDatastoreSACSuite) TestUpsert() {
	if env.PostgresDatastoreEnabled.BooleanSetting() {
		s.T().Skip("ClusterCVEEdge Upsert is only available in non-postgres mode")
	}
	err := s.dackboxTestStore.PushClusterToVulnerabilitiesGraph(waitForIndexing)
	defer s.cleanImageToVulnerabilitiesGraph(waitForIndexing)
	s.Require().NoError(err)
	validClusters := s.dackboxTestStore.GetStoredClusterIDs()
	s.Require().True(len(validClusters) >= 2)

	testEdgeIDs := getClusterCVEEdges(validClusters[0], validClusters[1])
	testCases := getClusterCVEEdgeWriteTestCases(s.T(), validClusters[0], validClusters[1])
	embeddedCVE := fixtures.GetEmbeddedClusterCVE4567x0002()
	storageCVE := utils.EmbeddedCVEToProtoCVE(clusterOS, embeddedCVE)
	cveFixVersion := embeddedCVE.GetFixedBy()
	cluster1 := &storage.Cluster{
		Id:   validClusters[0],
		Name: validClusters[0],
		Type: storage.ClusterType_OPENSHIFT_CLUSTER,
	}
	cluster1Only := []*storage.Cluster{cluster1}
	cveParts := cveConverterV1.NewClusterCVEParts(storageCVE, cluster1Only, cveFixVersion)
	targetCveEdgeID := testEdgeIDs[cluster1ToCVE2EdgeID]
	err = s.datastore.Delete(allAccessCtx, targetCveEdgeID)
	s.Require().NoError(err)
	s.waitForIndexing()
	for _, c := range testCases {
		s.Run(c.name, func() {
			ctx := c.ctx
			s.checkClusterCVEEdgePresence(targetCveEdgeID, false)
			err = s.datastore.Upsert(ctx, cveParts)
			if c.visibleEdges[cluster1ToCVE2EdgeID] {
				s.NoError(err)
				s.waitForIndexing()
				s.checkClusterCVEEdgePresence(targetCveEdgeID, true)
			} else {
				s.ErrorIs(err, sac.ErrResourceAccessDenied)
				s.checkClusterCVEEdgePresence(targetCveEdgeID, false)
			}
			err = s.datastore.Delete(allAccessCtx, targetCveEdgeID)
			s.NoError(err)
			s.waitForIndexing()
			s.checkClusterCVEEdgePresence(targetCveEdgeID, false)
		})
	}
}

func (s *clusterCVEEdgeDatastoreSACSuite) TestDelete() {
	if env.PostgresDatastoreEnabled.BooleanSetting() {
		s.T().Skip("ClusterCVEEdge Delete is only available in non-postgres mode")
	}
	err := s.dackboxTestStore.PushClusterToVulnerabilitiesGraph(waitForIndexing)
	defer s.cleanImageToVulnerabilitiesGraph(waitForIndexing)
	s.Require().NoError(err)
	validClusters := s.dackboxTestStore.GetStoredClusterIDs()
	s.Require().True(len(validClusters) >= 2)

	testEdgeIDs := getClusterCVEEdges(validClusters[0], validClusters[1])
	testCases := getClusterCVEEdgeWriteTestCases(s.T(), validClusters[0], validClusters[1])
	embeddedCVE := fixtures.GetEmbeddedClusterCVE4567x0002()
	storageCVE := utils.EmbeddedCVEToProtoCVE(clusterOS, embeddedCVE)
	cveFixVersion := embeddedCVE.GetFixedBy()
	cluster1 := &storage.Cluster{
		Id:   validClusters[0],
		Name: validClusters[0],
		Type: storage.ClusterType_OPENSHIFT_CLUSTER,
	}
	cluster1Only := []*storage.Cluster{cluster1}
	cveParts := cveConverterV1.NewClusterCVEParts(storageCVE, cluster1Only, cveFixVersion)
	targetCveEdgeID := testEdgeIDs[cluster1ToCVE2EdgeID]
	err = s.datastore.Delete(allAccessCtx, targetCveEdgeID)
	s.Require().NoError(err)
	s.waitForIndexing()
	s.checkClusterCVEEdgePresence(targetCveEdgeID, false)
	for _, c := range testCases {
		s.Run(c.name, func() {
			ctx := c.ctx
			err = s.datastore.Upsert(allAccessCtx, cveParts)
			s.Require().NoError(err)
			s.waitForIndexing()
			s.checkClusterCVEEdgePresence(targetCveEdgeID, true)
			err = s.datastore.Delete(ctx, targetCveEdgeID)
			if c.visibleEdges[cluster1ToCVE2EdgeID] {
				s.NoError(err)
				s.waitForIndexing()
				s.checkClusterCVEEdgePresence(targetCveEdgeID, false)
			} else {
				s.ErrorIs(err, sac.ErrResourceAccessDenied)
				s.checkClusterCVEEdgePresence(targetCveEdgeID, true)
			}
			err = s.datastore.Delete(allAccessCtx, targetCveEdgeID)
			s.Require().NoError(err)
			s.waitForIndexing()
			s.checkClusterCVEEdgePresence(targetCveEdgeID, false)
		})
	}
}

func (s *clusterCVEEdgeDatastoreSACSuite) TestExists() {
	err := s.dackboxTestStore.PushClusterToVulnerabilitiesGraph(waitForIndexing)
	defer s.cleanImageToVulnerabilitiesGraph(waitForIndexing)
	s.Require().NoError(err)
	validClusters := s.dackboxTestStore.GetStoredClusterIDs()
	s.Require().True(len(validClusters) >= 2)

	testCases := getClusterCVEEdgeReadTestCases(s.T(), validClusters[0], validClusters[1])
	for _, c := range testCases {
		s.Run(c.name, func() {
			ctx := c.ctx
			targetEdgeID := getClusterCVEEdges(validClusters[0], validClusters[1])[cluster1ToCVE1EdgeID]
			exists, err := s.datastore.Exists(ctx, targetEdgeID)
			s.NoError(err)
			s.Equal(c.visibleEdges[cluster1ToCVE1EdgeID], exists)
		})
	}
}

func (s *clusterCVEEdgeDatastoreSACSuite) TestGet() {
	err := s.dackboxTestStore.PushClusterToVulnerabilitiesGraph(waitForIndexing)
	defer s.cleanImageToVulnerabilitiesGraph(waitForIndexing)
	s.Require().NoError(err)
	validClusters := s.dackboxTestStore.GetStoredClusterIDs()
	s.Require().True(len(validClusters) >= 2)

	targetEdgeCveID := getCveID(embeddedCVE1, clusterOS)
	targetEdgeID := getClusterCVEEdges(validClusters[0], validClusters[1])[cluster1ToCVE1EdgeID]
	testCases := getClusterCVEEdgeReadTestCases(s.T(), validClusters[0], validClusters[1])
	for _, c := range testCases {
		s.Run(c.name, func() {
			ctx := c.ctx
			edge, exists, err := s.datastore.Get(ctx, targetEdgeID)
			s.NoError(err)
			if c.visibleEdges[cluster1ToCVE1EdgeID] {
				s.True(exists)
				s.NotNil(edge)
				if !env.PostgresDatastoreEnabled.BooleanSetting() {
					// In dackbox mode, the edge ID is base64(clusterID):base64(cveID)
					idParts := strings.Split(edge.GetId(), ":")
					s.Require().Equal(2, len(idParts))
					edgeClusterID, err := base64.RawURLEncoding.DecodeString(idParts[0])
					s.NoError(err)
					edgeCveID, err := base64.RawURLEncoding.DecodeString(idParts[1])
					s.NoError(err)
					s.Equal(validClusters[0], string(edgeClusterID))
					s.Equal(targetEdgeCveID, string(edgeCveID))
				} else {
					s.Equal(validClusters[0], edge.GetClusterId())
					s.Equal(targetEdgeCveID, edge.GetCveId())
				}
			} else {
				s.False(exists)
				s.Nil(edge)
			}
		})
	}
}

func (s *clusterCVEEdgeDatastoreSACSuite) TestGetBatch() {
	err := s.dackboxTestStore.PushClusterToVulnerabilitiesGraph(waitForIndexing)
	defer s.cleanImageToVulnerabilitiesGraph(waitForIndexing)
	s.Require().NoError(err)
	validClusters := s.dackboxTestStore.GetStoredClusterIDs()
	s.Require().True(len(validClusters) >= 2)

	testEdgeIDs := getClusterCVEEdges(validClusters[0], validClusters[1])
	testCases := getClusterCVEEdgeReadTestCases(s.T(), validClusters[0], validClusters[1])
	targetEdgeIDs := []string{
		testEdgeIDs[cluster1ToCVE1EdgeID],
		testEdgeIDs[cluster1ToCVE2EdgeID],
		testEdgeIDs[cluster2ToCVE2EdgeID],
		testEdgeIDs[cluster2ToCVE3EdgeID],
	}
	for _, c := range testCases {
		s.Run(c.name, func() {
			ctx := c.ctx
			edges, err := s.datastore.GetBatch(ctx, targetEdgeIDs)
			s.NoError(err)
			edgesByID := make(map[string]*storage.ClusterCVEEdge, 0)
			for ix := range edges {
				e := edges[ix]
				edgesByID[e.GetId()] = e
			}
			visibleEdgeCount := 0
			for id, visible := range c.visibleEdges {
				if visible {
					_, found := edgesByID[testEdgeIDs[id]]
					s.True(found)
					visibleEdgeCount++
				}
			}
			s.Equal(visibleEdgeCount, len(edges))
		})
	}
}

func (s *clusterCVEEdgeDatastoreSACSuite) TestCount() {
	if !env.PostgresDatastoreEnabled.BooleanSetting() {
		s.T().Skip("Count panics in non-postgres mode")
	}
	err := s.dackboxTestStore.PushClusterToVulnerabilitiesGraph(waitForIndexing)
	defer s.cleanImageToVulnerabilitiesGraph(waitForIndexing)
	s.Require().NoError(err)
	validClusters := s.dackboxTestStore.GetStoredClusterIDs()
	s.Require().True(len(validClusters) >= 2)

	testCases := getClusterCVEEdgeReadTestCases(s.T(), validClusters[0], validClusters[1])
	for _, c := range testCases {
		s.Run(c.name, func() {
			ctx := c.ctx
			expectedCount := 0
			for _, visible := range c.visibleEdges {
				if visible {
					expectedCount++
				}
			}
			count, err := s.datastore.Count(ctx, search.EmptyQuery())
			s.NoError(err)
			s.Equal(expectedCount, count)
		})
	}
}

func (s *clusterCVEEdgeDatastoreSACSuite) TestSearch() {
	err := s.dackboxTestStore.PushClusterToVulnerabilitiesGraph(waitForIndexing)
	defer s.cleanImageToVulnerabilitiesGraph(waitForIndexing)
	s.Require().NoError(err)
	validClusters := s.dackboxTestStore.GetStoredClusterIDs()
	s.Require().True(len(validClusters) >= 2)

	testEdgeIDs := getClusterCVEEdges(validClusters[0], validClusters[1])
	testCases := getClusterCVEEdgeReadTestCases(s.T(), validClusters[0], validClusters[1])
	for _, c := range testCases {
		s.Run(c.name, func() {
			ctx := c.ctx
			expectedIDs := make([]string, 0, len(c.visibleEdges))
			for id, visible := range c.visibleEdges {
				if visible {
					expectedIDs = append(expectedIDs, testEdgeIDs[id])
				}
			}
			results, err := s.datastore.Search(ctx, search.EmptyQuery())
			s.NoError(err)
			resultIDs := make([]string, 0, len(results))
			for _, r := range results {
				resultIDs = append(resultIDs, r.ID)
			}
			s.ElementsMatch(resultIDs, expectedIDs)
		})
	}
}

func (s *clusterCVEEdgeDatastoreSACSuite) TestSearchEdges() {
	if !env.PostgresDatastoreEnabled.BooleanSetting() {
		s.T().Skip("SearchEdges panics in non-postgres mode")
	}
	err := s.dackboxTestStore.PushClusterToVulnerabilitiesGraph(waitForIndexing)
	defer s.cleanImageToVulnerabilitiesGraph(waitForIndexing)
	s.Require().NoError(err)
	validClusters := s.dackboxTestStore.GetStoredClusterIDs()
	s.Require().True(len(validClusters) >= 2)

	testEdgeIDs := getClusterCVEEdges(validClusters[0], validClusters[1])
	testCases := getClusterCVEEdgeReadTestCases(s.T(), validClusters[0], validClusters[1])
	for _, c := range testCases {
		s.Run(c.name, func() {
			ctx := c.ctx
			expectedIDs := make([]string, 0, len(c.visibleEdges))
			for id, visible := range c.visibleEdges {
				if visible {
					expectedIDs = append(expectedIDs, testEdgeIDs[id])
				}
			}
			results, err := s.datastore.SearchEdges(ctx, search.EmptyQuery())
			s.NoError(err)
			resultIDs := make([]string, 0, len(results))
			for _, r := range results {
				resultIDs = append(resultIDs, r.GetId())
			}
			s.ElementsMatch(resultIDs, expectedIDs)
		})
	}
}

func (s *clusterCVEEdgeDatastoreSACSuite) TestSearchRawEdges() {
	err := s.dackboxTestStore.PushClusterToVulnerabilitiesGraph(waitForIndexing)
	defer s.cleanImageToVulnerabilitiesGraph(waitForIndexing)
	s.Require().NoError(err)
	validClusters := s.dackboxTestStore.GetStoredClusterIDs()
	s.Require().True(len(validClusters) >= 2)

	testEdgeIDs := getClusterCVEEdges(validClusters[0], validClusters[1])
	testCases := getClusterCVEEdgeReadTestCases(s.T(), validClusters[0], validClusters[1])
	for _, c := range testCases {
		s.Run(c.name, func() {
			ctx := c.ctx
			expectedIDs := make([]string, 0, len(c.visibleEdges))
			for id, visible := range c.visibleEdges {
				if visible {
					expectedIDs = append(expectedIDs, testEdgeIDs[id])
				}
			}
			results, err := s.datastore.SearchRawEdges(ctx, search.EmptyQuery())
			s.NoError(err)
			resultIDs := make([]string, 0, len(results))
			for _, r := range results {
				resultIDs = append(resultIDs, r.GetId())
			}
			s.ElementsMatch(resultIDs, expectedIDs)
		})
	}
}
