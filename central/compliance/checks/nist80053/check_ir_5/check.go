package checkir5

import (
	"github.com/stackrox/rox/central/compliance/checks/common"
	"github.com/stackrox/rox/central/compliance/framework"
	"github.com/stackrox/rox/generated/storage"
)

const (
	controlID = "NIST_SP_800_53:IR_5"

	interpretationText = `TODO`
)

func init() {
	framework.MustRegisterNewCheck(
		framework.CheckMetadata{
			ID:                 controlID,
			Scope:              framework.ClusterKind,
			DataDependencies:   []string{"Policies"},
			InterpretationText: interpretationText,
		},
		func(ctx framework.ComplianceContext) {
			framework.Pass(ctx, "StackRox is installed") // TODO(viswa): Get text
			common.CheckAnyPolicyInLifeCycle(ctx, storage.LifecycleStage_RUNTIME)
		})
}
