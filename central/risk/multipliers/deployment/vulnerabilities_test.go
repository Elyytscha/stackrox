package deployment

import (
	"context"
	"testing"

	"github.com/stackrox/rox/central/risk/multipliers"
	"github.com/stretchr/testify/assert"
)

func TestVulnerabilitiesScore(t *testing.T) {
	ctx := context.Background()

	mult := NewVulnerabilities()
	deployment := multipliers.GetMockDeployment()
	images := multipliers.GetMockImages()
	result := mult.Score(ctx, deployment, images)
	assert.Equal(t, float32(1.3), result.Score)

	images[0].GetScan().GetComponents()[0].GetVulns()[0].Cvss = 0
	result = mult.Score(ctx, deployment, images)
	assert.Equal(t, float32(1.225), result.Score)

	images[0].GetScan().GetComponents()[0].GetVulns()[0].Cvss = 10
	result = mult.Score(ctx, deployment, images)
	assert.Equal(t, float32(1.525), result.Score)

	// Set both CVSS to 0 and then there should be a nil RiskResult
	images[0].GetScan().GetComponents()[0].GetVulns()[0].Cvss = 0
	images[0].GetScan().GetComponents()[0].GetVulns()[1].Cvss = 0
	images[0].GetScan().GetComponents()[1].GetVulns()[0].Cvss = 0
	images[0].GetScan().GetComponents()[1].GetVulns()[1].Cvss = 0

	result = mult.Score(ctx, deployment, images)
	assert.Nil(t, result)
}
