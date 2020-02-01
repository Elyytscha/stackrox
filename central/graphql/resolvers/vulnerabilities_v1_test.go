package resolvers

import (
	"context"
	"testing"

	"github.com/facebookincubator/nvdtools/cvefeed/nvd/schema"
	"github.com/golang/mock/gomock"
	clusterMocks "github.com/stackrox/rox/central/cluster/datastore/mocks"
	"github.com/stackrox/rox/central/cve/converter"
	imageMocks "github.com/stackrox/rox/central/image/datastore/mocks"
	namespaceMocks "github.com/stackrox/rox/central/namespace/datastore/mocks"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stretchr/testify/assert"
)

func TestMapImagesToVulnerabilityResolvers(t *testing.T) {
	fakeRoot := &Resolver{}
	images := testImages()

	query := &v1.Query{}
	vulnerabilityResolvers, err := mapImagesToVulnerabilityResolvers(fakeRoot, images, query)
	assert.NoError(t, err)
	assert.Len(t, vulnerabilityResolvers, 5)

	query = search.NewQueryBuilder().AddExactMatches(search.FixedBy, "1.1").ProtoQuery()
	vulnerabilityResolvers, err = mapImagesToVulnerabilityResolvers(fakeRoot, images, query)
	assert.NoError(t, err)
	assert.Len(t, vulnerabilityResolvers, 1)

	query = search.NewQueryBuilder().AddExactMatches(search.CVE, "cve-2019-1", "cve-2019-2", "cve-2019-3").ProtoQuery()
	vulnerabilityResolvers, err = mapImagesToVulnerabilityResolvers(fakeRoot, images, query)
	assert.NoError(t, err)
	assert.Len(t, vulnerabilityResolvers, 2)
}

func TestIfSpecificVersionCVEAffectsCluster(t *testing.T) {
	cluster := &storage.Cluster{
		Id:   "test_cluster_id1",
		Name: "cluster1",
		Status: &storage.ClusterStatus{
			OrchestratorMetadata: &storage.OrchestratorMetadata{
				Version: "v1.15.3",
			},
		},
	}

	cve1 := &schema.NVDCVEFeedJSON10DefCVEItem{
		CVE: &schema.CVEJSON40{
			CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
				ID: "CVE-1",
			},
		},
		Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
			Nodes: []*schema.NVDCVEFeedJSON10DefNode{
				{
					Operator: "OR",
					CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
						{
							Vulnerable: true,
							Cpe23Uri:   "cpe:2.3:a:kubernetes:kubernetes:1.15.3:*:*:*:*:*:*:*",
						},
					},
				},
			},
		},
	}

	cve2 := &schema.NVDCVEFeedJSON10DefCVEItem{
		CVE: &schema.CVEJSON40{
			CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
				ID: "CVE-2",
			},
		},
		Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
			Nodes: []*schema.NVDCVEFeedJSON10DefNode{
				{
					Operator: "OR",
					CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
						{
							Vulnerable: true,
							Cpe23Uri:   "cpe:2.3:a:kubernetes:kubernetes:1.14.3:*:*:*:*:*:*:*",
						},
					},
				},
			},
		},
	}

	cve3 := &schema.NVDCVEFeedJSON10DefCVEItem{
		CVE: &schema.CVEJSON40{
			CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
				ID: "CVE-3",
			},
		},
		Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
			Nodes: []*schema.NVDCVEFeedJSON10DefNode{
				{
					Operator: "OR",
					CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
						{
							Vulnerable: true,
							Cpe23Uri:   "cpe:2.3:a:kubernetes:kubernetes:1.15.3:alpha1:*:*:*:*:*:*",
						},
					},
				},
			},
		},
	}

	ret := isClusterAffectedByK8sCVE(cluster, cve1)
	assert.Equal(t, true, ret)
	ret = isClusterAffectedByK8sCVE(cluster, cve2)
	assert.Equal(t, false, ret)
	ret = isClusterAffectedByK8sCVE(cluster, cve3)
	assert.Equal(t, false, ret)

	cluster.Status.OrchestratorMetadata.Version = "v1.15.3+build1"
	ret = isClusterAffectedByK8sCVE(cluster, cve1)
	assert.Equal(t, true, ret)
	ret = isClusterAffectedByK8sCVE(cluster, cve2)
	assert.Equal(t, false, ret)
	ret = isClusterAffectedByK8sCVE(cluster, cve3)
	assert.Equal(t, false, ret)

	cluster.Status.OrchestratorMetadata.Version = "v1.15.3-alpha1"
	ret = isClusterAffectedByK8sCVE(cluster, cve1)
	assert.Equal(t, true, ret)
	ret = isClusterAffectedByK8sCVE(cluster, cve2)
	assert.Equal(t, false, ret)
	ret = isClusterAffectedByK8sCVE(cluster, cve3)
	assert.Equal(t, true, ret)

	cluster.Status.OrchestratorMetadata.Version = "v1.15.3-alpha1+build1"
	ret = isClusterAffectedByK8sCVE(cluster, cve1)
	assert.Equal(t, true, ret)
	ret = isClusterAffectedByK8sCVE(cluster, cve2)
	assert.Equal(t, false, ret)
	ret = isClusterAffectedByK8sCVE(cluster, cve3)
	assert.Equal(t, true, ret)
}

