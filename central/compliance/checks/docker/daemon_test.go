package docker

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/golang/mock/gomock"
	"github.com/stackrox/rox/central/compliance/framework"
	"github.com/stackrox/rox/central/compliance/framework/mocks"
	"github.com/stackrox/rox/generated/internalapi/compliance"
	"github.com/stackrox/rox/generated/storage"
	internalTypes "github.com/stackrox/rox/pkg/docker/types"
	"github.com/stackrox/rox/pkg/netutil"
	"github.com/stackrox/rox/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	localIPNet = registry.NetIPNet(*netutil.MustParseCIDR("127.0.0.0/16"))

	nonLocalIPNet = registry.NetIPNet(*netutil.MustParseCIDR("0.0.0.0/24"))
)

func TestDockerInfoBasedChecks(t *testing.T) {
	cases := []struct {
		name   string
		info   types.Info
		status framework.Status
	}{
		{
			name: "CIS_Docker_v1_2_0:2_5",
			info: types.Info{
				Driver: "aufs",
			},
			status: framework.FailStatus,
		},
		{
			name: "CIS_Docker_v1_2_0:2_5",
			info: types.Info{
				Driver: "overlay2",
			},
			status: framework.PassStatus,
		},
		{
			name: "CIS_Docker_v1_2_0:2_15",
			info: types.Info{
				SecurityOptions: []string{"hello", "seccomp=default"},
			},
			status: framework.NoteStatus,
		},
		{
			name: "CIS_Docker_v1_2_0:2_16",
			info: types.Info{
				ExperimentalBuild: true,
			},
			status: framework.FailStatus,
		},
		{
			name:   "CIS_Docker_v1_2_0:2_16",
			info:   types.Info{},
			status: framework.PassStatus,
		},
		{
			name:   "CIS_Docker_v1_2_0:2_4",
			status: framework.PassStatus,
		},
		{
			name: "CIS_Docker_v1_2_0:2_4",
			info: types.Info{
				RegistryConfig: &registry.ServiceConfig{
					InsecureRegistryCIDRs: []*registry.NetIPNet{&localIPNet},
				},
			},
			status: framework.PassStatus,
		},
		{
			name: "CIS_Docker_v1_2_0:2_4",
			info: types.Info{
				RegistryConfig: &registry.ServiceConfig{
					InsecureRegistryCIDRs: []*registry.NetIPNet{&localIPNet, &nonLocalIPNet},
				},
			},
			status: framework.FailStatus,
		},
		{
			name: "CIS_Docker_v1_2_0:2_13",
			info: types.Info{
				LiveRestoreEnabled: true,
			},
			status: framework.PassStatus,
		},
		{
			name:   "CIS_Docker_v1_2_0:2_13",
			status: framework.FailStatus,
		},
		{
			name: "CIS_Docker_v1_2_0:2_12",
			info: types.Info{
				LoggingDriver: "json-file",
			},
			status: framework.FailStatus,
		},
		{
			name:   "CIS_Docker_v1_2_0:2_12",
			status: framework.PassStatus,
		},
	}

	for _, cIt := range cases {
		c := cIt
		t.Run(strings.Replace(c.name, ":", "-", -1), func(t *testing.T) {
			t.Parallel()

			registry := framework.RegistrySingleton()
			check := registry.Lookup(c.name)
			require.NotNil(t, check)

			testCluster := &storage.Cluster{
				Id: uuid.NewV4().String(),
			}
			testNodes := createTestNodes("A", "B")

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			domain := framework.NewComplianceDomain(testCluster, testNodes, nil)
			data := mocks.NewMockComplianceDataRepository(mockCtrl)

			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			err := json.NewEncoder(gz).Encode(&internalTypes.Data{
				Info: c.info,
			})
			require.NoError(t, err)
			require.NoError(t, gz.Close())

			data.EXPECT().HostScraped(nodeNameMatcher("A")).AnyTimes().Return(&compliance.ComplianceReturn{
				DockerData: &compliance.GZIPDataChunk{Gzip: buf.Bytes()},
			})
			data.EXPECT().HostScraped(nodeNameMatcher("B")).AnyTimes().Return(&compliance.ComplianceReturn{
				DockerData: &compliance.GZIPDataChunk{Gzip: buf.Bytes()},
			})

			run, err := framework.NewComplianceRun(check)
			require.NoError(t, err)
			err = run.Run(context.Background(), domain, data)
			require.NoError(t, err)

			results := run.GetAllResults()
			checkResults := results[c.name]
			require.NotNil(t, checkResults)

			require.Len(t, checkResults.Evidence(), 0)
			for _, node := range domain.Nodes() {
				nodeResults := checkResults.ForChild(node)
				require.NoError(t, nodeResults.Error())
				require.Len(t, nodeResults.Evidence(), 1)
				assert.Equal(t, c.status, nodeResults.Evidence()[0].Status)
			}
		})
	}
}
