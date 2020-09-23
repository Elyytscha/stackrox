package graph

import (
	"bytes"
	"context"
	"sort"
	"testing"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/labels"
	"github.com/stackrox/rox/pkg/networkgraph"
	networkPolicyConversion "github.com/stackrox/rox/pkg/protoconv/networkpolicy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8sV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var networkPolicyFixtures = map[string]*storage.NetworkPolicy{}

func init() {
	for _, policyYAML := range networkPolicyFixtureYAMLs {
		var k8sNp k8sV1.NetworkPolicy
		if err := yaml.NewYAMLToJSONDecoder(bytes.NewReader([]byte(policyYAML))).Decode(&k8sNp); err != nil {
			panic(err)
		}
		np := networkPolicyConversion.KubernetesNetworkPolicyWrap{NetworkPolicy: &k8sNp}.ToRoxNetworkPolicy()
		np.Id = k8sNp.GetName()
		if np.GetNamespace() == "" {
			np.Namespace = "default"
		}
		networkPolicyFixtures[np.GetName()] = np
	}
}

var networkPolicyFixtureYAMLs = []string{
	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: allow-only-egress-to-ipblock
  namespace: default
spec:
  policyTypes:
  - Egress
  - Ingress
  podSelector: {}
  egress:
  - to:
    - ipblock:
        cidr: 172.17.0.0/16
`,
	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: allow-traffic-from-apps-using-multiple-selectors
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: web
      role: db
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: bookstore
          role: search
    - podSelector:
            matchLabels:
              app: bookstore
              role: api

`,

	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: default-deny-all
  namespace: default
spec:
  podSelector: {}
  ingress: []
`,

	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: default-deny-all-egress
  namespace: default
spec:
  policyTypes:
  - Egress
  podSelector: {}
  egress: []
`,

	`
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: web-deny-external-egress
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: web
  policyTypes:
  - Egress
  egress:
  - ports:
    - port: 53
      protocol: UDP
    - port: 53
      protocol: TCP
    to:
    - namespaceSelector: {}
`,

	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  namespace: default
  name: deny-from-other-namespaces
spec:
  podSelector:
    matchLabels:
  ingress:
  - from:
    - podSelector: {}
`,

	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: limit-traffic
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: bookstore
      role: api
  ingress:
  - from:
      - podSelector:
          matchLabels:
            app: bookstore
`,

	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  namespace: default
  name: web-allow-all-namespaces
spec:
  podSelector:
    matchLabels:
      app: web
  ingress:
  - from:
    - namespaceSelector: {}
`,

	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: web-allow-all
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: web
  ingress:
  - {}
`,

	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: web-allow-stackrox
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: web
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: stackrox
`,

	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: web-deny-all
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: web
  ingress: []
`,

	`
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: web-deny-egress
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: web
  policyTypes:
  - Egress
  egress: []
`,
	`
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all-ingress
  namespace: qa
spec:
  podSelector: {}
  policyTypes:
  - Ingress
`,
	`
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-ingress-to-web
  namespace: qa
spec:
  ingress:
  - from:
    - namespaceSelector: {}
  podSelector:
    matchLabels:
      app: web

`,

	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: api-allow-5000
spec:
  podSelector:
    matchLabels:
      app: apiserver
  ingress:
  - ports:
    - port: 5000
    from:
    - podSelector:
        matchLabels:
          role: monitoring
`,

	// Custom network policies to test port-aware behavior
	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: allow-dns-egress-only
spec:
  podSelector:
    matchLabels:
      app: apiserver
  egress:
  - ports:
    - port: 53
      protocol: TCP
    - port: 53
      protocol: UDP
  policyTypes:
  - Egress
`,
	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: api-allow-named-api-port
spec:
  podSelector:
    matchLabels:
      app: apiserver
  ingress:
  - ports:
    - port: api
    from:
    - podSelector:
        matchLabels:
          role: monitoring
`,

	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: api-allow-all-udp-from-monitoring
spec:
  podSelector:
    matchLabels:
      app: apiserver
  ingress:
  - ports:
    - protocol: UDP
    from:
    - podSelector:
        matchLabels:
          role: monitoring
`,

	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: fully-isolate
spec:
  podSelector: {}
  ingress: []
  egress: []
  podSelector:
    matchExpressions: []
    matchLabels: {}
  policyTypes:
  - Ingress
  - Egress
`,

	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: a-ingress-tcp-8080
spec:
  podSelector:
    matchLabels:
      app: a
  ingress:
  - ports:
    - port: 8080
      protocol: TCP
  policyTypes:
  - Ingress
`,

	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: b-egress-a-tcp-ports-and-dns
spec:
  podSelector:
    matchLabels:
      app: b
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: a
    ports:
    - protocol: TCP
    - port: 53
      protocol: UDP
  policyTypes:
  - Egress
`,
	`
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: c-egress-a-tcp-8443-and-udp
spec:
  podSelector:
    matchLabels:
      app: c
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: a
    ports:
    - protocol: TCP
      port: 8443
    - protocol: UDP
  policyTypes:
  - Egress
`,
}

var (
	namespaces = []*storage.NamespaceMetadata{
		{
			Name: "default",
			Id:   "default",
			Labels: map[string]string{
				"name": "default",
			},
		},
		{
			Name: "stackrox",
			Id:   "stackrox",
			Labels: map[string]string{
				"name": "stackrox",
			},
		},
		{
			Name: "other",
			Id:   "other",
		},
	}

	namespacesByID = func() map[string]*storage.NamespaceMetadata {
		m := make(map[string]*storage.NamespaceMetadata)
		for _, ns := range namespaces {
			m[ns.GetId()] = ns
		}
		return m
	}()
)

type namespaceGetter struct{}

func (n *namespaceGetter) GetNamespaces(ctx context.Context) ([]*storage.NamespaceMetadata, error) {
	return namespaces, nil
}

func newMockGraphEvaluator() *evaluatorImpl {
	return newGraphEvaluator(&namespaceGetter{})
}

func TestDoesNamespaceMatchLabel(t *testing.T) {
	cases := []struct {
		name      string
		namespace *storage.NamespaceMetadata
		selector  *storage.LabelSelector
		expected  bool
	}{
		{
			name:      "No values in selector - no namespace labels",
			namespace: &storage.NamespaceMetadata{},
			selector:  &storage.LabelSelector{},
			expected:  true,
		},
		{
			name:      "No values in selector - some namespace labels",
			namespace: &storage.NamespaceMetadata{},
			selector:  &storage.LabelSelector{},
			expected:  true,
		},
		{
			name: "matching values in selector",
			namespace: &storage.NamespaceMetadata{
				Labels: map[string]string{
					"hello": "hi",
				},
			},
			selector: &storage.LabelSelector{
				MatchLabels: map[string]string{
					"hello": "hi",
				},
			},
			expected: true,
		},
		{
			name: "non matching values in selector",
			namespace: &storage.NamespaceMetadata{
				Labels: map[string]string{
					"hello": "hi1",
				},
			},
			selector: &storage.LabelSelector{
				MatchLabels: map[string]string{
					"hello": "hi",
				},
			},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cs, err := labels.CompileSelector(c.selector)
			require.NoError(t, err)
			assert.Equal(t, c.expected, cs.Matches(c.namespace.GetLabels()))
		})
	}
}

func TestDoesPodLabelsMatchLabel(t *testing.T) {
	cases := []struct {
		name       string
		deployment *storage.Deployment
		selector   *storage.LabelSelector
		expected   bool
	}{
		{
			name:       "No values in selector - no deployment labels",
			deployment: &storage.Deployment{},
			selector:   &storage.LabelSelector{},
			expected:   true,
		},
		{
			name:       "No values in selector - some deployment labels",
			deployment: &storage.Deployment{},
			selector:   &storage.LabelSelector{},
			expected:   true,
		},
		{
			name: "matching values in selector",
			deployment: &storage.Deployment{
				PodLabels: map[string]string{
					"hello": "hi",
				},
			},
			selector: &storage.LabelSelector{
				MatchLabels: map[string]string{
					"hello": "hi",
				},
			},
			expected: true,
		},
		{
			name: "non matching values in selector",
			deployment: &storage.Deployment{
				PodLabels: map[string]string{
					"hello": "hi1",
				},
			},
			selector: &storage.LabelSelector{
				MatchLabels: map[string]string{
					"hello": "hi",
				},
			},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cs, err := labels.CompileSelector(c.selector)
			require.NoError(t, err)
			assert.Equal(t, c.expected, cs.Matches(c.deployment.GetPodLabels()))
		})
	}
}

func getExamplePolicy(name string) *storage.NetworkPolicy {
	np, ok := networkPolicyFixtures[name]
	if !ok {
		panic(name)
	}
	return np
}

type testEdge struct {
	Source, Target string
	Ports          portDescs
}

func egressEdges(src string, dsts ...string) []testEdge {
	var edges []testEdge
	for _, d := range dsts {
		edges = append(edges, testEdge{Source: src, Target: d})
	}
	return edges
}

func egressEdgesWithPorts(src string, pds portDescs, dsts ...string) []testEdge {
	var edges []testEdge
	for _, d := range dsts {
		edges = append(edges, testEdge{Source: src, Target: d, Ports: pds})
	}
	return edges
}

func ingressEdges(dst string, srcs ...string) []testEdge {
	var edges []testEdge
	for _, s := range srcs {
		edges = append(edges, testEdge{Source: s, Target: dst})
	}
	return edges
}

func ingressEdgesWithPort(dst string, pds portDescs, srcs ...string) []testEdge {
	var edges []testEdge
	for _, s := range srcs {
		edges = append(edges, testEdge{Source: s, Target: dst, Ports: pds})
	}
	return edges
}

func fullyConnectedEdges(values ...string) []testEdge {
	var edges []testEdge
	for i, value1 := range values {
		for j, value2 := range values {
			if i == j {
				continue
			}
			edges = append(edges, testEdge{Source: value1, Target: value2})
		}
	}
	return edges
}

func flattenEdges(edges ...[]testEdge) []testEdge {
	var finalEdges []testEdge
	for _, e := range edges {
		finalEdges = append(finalEdges, e...)
	}
	return finalEdges
}

func mockNode(node string, namespace string, internetAccess, nonIsolatedIngress, nonIsolatedEgress bool, policies ...string) *v1.NetworkNode {
	sort.Strings(policies)
	return &v1.NetworkNode{
		Entity: &storage.NetworkEntityInfo{
			Type: storage.NetworkEntityInfo_DEPLOYMENT,
			Id:   node,
			Desc: &storage.NetworkEntityInfo_Deployment_{
				Deployment: &storage.NetworkEntityInfo_Deployment{
					Namespace: namespace,
				},
			},
		},
		PolicyIds:          policies,
		InternetAccess:     internetAccess,
		NonIsolatedIngress: nonIsolatedIngress,
		NonIsolatedEgress:  nonIsolatedEgress,
		OutEdges:           make(map[int32]*v1.NetworkEdgePropertiesBundle),
	}
}

func mockExternalNode(node string, cidr string) *v1.NetworkNode {
	return &v1.NetworkNode{
		Entity: &storage.NetworkEntityInfo{
			Type: storage.NetworkEntityInfo_EXTERNAL_SOURCE,
			Id:   node,
			Desc: &storage.NetworkEntityInfo_ExternalSource_{
				ExternalSource: &storage.NetworkEntityInfo_ExternalSource{
					Source: &storage.NetworkEntityInfo_ExternalSource_Cidr{
						Cidr: cidr,
					},
				},
			},
		},
		InternetAccess:     true,
		NonIsolatedIngress: true,
		NonIsolatedEgress:  true,
		OutEdges:           make(map[int32]*v1.NetworkEdgePropertiesBundle),
	}
}

func mockInternetNode() *v1.NetworkNode {
	return &v1.NetworkNode{
		Entity: &storage.NetworkEntityInfo{
			Type: storage.NetworkEntityInfo_INTERNET,
			Id:   networkgraph.InternetExternalSourceID,
		},
		InternetAccess:     true,
		NonIsolatedIngress: true,
		NonIsolatedEgress:  true,
		OutEdges:           make(map[int32]*v1.NetworkEdgePropertiesBundle),
	}
}

func deploymentLabels(values ...string) map[string]string {
	if len(values)%2 != 0 {
		panic("values for deployments labels must be even")
	}
	m := make(map[string]string)
	for i := 0; i < len(values)/2; i++ {
		m[values[2*i]] = values[2*i+1]
	}
	return m
}

func TestEvaluateClusters(t *testing.T) {
	g := newMockGraphEvaluator()

	// These are the k8s examples from https://github.com/ahmetb/kubernetes-network-policy-recipes
	// Seems like a good way to verify that the logic is correct
	cases := []struct {
		name         string
		deployments  []*storage.Deployment
		externalSrcs []*storage.NetworkEntityInfo
		nps          []*storage.NetworkPolicy
		edges        []testEdge
		nodes        []*v1.NetworkNode
	}{
		{
			name: "No policies - fully connected",
			deployments: []*storage.Deployment{
				{
					Id: "d1",
				},
				{
					Id: "d2",
				},
			},
			externalSrcs: []*storage.NetworkEntityInfo{
				{
					Id:   "es1",
					Type: storage.NetworkEntityInfo_EXTERNAL_SOURCE,
					Desc: &storage.NetworkEntityInfo_ExternalSource_{
						ExternalSource: &storage.NetworkEntityInfo_ExternalSource{
							Source: &storage.NetworkEntityInfo_ExternalSource_Cidr{
								Cidr: "172.17.0.0/24",
							},
						},
					},
				},
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "", true, true, true),
				mockNode("d2", "", true, true, true),
				mockInternetNode(),
				mockExternalNode("es1", "172.17.0.0/24"),
			},
		},
		{
			name: "deny all to app=web",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "web"),
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
				},
				{
					Id:          "d3",
					Namespace:   "default",
					NamespaceId: "default",
				},
			},
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("web-deny-all"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", true, false, true, "web-deny-all"),
				mockNode("d2", "default", true, true, true),
				mockNode("d3", "default", true, true, true),
				mockInternetNode(),
			},
		},
		{
			name: "limit traffic to application",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "bookstore", "role", "api"),
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "bookstore", "role", "frontend"),
				},
				{
					Id:          "d3",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "coffeeshop", "role", "api"),
				},
			},
			edges: flattenEdges(
				ingressEdges("d1", "d2"),
			),
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("limit-traffic"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", true, false, true, "limit-traffic"),
				mockNode("d2", "default", true, true, true),
				mockNode("d3", "default", true, true, true),
				mockInternetNode(),
			},
		},
		{
			name: "allow all ingress even if deny all",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "web"),
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
				},
				{
					Id:          "d3",
					Namespace:   "default",
					NamespaceId: "default",
				},
			},
			edges: flattenEdges(
				ingressEdges("d1", "d2", "d3"),
			),
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("web-deny-all"),
				getExamplePolicy("web-allow-all"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", true, false, true, "web-allow-all", "web-deny-all"),
				mockNode("d2", "default", true, true, true),
				mockNode("d3", "default", true, true, true),
				mockInternetNode(),
			},
		},
		{
			name: "DENY all non-whitelisted traffic to a namespace", // TODO: update to inclusive language when updating actual code
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "web"),
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
				},
				{
					Id:          "d3",
					Namespace:   "stackrox",
					NamespaceId: "stackrox",
				},
			},
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("default-deny-all"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", true, false, true, "default-deny-all"),
				mockNode("d2", "default", true, false, true, "default-deny-all"),
				mockNode("d3", "stackrox", true, true, true),
				mockInternetNode(),
			},
		},
		{
			name: "DENY all traffic from other namespaces",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
				},
				{
					Id:          "d3",
					Namespace:   "stackrox",
					NamespaceId: "stackrox",
				},
			},
			edges: flattenEdges(
				fullyConnectedEdges("d1", "d2"),
			),
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("deny-from-other-namespaces"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", true, false, true, "deny-from-other-namespaces"),
				mockNode("d2", "default", true, false, true, "deny-from-other-namespaces"),
				mockNode("d3", "stackrox", true, true, true),
				mockInternetNode(),
			},
		},
		{
			name: "Web allow all traffic from other namespaces",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "web"),
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
				},
				{
					Id:          "d3",
					Namespace:   "stackrox",
					NamespaceId: "stackrox",
				},
			},
			edges: flattenEdges(
				fullyConnectedEdges("d1", "d2"),
				ingressEdges("d1", "d3"),
			),
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("deny-from-other-namespaces"),
				getExamplePolicy("web-allow-all-namespaces"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", true, false, true, "deny-from-other-namespaces", "web-allow-all-namespaces"),
				mockNode("d2", "default", true, false, true, "deny-from-other-namespaces"),
				mockNode("d3", "stackrox", true, true, true),
				mockInternetNode(),
			},
		},
		{
			name: "Web allow all traffic from stackrox namespace",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "web"),
				},
				{
					Id:          "d2",
					Namespace:   "other",
					NamespaceId: "other",
				},
				{
					Id:          "d3",
					Namespace:   "stackrox",
					NamespaceId: "stackrox",
				},
			},
			edges: flattenEdges(
				ingressEdges("d1", "d3"),
			),
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("web-allow-stackrox"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", true, false, true, "web-allow-stackrox"),
				mockNode("d2", "other", true, true, true),
				mockNode("d3", "stackrox", true, true, true),
				mockInternetNode(),
			},
		},
		{
			name: "Allow traffic from apps using multiple selectors",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "web", "role", "db"),
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "bookstore", "role", "search"),
				},
				{
					Id:          "d3",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "bookstore", "role", "api"),
				},
				{
					Id:          "d4",
					Namespace:   "default",
					NamespaceId: "default",
				},
			},
			edges: flattenEdges(
				ingressEdges("d1", "d2", "d3"),
			),
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("web-deny-all"),
				getExamplePolicy("allow-traffic-from-apps-using-multiple-selectors"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", true, false, true, "web-deny-all", "allow-traffic-from-apps-using-multiple-selectors"),
				mockNode("d2", "default", true, true, true),
				mockNode("d3", "default", true, true, true),
				mockNode("d4", "default", true, true, true),
				mockInternetNode(),
			},
		},
		{
			name: "web deny egress",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "web"),
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
				},
			},
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("web-deny-egress"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", false, true, false, "web-deny-egress"),
				mockNode("d2", "default", true, true, true),
				mockInternetNode(),
			},
		},
		{
			name: "deny egress from namespace",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "web"),
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
				},
				{
					Id:          "d3",
					Namespace:   "stackrox",
					NamespaceId: "stackrox",
				},
			},
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("default-deny-all-egress"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", false, true, false, "default-deny-all-egress"),
				mockNode("d2", "default", false, true, false, "default-deny-all-egress"),
				mockNode("d3", "stackrox", true, true, true),
				mockInternetNode(),
			},
		},
		{
			name: "deny internetAccess egress from cluster",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "web"),
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
				},
				{
					Id:          "d3",
					Namespace:   "stackrox",
					NamespaceId: "stackrox",
				},
			},
			edges: flattenEdges(
				egressEdges("d1", "d2", "d3"),
			),
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("web-deny-external-egress"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", false, true, false, "web-deny-external-egress"),
				mockNode("d2", "default", true, true, true),
				mockNode("d3", "stackrox", true, true, true),
				mockInternetNode(),
			},
		},
		{
			name: "deny all ingress except for app = web",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "qa",
					NamespaceId: "qa",
					PodLabels:   deploymentLabels("app", "web"),
				},
				{
					Id:          "d2",
					Namespace:   "qa",
					NamespaceId: "qa",
					PodLabels:   deploymentLabels("app", "client"),
				},
				{
					Id:          "d3",
					Namespace:   "stackrox",
					NamespaceId: "stackrox",
				},
				{
					Id:          "d4",
					Namespace:   "default",
					NamespaceId: "default",
				},
			},
			edges: flattenEdges(
				ingressEdges("d1", "d2", "d3", "d4"),
			),
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("deny-all-ingress"),
				getExamplePolicy("allow-ingress-to-web"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "qa", true, false, true, "allow-ingress-to-web", "deny-all-ingress"),
				mockNode("d2", "qa", true, false, true, "deny-all-ingress"),
				mockNode("d3", "stackrox", true, true, true),
				mockNode("d4", "default", true, true, true),
				mockInternetNode(),
			},
		},
		{
			name: "fully isolate all pods",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
				},
			},
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("fully-isolate"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", false, false, false, "fully-isolate"),
				mockNode("d2", "default", false, false, false, "fully-isolate"),
			},
		},
		{
			name: "allow only egress to ipblock",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
				},
			},
			externalSrcs: []*storage.NetworkEntityInfo{
				{
					Id:   "es1",
					Type: storage.NetworkEntityInfo_EXTERNAL_SOURCE,
					Desc: &storage.NetworkEntityInfo_ExternalSource_{
						ExternalSource: &storage.NetworkEntityInfo_ExternalSource{
							Source: &storage.NetworkEntityInfo_ExternalSource_Cidr{
								Cidr: "172.17.0.0/24",
							},
						},
					},
				},
			},
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("allow-only-egress-to-ipblock"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", true, false, false, "allow-only-egress-to-ipblock"),
				mockInternetNode(),
				mockExternalNode("es1", "172.17.0.0/24"),
			},
			edges: flattenEdges(
				egressEdges("d1", "es1", networkgraph.InternetExternalSourceID),
			),
		},
		{
			name: "ingress and egress combination",
			deployments: []*storage.Deployment{
				{
					Id:          "a",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "a"),
				},
				{
					Id:          "b",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "b"),
				},
				{
					Id:          "c",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "c"),
				},
			},
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("a-ingress-tcp-8080"),
				getExamplePolicy("b-egress-a-tcp-ports-and-dns"),
				getExamplePolicy("c-egress-a-tcp-8443-and-udp"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("a", "default", true, false, true, "a-ingress-tcp-8080"),
				mockNode("b", "default", false, true, false, "b-egress-a-tcp-ports-and-dns"),
				mockNode("c", "default", false, true, false, "c-egress-a-tcp-8443-and-udp"),
				mockInternetNode(),
			},
			edges: flattenEdges(
				ingressEdges("a", "b"),
			),
		},
	}
	for _, c := range cases {
		testCase := c
		populateOutEdges(testCase.nodes, testCase.edges)
		t.Run(c.name, func(t *testing.T) {
			nodes := g.GetGraph("", testCase.deployments, testCase.externalSrcs, testCase.nps, false)
			assert.ElementsMatch(t, testCase.nodes, nodes.GetNodes())
		})
	}
}

func TestGetApplicable(t *testing.T) {
	g := newMockGraphEvaluator()

	// These are the k8s examples from https://github.com/ahmetb/kubernetes-network-policy-recipes
	// Seems like a good way to verify that the logic is correct
	cases := []struct {
		name        string
		deployments []*storage.Deployment
		policies    []*storage.NetworkPolicy
		expected    []*storage.NetworkPolicy
	}{
		{
			name: "No policies",
			deployments: []*storage.Deployment{
				{
					Id: "d1",
				},
				{
					Id: "d2",
				},
			},
		},
		{
			name: "deny all to app=web with match",
			deployments: []*storage.Deployment{
				{
					Id:        "d1",
					Namespace: "default",
					PodLabels: deploymentLabels("app", "web"),
				},
				{
					Id:        "d2",
					Namespace: "default",
				},
				{
					Id:        "d3",
					Namespace: "default",
				},
			},
			policies: []*storage.NetworkPolicy{
				getExamplePolicy("web-deny-all"),
			},
			expected: []*storage.NetworkPolicy{
				getExamplePolicy("web-deny-all"),
			},
		},
		{
			name: "limit traffic to application with match",
			deployments: []*storage.Deployment{
				{
					Id:        "d1",
					Namespace: "default",
					PodLabels: deploymentLabels("app", "bookstore", "role", "api"),
				},
				{
					Id:        "d2",
					Namespace: "default",
					PodLabels: deploymentLabels("app", "bookstore", "role", "frontend"),
				},
				{
					Id:        "d3",
					Namespace: "default",
					PodLabels: deploymentLabels("app", "coffeeshop", "role", "api"),
				},
			},
			policies: []*storage.NetworkPolicy{
				getExamplePolicy("limit-traffic"),
			},
			expected: []*storage.NetworkPolicy{
				getExamplePolicy("limit-traffic"),
			},
		},
		{
			name: "limit traffic to application no match",
			deployments: []*storage.Deployment{
				{
					Id:        "d1",
					Namespace: "default",
					PodLabels: deploymentLabels("app", "web"),
				},
				{
					Id:        "d2",
					Namespace: "default",
				},
				{
					Id:        "d3",
					Namespace: "default",
				},
			},
			policies: []*storage.NetworkPolicy{
				getExamplePolicy("limit-traffic"),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := g.GetAppliedPolicies(c.deployments, nil, c.policies)
			assert.ElementsMatch(t, c.expected, actual)
		})
	}
}

func populateOutEdges(nodes []*v1.NetworkNode, edges []testEdge) {
	indexMap := make(map[string]int)
	for i, node := range nodes {
		indexMap[node.Entity.Id] = i
	}

	for _, e := range edges {
		if e.Source == e.Target {
			continue
		}
		srcIndex := indexMap[e.Source]
		srcNode := nodes[srcIndex]
		tgtIndex := indexMap[e.Target]
		bundle := &v1.NetworkEdgePropertiesBundle{}
		pds := e.Ports.Clone()
		pds.normalizeInPlace()
		bundle.Properties = pds.ToProto()
		srcNode.OutEdges[int32(tgtIndex)] = bundle
	}
}

func TestEvaluateClustersWithPorts(t *testing.T) {
	g := newMockGraphEvaluator()

	cases := []struct {
		name        string
		deployments []*storage.Deployment
		nps         []*storage.NetworkPolicy
		edges       []testEdge
		nodes       []*v1.NetworkNode
	}{
		{
			name: "only allow port 5000 on API server",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "apiserver"),
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("role", "monitoring"),
				},
				{
					Id:          "d3",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("role", "other"),
				},
			},
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("api-allow-5000"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", true, false, true, "api-allow-5000"),
				mockNode("d2", "default", true, true, true),
				mockNode("d3", "default", true, true, true),
				mockInternetNode(),
			},
			edges: flattenEdges(
				ingressEdgesWithPort("d1", portDescs{{l4proto: storage.Protocol_TCP_PROTOCOL, port: 5000}}, "d2"),
			),
		},
		{
			name: "only allow DNS egress",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "apiserver"),
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("role", "app"),
				},
				{
					Id:          "d3",
					Namespace:   "kube-system",
					NamespaceId: "kube-system",
					PodLabels:   deploymentLabels("role", "kube-dns"),
				},
			},
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("allow-dns-egress-only"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", true, true, false, "allow-dns-egress-only"),
				mockNode("d2", "default", true, true, true),
				mockNode("d3", "kube-system", true, true, true),
				mockInternetNode(),
			},
			edges: flattenEdges(
				egressEdgesWithPorts("d1", portDescs{{l4proto: storage.Protocol_TCP_PROTOCOL, port: 53}, {l4proto: storage.Protocol_UDP_PROTOCOL, port: 53}}, "d2", "d3", networkgraph.InternetExternalSourceID),
			),
		},
		{
			name: "allow traffic on named API port",
			deployments: []*storage.Deployment{
				{
					Id:          "d1",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "apiserver"),
					Ports: []*storage.PortConfig{
						{
							Name:          "api",
							ContainerPort: 8443,
							Protocol:      "TCP",
						},
					},
				},
				{
					Id:          "d2",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "apiserver"),
				},
				{
					Id:          "d3",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("role", "monitoring"),
				},
			},
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("api-allow-named-api-port"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("d1", "default", true, false, true, "api-allow-named-api-port"),
				mockNode("d2", "default", true, false, true, "api-allow-named-api-port"),
				mockNode("d3", "default", true, true, true),
				mockInternetNode(),
			},
			edges: flattenEdges(
				ingressEdgesWithPort("d1", portDescs{{l4proto: storage.Protocol_TCP_PROTOCOL, port: 8443}}, "d3"),
			),
		},
		{
			name: "ingress and egress combination",
			deployments: []*storage.Deployment{
				{
					Id:          "a",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "a"),
				},
				{
					Id:          "b",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "b"),
				},
				{
					Id:          "c",
					Namespace:   "default",
					NamespaceId: "default",
					PodLabels:   deploymentLabels("app", "c"),
				},
			},
			nps: []*storage.NetworkPolicy{
				getExamplePolicy("a-ingress-tcp-8080"),
				getExamplePolicy("b-egress-a-tcp-ports-and-dns"),
				getExamplePolicy("c-egress-a-tcp-8443-and-udp"),
			},
			nodes: []*v1.NetworkNode{
				mockNode("a", "default", true, false, true, "a-ingress-tcp-8080"),
				mockNode("b", "default", false, true, false, "b-egress-a-tcp-ports-and-dns"),
				mockNode("c", "default", false, true, false, "c-egress-a-tcp-8443-and-udp"),
				mockInternetNode(),
			},
			edges: flattenEdges(
				ingressEdgesWithPort("a", portDescs{{l4proto: storage.Protocol_TCP_PROTOCOL, port: 8080}}, "b"),
			),
		},
	}
	for _, c := range cases {
		testCase := c
		populateOutEdges(testCase.nodes, testCase.edges)
		t.Run(c.name, func(t *testing.T) {
			nodes := g.GetGraph("", testCase.deployments, nil, testCase.nps, true)
			assert.ElementsMatch(t, testCase.nodes, nodes.GetNodes())
		})
	}
}
