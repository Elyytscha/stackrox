package extensions

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/joelanford/helm-operator/pkg/extensions"
	"github.com/pkg/errors"
	securedClusterv1Alpha1 "github.com/stackrox/rox/operator/api/securedcluster/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

// CheckClusterNameExtension is an extension that ensures the spec.clusterName and status.clusterName fields are
// in sync.
func CheckClusterNameExtension(k8sClient kubernetes.Interface) extensions.ReconcileExtension {
	return wrapExtension(checkClusterName, k8sClient)
}

func checkClusterName(ctx context.Context, sc *securedClusterv1Alpha1.SecuredCluster, _ kubernetes.Interface, statusUpdater func(statusFunc updateStatusFunc), _ logr.Logger) error {
	if sc.DeletionTimestamp != nil {
		return nil // doesn't matter on deletion
	}
	if sc.Spec.ClusterName == "" {
		return errors.New("spec.clusterName is a required field")
	}
	if sc.Spec.ClusterName == sc.Status.ClusterName {
		return nil
	}
	if sc.Status.ClusterName != "" {
		return errors.Errorf("SecuredCluster instance was initially created with clusterName %q, but current value is %q. "+
			"Renaming clusters is not supported - you need to delete this object, and then recreate one with the correct cluster name.",
			sc.Status.ClusterName, sc.Spec.ClusterName)
	}

	statusUpdater(func(status *securedClusterv1Alpha1.SecuredClusterStatus) bool {
		status.ClusterName = sc.Spec.ClusterName
		return true
	})
	return nil
}
