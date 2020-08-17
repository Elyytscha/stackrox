package check85

import (
	"github.com/stackrox/rox/central/compliance/checks/common"
	"github.com/stackrox/rox/central/compliance/framework"
	pkgFramework "github.com/stackrox/rox/pkg/compliance/framework"
	"github.com/stackrox/rox/pkg/features"
)

const checkID = "PCI_DSS_3_2:8_5"

func init() {
	framework.MustRegisterNewCheck(
		framework.CheckMetadata{
			ID:                 checkID,
			Scope:              pkgFramework.ClusterKind,
			AdditionalScopes:   []pkgFramework.TargetKind{pkgFramework.DeploymentKind},
			DataDependencies:   []string{"Deployments", "K8sRoles", "K8sRoleBindings", "HostScraped"},
			InterpretationText: interpretationText,
		},
		clusterIsCompliant)
}

func clusterIsCompliant(ctx framework.ComplianceContext) {
	if !features.ComplianceInNodes.Enabled() {
		common.IsRBACConfiguredCorrectly(ctx)
	}
	common.LimitedUsersAndGroupsWithClusterAdmin(ctx)
	common.AdministratorUsersPresent(ctx)
	common.CheckDeploymentsDoNotHaveClusterAccess(ctx, common.EffectiveAdmin)
}
