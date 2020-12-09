package imageintegration

import (
	"github.com/stackrox/rox/central/integrationhealth/reporter"
	"github.com/stackrox/rox/pkg/images/integration"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	once sync.Once

	is integration.Set
	no integration.ToNotify
)

func initialize() {
	// This is the set of image integrations currently active, and the ToNotify that updates that set.
	is = integration.NewSet(reporter.Singleton())
	no = integration.NewToNotify(is)
}

// Set provides the set of image integrations currently in use by central.
func Set() integration.Set {
	once.Do(initialize)
	return is
}

// ToNotify returns the ToNotify instance that updates the integration Set.
func ToNotify() integration.ToNotify {
	once.Do(initialize)
	return no
}