func TestMultipleVersionCVEAffectsCluster(t *testing.T) {
	cluster := &storage.Cluster{
		Id:   "test_cluster_id1",
		Name: "cluster1",
		Status: &storage.ClusterStatus{
			OrchestratorMetadata: &storage.OrchestratorMetadata{
				Version: "v1.15.3",
			},
		},
	}

	cve1 := &schema.NVDCVEFeedJSON10DefCVEItem{
		CVE: &schema.CVEJSON40{
			CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
				ID: "CVE-1",
			},
		},
		Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
			Nodes: []*schema.NVDCVEFeedJSON10DefNode{
				{
					Operator: "OR",
					CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
						{
							Vulnerable:            true,
							Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
							VersionStartIncluding: "1.10.1",
							VersionEndExcluding:   "1.10.9",
						},
						{
							Vulnerable:            true,
							Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
							VersionStartIncluding: "1.11.1",
							VersionEndExcluding:   "1.11.9",
						},
					},
				},
			},
		},
	}

	cve2 := &schema.NVDCVEFeedJSON10DefCVEItem{
		CVE: &schema.CVEJSON40{
			CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
				ID: "CVE-2",
			},
		},
		Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
			Nodes: []*schema.NVDCVEFeedJSON10DefNode{
				{
					Operator: "OR",
					CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
						{
							Vulnerable:            true,
							Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
							VersionStartIncluding: "1.14.1",
							VersionEndExcluding:   "1.14.9",
						},
						{
							Vulnerable:            true,
							Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
							VersionStartIncluding: "1.15.1",
							VersionEndExcluding:   "1.15.9",
						},
					},
				},
			},
		},
	}

	cve3 := &schema.NVDCVEFeedJSON10DefCVEItem{
		CVE: &schema.CVEJSON40{
			CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
				ID: "CVE-3",
			},
		},
		Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
			Nodes: []*schema.NVDCVEFeedJSON10DefNode{
				{
					Operator: "OR",
					CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
						{
							Vulnerable:            true,
							Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
							VersionStartIncluding: "1.16.5",
							VersionEndIncluding:   "1.16.9",
						},
						{
							Vulnerable:            true,
							Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
							VersionStartIncluding: "1.17.1",
							VersionEndIncluding:   "1.17.9",
						},
					},
				},
			},
		},
	}

	ret := isClusterAffectedByK8sCVE(cluster, cve1)
	assert.Equal(t, false, ret)

	ret = isClusterAffectedByK8sCVE(cluster, cve2)
	assert.Equal(t, true, ret)

	cluster.Status.OrchestratorMetadata.Version = "v1.15.1-beta1"
	ret = isClusterAffectedByK8sCVE(cluster, cve1)
	assert.Equal(t, false, ret)

	ret = isClusterAffectedByK8sCVE(cluster, cve2)
	assert.Equal(t, true, ret)

	cluster.Status.OrchestratorMetadata.Version = "v1.15.9"
	ret = isClusterAffectedByK8sCVE(cluster, cve1)
	assert.Equal(t, false, ret)

	ret = isClusterAffectedByK8sCVE(cluster, cve2)
	assert.Equal(t, false, ret)

	cluster.Status.OrchestratorMetadata.Version = "v1.15.1"
	ret = isClusterAffectedByK8sCVE(cluster, cve1)
	assert.Equal(t, false, ret)

	ret = isClusterAffectedByK8sCVE(cluster, cve2)
	assert.Equal(t, true, ret)

	cluster.Status.OrchestratorMetadata.Version = "v1.16.4"
	ret = isClusterAffectedByK8sCVE(cluster, cve3)
	assert.Equal(t, false, ret)

	cluster.Status.OrchestratorMetadata.Version = "v1.17.4"
	ret = isClusterAffectedByK8sCVE(cluster, cve3)
	assert.Equal(t, true, ret)
}

