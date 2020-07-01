package check432

import (
	"github.com/stackrox/rox/central/compliance/checks/common"
	"github.com/stackrox/rox/central/compliance/checks/kubernetes"
	"github.com/stackrox/rox/central/compliance/framework"
)

const (
	standardID = "NIST_800_190:4_3_2"
)

func init() {
	framework.MustRegisterNewCheck(
		framework.CheckMetadata{
			ID:                 standardID,
			Scope:              framework.ClusterKind,
			AdditionalScopes:   []framework.TargetKind{framework.DeploymentKind},
			DataDependencies:   []string{"Deployments", "K8sRoles", "K8sRoleBindings"},
			InterpretationText: interpretationText,
		},
		clusterIsCompliant)
}

func clusterIsCompliant(ctx framework.ComplianceContext) {
	kubernetes.MasterAPIServerCommandLine(standardID, "authorization-mode", "RBAC", "RBAC", common.Contains).Run(ctx)
	common.LimitedUsersAndGroupsWithClusterAdmin(ctx)
	common.CheckDeploymentsDoNotHaveClusterAccess(ctx, common.EffectiveAdmin)
}
