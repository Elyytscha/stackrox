// +build nobazel

package central

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/image/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func getBaseConfig() Config {
	return Config{
		ClusterType: v1.ClusterType_KUBERNETES_CLUSTER,
		K8sConfig: &K8sConfig{
			CommonConfig: CommonConfig{
				MainImage:     "stackrox/main:2.2.11.0-57-g392c0f5bed-dirty",
				ClairifyImage: "stackrox.io/clairify:0.4.2",
			},
			Namespace: "stackrox",
		},
		SecretsByteMap: map[string][]byte{
			"ca.pem":                     {1},
			"ca-key.pem":                 {1},
			"jwt-key.der":                {1},
			"monitoring-client-cert.pem": {1},
			"monitoring-client-key.pem":  {1},
		},
	}
}

func TestRender(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(renderSuite))
}

type renderSuite struct {
	suite.Suite
	*kubernetes
}

func (suite *renderSuite) SetupSuite() {
	suite.kubernetes = &kubernetes{}
	dir := templates.Directory()
	monitoringChartPath = filepath.Join(dir, monitoringChartSuffix)
	centralChartPath = filepath.Join(dir, centralChartSuffix)
	clairifyChartPath = filepath.Join(dir, clairifyChartSuffix)
	dockerAuthPath = filepath.Join(filepath.Dir(dir), "assets/docker-auth.sh")
}

func (suite *renderSuite) testWithHostPath(t *testing.T, c Config) {
	c.HostPath = &HostPathPersistence{
		Name:     "stackrox-db",
		HostPath: "/var/lib/stackrox",
	}
	_, err := suite.Render(c)
	assert.NoError(t, err)

	c.HostPath = &HostPathPersistence{
		Name:              "stackrox-db",
		HostPath:          "/var/lib/stackrox",
		NodeSelectorKey:   "key",
		NodeSelectorValue: "value",
	}
	_, err = suite.Render(c)
	assert.NoError(t, err)
}

func (suite *renderSuite) testWithPV(t *testing.T, c Config) {
	c.External = &ExternalPersistence{
		Name: "name",
	}
	_, err := suite.Render(c)
	assert.NoError(t, err)

	c.External = &ExternalPersistence{
		Name:         "name",
		StorageClass: "storageClass",
	}
	_, err = suite.Render(c)
	assert.NoError(t, err)
}

func (suite *renderSuite) testWithLoadBalancers(t *testing.T, c Config) {
	c.K8sConfig.LoadBalancerType = v1.LoadBalancerType_NODE_PORT
	_, err := suite.Render(c)
	assert.NoError(t, err)

	c.K8sConfig.LoadBalancerType = v1.LoadBalancerType_LOAD_BALANCER
	_, err = suite.Render(c)
	assert.NoError(t, err)
}

func (suite *renderSuite) TestRenderMultiple() {
	for _, orch := range []v1.ClusterType{v1.ClusterType_KUBERNETES_CLUSTER, v1.ClusterType_OPENSHIFT_CLUSTER} {
		for _, format := range []v1.DeploymentFormat{v1.DeploymentFormat_KUBECTL, v1.DeploymentFormat_HELM} {
			suite.T().Run(fmt.Sprintf("%s-%s", orch, format), func(t *testing.T) {
				conf := getBaseConfig()
				conf.ClusterType = orch
				conf.K8sConfig.DeploymentFormat = format

				suite.testWithHostPath(t, conf)
				suite.testWithPV(t, conf)
				suite.testWithLoadBalancers(t, conf)

				suite.testWithMonitoring(t, conf)
			})
		}
	}
}

func (suite *renderSuite) testWithMonitoring(t *testing.T, c Config) {
	c.K8sConfig.MonitoringType = OnPrem
	c.K8sConfig.MonitoringEndpoint = "monitoring.stackrox:443"
	_, err := suite.Render(c)
	suite.NoError(err)

	c.K8sConfig.MonitoringLoadBalancerType = v1.LoadBalancerType_NODE_PORT
	_, err = suite.Render(c)
	suite.NoError(err)

	c.K8sConfig.MonitoringLoadBalancerType = v1.LoadBalancerType_LOAD_BALANCER
	_, err = suite.Render(c)
	suite.NoError(err)
}
