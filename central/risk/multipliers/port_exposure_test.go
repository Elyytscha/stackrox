package multipliers

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/risk"
	"github.com/stretchr/testify/assert"
)

func TestPortExposureScore(t *testing.T) {
	portMultiplier := NewReachability()

	deployment := getMockDeployment()
	expectedScore := &storage.Risk_Result{
		Name: risk.PortExposure.DisplayTitle,
		Factors: []*storage.Risk_Result_Factor{
			{Message: "Port 22 is exposed to external clients"},
			{Message: "Port 23 is exposed in the cluster"},
			{Message: "Port 24 is exposed on node interfaces"},
		},
		Score: 1.6,
	}
	score := portMultiplier.Score(context.Background(), deployment)
	assert.Equal(t, expectedScore, score)
}