func TestSingleAndMultipleVersionCVEAffectsCluster(t *testing.T) {
	cluster := &storage.Cluster{
		Id:   "test_cluster_id1",
		Name: "cluster1",
		Status: &storage.ClusterStatus{
			OrchestratorMetadata: &storage.OrchestratorMetadata{
				Version: "v1.10.6",
			},
		},
	}
	cve := &schema.NVDCVEFeedJSON10DefCVEItem{
		CVE: &schema.CVEJSON40{
			CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
				ID: "CVE-1",
			},
		},
		Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
			Nodes: []*schema.NVDCVEFeedJSON10DefNode{
				{
					Operator: "OR",
					CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
						{
							Vulnerable:            true,
							Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
							VersionStartIncluding: "1.10.1",
							VersionEndExcluding:   "1.10.9",
						},
						{
							Vulnerable:            true,
							Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
							VersionStartIncluding: "1.11.1",
							VersionEndExcluding:   "1.11.9",
						},
						{
							Vulnerable: true,
							Cpe23Uri:   "cpe:2.3:a:kubernetes:kubernetes:1.10.3:alpha1:*:*:*:*:*:*",
						},
						{
							Vulnerable: true,
							Cpe23Uri:   "cpe:2.3:a:kubernetes:kubernetes:1.10.3:beta1:*:*:*:*:*:*",
						},
					},
				},
			},
		},
	}

	ret := isClusterAffectedByK8sCVE(cluster, cve)
	assert.Equal(t, true, ret)

	cluster.Status.OrchestratorMetadata.Version = "v1.10.3-alpha1"
	ret = isClusterAffectedByK8sCVE(cluster, cve)
	assert.Equal(t, true, ret)

	cluster.Status.OrchestratorMetadata.Version = "v1.10.3-beta1"
	ret = isClusterAffectedByK8sCVE(cluster, cve)
	assert.Equal(t, true, ret)

	cluster.Status.OrchestratorMetadata.Version = "v1.10.3-rc1"
	ret = isClusterAffectedByK8sCVE(cluster, cve)
	assert.Equal(t, true, ret)
}

