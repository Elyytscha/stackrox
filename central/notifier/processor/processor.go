package processor

import (
	"github.com/stackrox/rox/central/notifier/store"
	"github.com/stackrox/rox/central/notifiers"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	log = logging.LoggerForModule()
)

// Processor is the interface for processing benchmarks, notifiers, and policies.
//go:generate mockgen-wrapper Processor
type Processor interface {
	ProcessAlert(alert *storage.Alert)
	ProcessAuditMessage(msg *v1.Audit_Message)

	HasNotifiers() bool
	UpdateNotifier(notifier notifiers.Notifier)
	RemoveNotifier(id string)

	GetIntegratedPolicies(notifierID string) (output []*storage.Policy)
	UpdatePolicy(policy *storage.Policy)
	RemovePolicy(policy *storage.Policy)
}

// New returns a new Processor
func New(s store.Store) (Processor, error) {
	processor := &processorImpl{
		notifiers:           make(map[string]notifiers.Notifier),
		notifiersToPolicies: make(map[string]map[string]*storage.Policy),
		storage:             s,
	}
	err := processor.initializeNotifiers()
	return processor, err
}
