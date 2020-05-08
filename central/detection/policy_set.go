package detection

import (
	policyDatastore "github.com/stackrox/rox/central/policy/datastore"
	"github.com/stackrox/rox/pkg/detection"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	log = logging.LoggerForModule()
)

// PolicySet is a set of policies.
type PolicySet interface {
	detection.PolicySet

	RemoveNotifier(notifierID string) error
}

// NewPolicySet returns a new instance of a PolicySet.
func NewPolicySet(store policyDatastore.DataStore, compiler detection.PolicyCompiler) PolicySet {
	return &setImpl{
		PolicySet:   detection.NewPolicySet(compiler),
		policyStore: store,
	}
}