func TestCountCVEsAffectsCluster(t *testing.T) {
	cluster := &storage.Cluster{
		Id:   "test_cluster_id1",
		Name: "cluster1",
		Status: &storage.ClusterStatus{
			OrchestratorMetadata: &storage.OrchestratorMetadata{
				Version: "v1.10.6",
			},
		},
	}
	cves := []*schema.NVDCVEFeedJSON10DefCVEItem{
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-1",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.10.1",
								VersionEndExcluding:   "1.10.9",
							},
						},
					},
				},
			},
		},
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-2",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable: true,
								Cpe23Uri:   "cpe:2.3:a:kubernetes:kubernetes:1.10.6:*:*:*:*:*:*:*",
							},
						},
					},
				},
			},
		},
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-3",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.10.3",
								VersionEndIncluding:   "1.10.7",
							},
						},
					},
				},
			},
		},
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-4",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.11.3",
								VersionEndIncluding:   "1.11.7",
							},
						},
					},
				},
			},
		},
	}

	var countAffectedClusters, countFixableCVEs int
	for _, cve := range cves {
		if isClusterAffectedByK8sCVE(cluster, cve) {
			countAffectedClusters++
		}
		if isK8sCVEFixable(cve) {
			countFixableCVEs++
		}
	}
	assert.Equal(t, countAffectedClusters, 3)
	assert.Equal(t, countFixableCVEs, 1)
}

func TestNonK8sCPEMatch(t *testing.T) {
	cluster := &storage.Cluster{
		Id:   "test_cluster_id1",
		Name: "cluster1",
		Status: &storage.ClusterStatus{
			OrchestratorMetadata: &storage.OrchestratorMetadata{
				Version: "v1.10.6",
			},
		},
	}

	cves := []*schema.NVDCVEFeedJSON10DefCVEItem{
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-2019-1",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:vendorfoo:projectbar:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.10.1",
								VersionEndExcluding:   "1.10.9",
							},
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:vendorfoo:projectbar:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.11.1",
								VersionEndExcluding:   "1.11.9",
							},
							{
								Vulnerable: true,
								Cpe23Uri:   "cpe:2.3:a:vendorfoo:projectbar:1.13.1:*:*:*:*:*:*:*",
							},
						},
					},
				},
			},
		},
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-2019-2",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:vendorfoo:projectbar:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.10.1",
								VersionEndExcluding:   "1.10.9",
							},
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:vendorfoo:projectbar:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.11.1",
								VersionEndExcluding:   "1.11.9",
							},
							{
								Vulnerable: true,
								Cpe23Uri:   "cpe:2.3:a:vendorfoo:projectbar:1.13.1:*:*:*:*:*:*:*",
							},
						},
					},
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.10.1",
								VersionEndExcluding:   "1.10.9",
							},
							{
								Vulnerable: true,
								Cpe23Uri:   "cpe:2.3:a:kubernetes:kubernetes:1.13.1:*:*:*:*:*:*:*",
							},
						},
					},
				},
			},
		},
	}

	ret := isClusterAffectedByK8sCVE(cluster, cves[0])
	assert.Equal(t, false, ret)
	ret = isClusterAffectedByK8sCVE(cluster, cves[1])
	assert.Equal(t, true, ret)
}

func TestFixableCVEs(t *testing.T) {
	cves := []*schema.NVDCVEFeedJSON10DefCVEItem{
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-2019-1",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.10.1",
								VersionEndExcluding:   "1.10.9",
							},
						},
					},
				},
			},
		},
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-2",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable: true,
								Cpe23Uri:   "cpe:2.3:a:kubernetes:kubernetes:1.10.6:*:*:*:*:*:*:*",
							},
						},
					},
				},
			},
		},
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-3",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.10.3",
								VersionEndIncluding:   "1.10.7",
							},
						},
					},
				},
			},
		},
	}
	actual := isK8sCVEFixable(cves[0])
	assert.Equal(t, actual, true)
	actual = isK8sCVEFixable(cves[1])
	assert.Equal(t, actual, false)
	actual = isK8sCVEFixable(cves[2])
	assert.Equal(t, actual, false)
}

