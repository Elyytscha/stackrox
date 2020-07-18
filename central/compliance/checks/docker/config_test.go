package docker

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stackrox/rox/central/compliance/framework"
	"github.com/stackrox/rox/central/compliance/framework/mocks"
	"github.com/stackrox/rox/generated/internalapi/compliance"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func nodeNameMatcher(nodeName string) gomock.Matcher {
	return testutils.PredMatcher(fmt.Sprintf("node %s", nodeName), func(node *storage.Node) bool {
		return node.GetName() == nodeName
	})
}

func TestOwnershipAndPermissionChecks(t *testing.T) {
	cases := []struct {
		name   string
		file   *compliance.File
		status framework.Status
	}{
		{
			name: "CIS_Docker_v1_2_0:3_1",
			file: &compliance.File{
				Path:      "docker.service",
				GroupName: "root",
				UserName:  "root",
			},
			status: framework.PassStatus,
		},
		{
			name: "CIS_Docker_v1_2_0:3_1",
			file: &compliance.File{
				Path:      "docker.service",
				GroupName: "docker",
				UserName:  "root",
			},
			status: framework.FailStatus,
		},
		{
			name: "CIS_Docker_v1_2_0:3_1",
			file: &compliance.File{
				Path:      "docker.service",
				GroupName: "root",
				UserName:  "docker",
			},
			status: framework.FailStatus,
		},
		{
			name: "CIS_Docker_v1_2_0:3_2",
			file: &compliance.File{
				Path:        "docker.service",
				Permissions: 0644,
			},
			status: framework.PassStatus,
		},
		{
			name: "CIS_Docker_v1_2_0:3_2",
			file: &compliance.File{
				Path:        "docker.service",
				Permissions: 0643,
			},
			status: framework.PassStatus,
		},
		{
			name: "CIS_Docker_v1_2_0:3_2",
			file: &compliance.File{
				Path:        "docker.service",
				Permissions: 0645,
			},
			status: framework.FailStatus,
		},
	}

	for _, cIt := range cases {
		c := cIt
		t.Run(strings.Replace(c.name, ":", "-", -1), func(t *testing.T) {
			if features.ComplianceInNodes.Enabled() {
				return
			}
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

			domain := framework.NewComplianceDomain(testCluster, testNodes, nil, nil)
			data := mocks.NewMockComplianceDataRepository(mockCtrl)

			data.EXPECT().HostScraped(nodeNameMatcher("A")).AnyTimes().Return(&compliance.ComplianceReturn{
				SystemdFiles: map[string]*compliance.File{
					c.file.Path: c.file,
				},
			})

			data.EXPECT().HostScraped(nodeNameMatcher("B")).AnyTimes().Return(&compliance.ComplianceReturn{
				SystemdFiles: map[string]*compliance.File{
					c.file.Path: c.file,
				},
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
