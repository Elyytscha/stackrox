package storagetov2

import (
	"github.com/stackrox/rox/central/complianceoperator/v2/checkresults/datastore"
	v2 "github.com/stackrox/rox/generated/api/v2"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	log = logging.LoggerForModule()
)

func ComplianceV2CheckResult(incoming *storage.ComplianceOperatorCheckResultV2) *v2.ComplianceCheckResult {

	converted := &v2.ComplianceCheckResult{
		CheckId:      incoming.GetCheckId(),
		CheckName:    incoming.GetCheckName(),
		Description:  incoming.GetDescription(),
		Instructions: incoming.GetInstructions(),
		Clusters: []*v2.ComplianceCheckResult_ClusterCheckStatus{
			clusterStatus(incoming),
		},
	}

	return converted
}

func ComplianceV2CheckResults(incoming []*storage.ComplianceOperatorCheckResultV2) []*v2.ComplianceScanResult {
	// Map of []*ComplianceCheckResult by scan name
	// append check result to the map
	// build slice of results by scan name
	// build slice of status by check

	resultsByScan := make(map[string][]*v2.ComplianceCheckResult)
	resultsByScanByCheck := make(map[string]map[string][]*v2.ComplianceCheckResult_ClusterCheckStatus)
	var orderedKeys []string

	for _, result := range incoming {
		workingResult, found := resultsByScan[result.GetScanConfigId()]
		// First time seeing this scan config in the results.
		if !found {
			orderedKeys = append(orderedKeys, result.GetScanConfigId())
			resultsByScan[result.GetScanConfigId()] = []*v2.ComplianceCheckResult{
				ComplianceV2CheckResult(result),
			}
			resultsByScanByCheck[result.GetScanConfigId()][result.GetCheckName()] =
		} else {
			workingResult.
				resultsByScan[result.GetScanConfigId()] = append(resultsByScan[result.GetScanConfigId()], ComplianceV2CheckResult(result))
		}
	}

	var convertedResults []*v2.ComplianceScanResult
	for _, key := range orderedKeys {
		convertedResults = append(convertedResults, &v2.ComplianceScanResult{
			ScanName:     key,
			ProfileName:  "test for now",
			CheckResults: resultsByScan[key],
		})
	}

	log.Infof("SHREWS -- convertResultStorageSliceToV2 -- %v", convertedResults)

	return convertedResults
}

func ComplianceV2ClusterStats(resultCounts []*datastore.ResourceCountByResultByCluster) []*v2.ComplianceClusterScanStats {
	convertedResults := make([]*v2.ComplianceClusterScanStats, len(resultCounts))

	for _, resultCount := range resultCounts {
		convertedResults = append(convertedResults, &v2.ComplianceClusterScanStats{
			Cluster: &v2.ComplianceScanCluster{
				ClusterId:   resultCount.ClusterID,
				ClusterName: resultCount.ClusterName,
			},
			ScanStats: &v2.ComplianceScanStatsShim{
				ScanName: resultCount.ScanConfigName,
				CheckStats: []*v2.ComplianceScanStatsShim_ComplianceCheckStatusCount{
					{
						Count:  int32(resultCount.FailCount),
						Status: v2.ComplianceCheckStatus_FAIL,
					},
					{
						Count:  int32(resultCount.InfoCount),
						Status: v2.ComplianceCheckStatus_INFO,
					},
					{
						Count:  int32(resultCount.PassCount),
						Status: v2.ComplianceCheckStatus_PASS,
					},
					{
						Count:  int32(resultCount.ErrorCount),
						Status: v2.ComplianceCheckStatus_ERROR,
					},
					{
						Count:  int32(resultCount.ManualCount),
						Status: v2.ComplianceCheckStatus_MANUAL,
					},
					{
						Count:  int32(resultCount.InconsistentCount),
						Status: v2.ComplianceCheckStatus_INCONSISTENT,
					},
					{
						Count:  int32(resultCount.NotApplicableCount),
						Status: v2.ComplianceCheckStatus_NOT_APPLICABLE,
					},
				},
			},
		})
	}
	return nil
}

func clusterStatus(incoming *storage.ComplianceOperatorCheckResultV2) *v2.ComplianceCheckResult_ClusterCheckStatus {
	return &v2.ComplianceCheckResult_ClusterCheckStatus{
		Cluster: &v2.ComplianceScanCluster{
			ClusterId:   incoming.GetClusterId(),
			ClusterName: incoming.GetClusterName(),
		},
		Status:      convertComplianceCheckStatus(incoming.Status),
		CreatedTime: incoming.GetCreatedTime(),
	}
}
