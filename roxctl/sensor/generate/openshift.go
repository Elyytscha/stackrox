package generate

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stackrox/rox/generated/storage"
	clusterValidation "github.com/stackrox/rox/pkg/cluster"
	"github.com/stackrox/rox/pkg/pointers"
	"github.com/stackrox/rox/roxctl/common/flags"
	"github.com/stackrox/rox/roxctl/common/util"
)

const (
	errorNotSupportedOnOpenShift3x   = `ERROR: The --admission-controller-listen-on-events flag is not supported for OpenShift 3.11. Set --openshift-version=4 to indicate that you are deploying on OpenShift 4.x in order to use this flag.`
	noteOpenShift3xCompatibilityMode = `NOTE: Deployment files are generated in OpenShift 3.x compatibility mode. Set the --openshift-version flag to 3 to suppress this note, or to 4 to take advantage of OpenShift 4.x features.`
)

func openshift() *cobra.Command {
	var openshiftVersion int
	var admissionControllerEvents *bool

	c := &cobra.Command{
		Use: "openshift",
		RunE: util.RunENoArgs(func(c *cobra.Command) error {
			cluster.Type = storage.ClusterType_OPENSHIFT_CLUSTER
			switch openshiftVersion {
			case 0:
				fmt.Fprintf(os.Stderr, "%s\n\n", noteOpenShift3xCompatibilityMode)
			case 3:
			case 4:
				cluster.Type = storage.ClusterType_OPENSHIFT4_CLUSTER
			default:
				return errors.Errorf("invalid OpenShift version %d, supported values are '3' and '4'", openshiftVersion)
			}

			if admissionControllerEvents == nil {
				admissionControllerEvents = pointers.Bool(cluster.Type == storage.ClusterType_OPENSHIFT4_CLUSTER) // enable for OpenShift 4 only
			} else if *admissionControllerEvents && cluster.Type == storage.ClusterType_OPENSHIFT_CLUSTER {
				// The below `Validate` call would also catch this, but catching it here allows us to print more
				// CLI-relevant error messages that reference flag names.
				fmt.Fprintf(os.Stderr, "%s\n\n", errorNotSupportedOnOpenShift3x)
				return errors.New("incompatible flag settings")
			}

			cluster.AdmissionControllerEvents = *admissionControllerEvents

			if err := clusterValidation.Validate(&cluster).ToError(); err != nil {
				return err
			}
			return fullClusterCreation(flags.Timeout(c))
		}),
	}
	c.PersistentFlags().IntVar(&openshiftVersion, "openshift-version", 0, "OpenShift major version to generate deployment files for")
	flags.OptBoolFlagVarPF(c.PersistentFlags(), &admissionControllerEvents, "admission-controller-listen-on-events", "", "enable admission controller webhook to listen on Kubernetes events", "auto")
	return c
}