func TestK8sCVEEnvImpact(t *testing.T) {
	expected := []float64{0.6, 0.4, 0.4}

	clusters := []*storage.Cluster{
		{
			Id:   "test_cluster_id1",
			Name: "cluster1",
			Status: &storage.ClusterStatus{
				OrchestratorMetadata: &storage.OrchestratorMetadata{
					Version: "v1.14.2",
				},
			},
		},
		{
			Id:   "test_cluster_id2",
			Name: "cluster2",
			Status: &storage.ClusterStatus{
				OrchestratorMetadata: &storage.OrchestratorMetadata{
					Version: "v1.14.5+build1",
				},
			},
		},
		{
			Id:   "test_cluster_id3",
			Name: "cluster3",
			Status: &storage.ClusterStatus{
				OrchestratorMetadata: &storage.OrchestratorMetadata{
					Version: "v1.15.4-beta1",
				},
			},
		},
		{
			Id:   "test_cluster_id4",
			Name: "cluster4",
			Status: &storage.ClusterStatus{
				OrchestratorMetadata: &storage.OrchestratorMetadata{
					Version: "v1.16.3-alpha1+build2",
				},
			},
		},
		{
			Id:   "test_cluster_id5",
			Name: "cluster4",
			Status: &storage.ClusterStatus{
				OrchestratorMetadata: &storage.OrchestratorMetadata{
					Version: "v1.17.5",
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	clusterDataStore := clusterMocks.NewMockDataStore(ctrl)
	clusterDataStore.EXPECT().GetClusters(gomock.Any()).Return(clusters, nil).AnyTimes()

	resolver := &Resolver{
		ClusterDataStore: clusterDataStore,
	}

	cves := []*schema.NVDCVEFeedJSON10DefCVEItem{
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-2019-1",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable: true,
								Cpe23Uri:   "cpe:2.3:a:kubernetes:kubernetes:1.15.4:*:*:*:*:*:*:*",
							},
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.16.1",
								VersionEndIncluding:   "1.16.9",
							},
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.17.1",
								VersionEndExcluding:   "1.17.7",
							},
						},
					},
				},
			},
		},
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-2019-2",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.14.1",
								VersionEndExcluding:   "1.14.9",
							},
						},
					},
				},
			},
		},
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-2019-3",
				},
			},

			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.10.1",
								VersionEndIncluding:   "1.10.9",
							},
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:kubernetes:kubernetes:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.11.1",
								VersionEndExcluding:   "1.11.7",
							},
							{
								Vulnerable: true,
								Cpe23Uri:   "cpe:2.3:a:kubernetes:kubernetes:1.14.5:*:*:*:*:*:*:*",
							},
							{
								Vulnerable: true,
								Cpe23Uri:   "cpe:2.3:a:kubernetes:kubernetes:1.15.4:alpha1:*:*:*:*:*:*",
							},
							{
								Vulnerable: true,
								Cpe23Uri:   "cpe:2.3:a:kubernetes:kubernetes:1.15.4:beta1:*:*:*:*:*:*",
							},
						},
					},
				},
			},
		},
	}

	for i, cve := range cves {
		actual, err := resolver.getAffectedClusterPercentage(context.Background(), cve, converter.K8s)
		assert.Nil(t, err)
		assert.Equal(t, actual, expected[i])
	}
}

