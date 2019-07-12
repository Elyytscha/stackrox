package processor

import (
	"github.com/stackrox/rox/central/notifiers"
	v1 "github.com/stackrox/rox/generated/api/v1"
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
	HasEnabledAuditNotifiers() bool
	UpdateNotifier(notifier notifiers.Notifier)
	RemoveNotifier(id string)

	UpdatePolicy(policy *storage.Policy)
	RemovePolicy(policy *storage.Policy)
}

// New returns a new Processor
func New(pns policyNotifierSet) Processor {
	return &processorImpl{
		pns: pns,
	}
}
