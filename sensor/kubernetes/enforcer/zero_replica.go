package enforcer

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/internalapi/central"
	pkgKubernetes "github.com/stackrox/rox/pkg/kubernetes"
	"github.com/stackrox/rox/sensor/kubernetes/enforcer/common"
	"github.com/stackrox/rox/sensor/kubernetes/enforcer/daemonset"
	"github.com/stackrox/rox/sensor/kubernetes/enforcer/deployment"
	"github.com/stackrox/rox/sensor/kubernetes/enforcer/replicaset"
	"github.com/stackrox/rox/sensor/kubernetes/enforcer/replicationcontroller"
	"github.com/stackrox/rox/sensor/kubernetes/enforcer/statefulset"
)

func (e *enforcerImpl) scaleToZero(enforcement *central.SensorEnforcement) (err error) {
	deploymentInfo := enforcement.GetDeployment()
	if deploymentInfo == nil {
		return errors.New("unable to apply constraint to non-deployment")
	}

	// Set enforcement function based on deployment type.
	var function func() error
	switch deploymentInfo.GetDeploymentType() {
	case pkgKubernetes.Deployment:
		function = func() error {
			return deployment.EnforceZeroReplica(e.client, deploymentInfo)
		}
	case pkgKubernetes.DaemonSet:
		function = func() error {
			return daemonset.EnforceZeroReplica(e.client, deploymentInfo)
		}
	case pkgKubernetes.ReplicaSet:
		function = func() error {
			return replicaset.EnforceZeroReplica(e.client, deploymentInfo)
		}
	case pkgKubernetes.ReplicationController:
		function = func() error {
			return replicationcontroller.EnforceZeroReplica(e.client, deploymentInfo)
		}
	case pkgKubernetes.StatefulSet:
		function = func() error {
			return statefulset.EnforceZeroReplica(e.client, deploymentInfo)
		}
	default:
		return fmt.Errorf("unknown type: %s", deploymentInfo.GetDeploymentType())
	}

	// Retry any retriable errors encountered when trying to run the enforcement function.
	err = withReasonableRetry(function)
	if err != nil {
		return
	}

	// Mark the deployment as having been scaled to zero.
	return withReasonableRetry(func() error {
		return common.MarkScaledToZero(e.recorder, enforcement.GetDeployment().GetPolicyName(), getRef(enforcement))
	})
}