func TestIstioCVEImpactsCluster(t *testing.T) {

	expected := []bool{true, true, true, false}

	clusters := []*storage.Cluster{
		{
			Id:   "test_cluster_id1",
			Name: "cluster1",
			Status: &storage.ClusterStatus{
				OrchestratorMetadata: &storage.OrchestratorMetadata{
					Version: "v1.14.2",
				},
			},
		},
	}

	namespaces := []*storage.NamespaceMetadata{
		{
			Id:          "test_namespace1",
			Name:        "istio-system",
			ClusterId:   "test_cluster_id1",
			ClusterName: "cluster1",
		},
	}

	images := []*storage.Image{
		{
			Id: "test_image_id1",
			Name: &storage.ImageName{
				Tag:      "1.2.6",
				Remote:   "istio/proxyv2",
				Registry: "docker.io",
				FullName: "docker.io/istio/proxyv2:1.2.6",
			},
		},
		{
			Id: "test_image_id2",
			Name: &storage.ImageName{
				Tag:      "1.4.4",
				Remote:   "istio/node-agent-k8s",
				Registry: "docker.io",
				FullName: "docker.io/istio/node-agent-k8s:1.4.4",
			},
		},
		{
			Id: "test_image_id3",
			Name: &storage.ImageName{
				Tag:      "v1.13.11-gke.14",
				Remote:   "kube-proxy",
				Registry: "gke.gcr.io",
				FullName: "gke.gcr.io/kube-proxy:v1.13.11-gke.14",
			},
		},
		{
			Id: "test_image_id4",
			Name: &storage.ImageName{
				Tag:      "v0.11.0",
				Remote:   "jetstack/cert-manager-controller",
				Registry: "quay.io",
				FullName: "quay.io/jetstack/cert-manager-controller:v0.11.0",
			},
		},
	}

	ctrl := gomock.NewController(t)
	clusterDataStore := clusterMocks.NewMockDataStore(ctrl)
	clusterDataStore.EXPECT().GetClusters(gomock.Any()).Return(clusters, nil).AnyTimes()

	imageDataStore := imageMocks.NewMockDataStore(ctrl)
	imageDataStore.EXPECT().SearchRawImages(gomock.Any(), gomock.Any()).Return(images, nil).AnyTimes()

	namespaceDataStore := namespaceMocks.NewMockDataStore(ctrl)
	namespaceDataStore.EXPECT().SearchNamespaces(gomock.Any(), gomock.Any()).Return(namespaces, nil).AnyTimes()

	resolver := &Resolver{
		ClusterDataStore:   clusterDataStore,
		ImageDataStore:     imageDataStore,
		NamespaceDataStore: namespaceDataStore,
	}

	ok, err := resolver.isIstioControlPlaneRunning(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, ok, true)

	istioCVEs := []*schema.NVDCVEFeedJSON10DefCVEItem{
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-2019-1",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:istio:istio:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "",
								VersionEndIncluding:   "",
								VersionEndExcluding:   "1.1.13",
							},
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:istio:istio:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.2.0",
								VersionEndIncluding:   "",
								VersionEndExcluding:   "1.2.8",
							},
						},
					},
				},
			},
		},
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-2019-2",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:istio:istio:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.4.1",
								VersionEndIncluding:   "",
								VersionEndExcluding:   "1.4.9",
							},
						},
					},
				},
			},
		},
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-2019-3",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:istio:istio:1.4.4:*:*:*:*:*:*:*",
								VersionStartIncluding: "",
								VersionEndIncluding:   "",
								VersionEndExcluding:   "",
							},
						},
					},
				},
			},
		},
		{
			CVE: &schema.CVEJSON40{
				CVEDataMeta: &schema.CVEJSON40CVEDataMeta{
					ID: "CVE-2019-4",
				},
			},
			Configurations: &schema.NVDCVEFeedJSON10DefConfigurations{
				Nodes: []*schema.NVDCVEFeedJSON10DefNode{
					{
						Operator: "OR",
						CPEMatch: []*schema.NVDCVEFeedJSON10DefCPEMatch{
							{
								Vulnerable:            true,
								Cpe23Uri:              "cpe:2.3:a:istio:istio:*:*:*:*:*:*:*:*",
								VersionStartIncluding: "1.3.1",
								VersionEndIncluding:   "",
								VersionEndExcluding:   "1.3.9",
							},
						},
					},
				},
			},
		},
	}

	for i, cve := range istioCVEs {
		actual, err := resolver.isClusterAffectedByIstioCVE(context.Background(), clusters[0], cve)
		assert.Nil(t, err)
		assert.Equal(t, expected[i], actual)
	}
}
