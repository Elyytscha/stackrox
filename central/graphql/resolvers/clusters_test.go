package resolvers

import (
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stretchr/testify/assert"
)

func TestPolicyAppliesToCluster(t *testing.T) {
	cases := []struct {
		policy     *storage.Policy
		cluster    *storage.Cluster
		namespaces []*storage.NamespaceMetadata
		matches    bool
	}{
		{
			cluster: &storage.Cluster{
				Id: "cluster1",
			},
			namespaces: []*storage.NamespaceMetadata{
				{
					ClusterId: "cluster1",
					Name:      "ns1",
				},
			},
			policy: &storage.Policy{
				Scope: []*storage.Scope{
					{
						Cluster: "cluster1",
					},
				},
			},
			matches: true,
		},
		{
			cluster: &storage.Cluster{
				Id: "cluster1",
			},
			namespaces: []*storage.NamespaceMetadata{
				{
					ClusterId: "cluster1",
					Name:      "ns1",
				},
			},
			policy: &storage.Policy{
				Scope: []*storage.Scope{
					{
						Cluster:   "cluster2",
						Namespace: "ns1",
					},
				},
			},
			matches: false,
		},
		{
			cluster: &storage.Cluster{
				Id: "cluster1",
			},
			namespaces: []*storage.NamespaceMetadata{
				{
					ClusterId: "cluster1",
					Name:      "ns1",
				},
			},
			policy: &storage.Policy{
				Scope: []*storage.Scope{
					{
						Namespace: "ns1",
					},
				},
			},
			matches: true,
		},
		{
			cluster: &storage.Cluster{
				Id: "cluster1",
			},
			namespaces: []*storage.NamespaceMetadata{
				{
					ClusterId: "cluster1",
					Name:      "ns1",
				},
			},
			policy:  &storage.Policy{},
			matches: true,
		},
	}

	for _, c := range cases {
		actual := policyAppliesToCluster(c.cluster, c.namespaces, c.policy)
		assert.Equal(t, c.matches, actual)
	}
}

func TestPolicyAppliesToClusterWithWhitelist(t *testing.T) {
	cases := []struct {
		policy     *storage.Policy
		cluster    *storage.Cluster
		namespaces []*storage.NamespaceMetadata
		matches    bool
	}{
		{
			cluster: &storage.Cluster{
				Id: "cluster1",
			},
			namespaces: []*storage.NamespaceMetadata{
				{
					ClusterId: "cluster1",
					Name:      "ns1",
				},
			},
			policy: &storage.Policy{
				Scope: []*storage.Scope{
					{
						Cluster: "cluster1",
					},
				},
				Whitelists: []*storage.Whitelist{
					{
						Deployment: &storage.Whitelist_Deployment{
							Scope: &storage.Scope{
								Namespace: "ns.*",
							},
						},
					},
				},
			},
			matches: false,
		},
		{
			cluster: &storage.Cluster{
				Id: "cluster1",
			},
			namespaces: []*storage.NamespaceMetadata{
				{
					ClusterId: "cluster1",
					Name:      "ns1",
				},
			},
			policy: &storage.Policy{
				Scope: []*storage.Scope{
					{
						Cluster:   "cluster2",
						Namespace: "ns1",
					},
				},
				Whitelists: []*storage.Whitelist{
					{
						Deployment: &storage.Whitelist_Deployment{
							Scope: &storage.Scope{
								Namespace: "ns.*",
							},
						},
					},
				},
			},
			matches: false,
		},
		{
			cluster: &storage.Cluster{
				Id: "cluster1",
			},
			namespaces: []*storage.NamespaceMetadata{
				{
					ClusterId: "cluster1",
					Name:      "ns1",
				},
			},
			policy: &storage.Policy{
				Scope: []*storage.Scope{
					{
						Namespace: "ns1",
					},
				},
				Whitelists: []*storage.Whitelist{
					{
						Deployment: &storage.Whitelist_Deployment{
							Scope: &storage.Scope{
								Cluster: "cluster1",
							},
						},
					},
				},
			},
			matches: false,
		},
		{
			cluster: &storage.Cluster{
				Id: "cluster1",
			},
			namespaces: []*storage.NamespaceMetadata{
				{
					ClusterId: "cluster1",
					Name:      "ns1",
				},
			},
			policy: &storage.Policy{
				Whitelists: []*storage.Whitelist{
					{
						Deployment: &storage.Whitelist_Deployment{
							Scope: &storage.Scope{
								Namespace: "ns2.*",
							},
						},
					},
				},
			},
			matches: true,
		},
	}

	for _, c := range cases {
		actual := policyAppliesToCluster(c.cluster, c.namespaces, c.policy)
		assert.Equal(t, c.matches, actual)
	}
}
