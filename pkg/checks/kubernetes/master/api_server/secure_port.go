package apiserver

import (
	"strconv"

	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/checks"
	"github.com/stackrox/rox/pkg/checks/utils"
)

type securePort struct{}

func (a *securePort) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS Kubernetes v1.2.0 - 1.1.7",
			Description: "Ensure that the --secure-port argument is not set to 0",
		}, Dependencies: []utils.Dependency{utils.InitKubeAPIServerConfig},
	}
}

func (a *securePort) Run() (result v1.CheckResult) {
	if params, ok := utils.KubeAPIServerConfig["secure-port"]; ok {
		port, err := strconv.Atoi(params.String())
		if err != nil || port < 1 || port > 65535 {
			utils.Warn(&result)
			utils.AddNotef(&result, "secure-port on kube-apiserver is set to '%v', but it must be set to a valid integer between 1-65535", params.String())
			return
		}
	}
	utils.Pass(&result)
	return
}

// NewSecurePort implements CIS Kubernetes v1.2.0 1.1.7
func NewSecurePort() utils.Check {
	return &securePort{}
}

func init() {
	checks.AddToRegistry(
		NewSecurePort(),
	)
}
