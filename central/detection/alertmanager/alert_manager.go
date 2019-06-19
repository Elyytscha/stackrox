package alertmanager

import (
	"context"

	alertDataStore "github.com/stackrox/rox/central/alert/datastore"
	"github.com/stackrox/rox/central/detection/runtime"
	notifierProcessor "github.com/stackrox/rox/central/notifier/processor"
	"github.com/stackrox/rox/generated/storage"
)

// AlertManager is a simplified interface for fetching and updating alerts.
type AlertManager interface {
	// AlertAndNotify takes in a list of alerts being produced, and a bunch of filters that specify what subset of alerts
	// we're looking at. It then pulls out the alerts matching the filters, and compares the alerts in the DB with the ones
	// that have been produced, and takes care of the logic of marking alerts no longer being produced as resolved,
	// notifying of new alerts, and updating the timestamp of updated alerts.
	// It returns true if it has actually added/removed/updated alerts.
	AlertAndNotify(ctx context.Context, alerts []*storage.Alert, oldAlertFilters ...AlertFilterOption) (modified bool, err error)
}

// New returns a new instance of AlertManager. You should just use the singleton instance instead.
func New(notifier notifierProcessor.Processor, alerts alertDataStore.DataStore, detector runtime.Detector) AlertManager {
	return &alertManagerImpl{
		notifier:        notifier,
		alerts:          alerts,
		runtimeDetector: detector,
	}
}
