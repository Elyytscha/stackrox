package processor

import (
	"context"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	pkgNotifier "github.com/stackrox/rox/pkg/notifier"
	"github.com/stackrox/rox/pkg/notifiers"
	"github.com/stackrox/rox/pkg/set"
)

var (
	// Replacing with a background context such that outside context cancellation
	// does not affect long-running go routines.
	ctxBackground = context.Background()
	log           = logging.LoggerForModule()
)

// Processor takes in alerts and sends the notifications tied to that alert
type processorImpl struct {
	ns pkgNotifier.Set
}

func (p *processorImpl) HasNotifiers() bool {
	return p.ns.HasNotifiers()
}

func (p *processorImpl) HasEnabledAuditNotifiers() bool {
	return p.ns.HasEnabledAuditNotifiers()
}

// RemoveNotifier removes the in memory copy of the specified notifier
func (p *processorImpl) RemoveNotifier(ctx context.Context, id string) {
	p.ns.RemoveNotifier(ctx, id)
}

// GetNotifier gets the in memory copy of the specified notifier id
func (p *processorImpl) GetNotifier(ctx context.Context, id string) (notifier notifiers.Notifier) {
	return p.ns.GetNotifier(ctx, id)
}

// GetNotifiers gets the in memory copies of all notifiers
func (p *processorImpl) GetNotifiers(ctx context.Context) (notifiers []notifiers.Notifier) {
	return p.ns.GetNotifiers(ctx)
}

// UpdateNotifier updates or adds the passed notifier into memory
func (p *processorImpl) UpdateNotifier(ctx context.Context, notifier notifiers.Notifier) {
	p.ns.UpsertNotifier(ctx, notifier)
}

// ProcessAlert pushes the alert into a channel to be processed
func (p *processorImpl) ProcessAlert(ctx context.Context, alert *storage.Alert) {
	if len(alert.GetPolicy().GetNotifiers()) == 0 {
		return
	}
	alertNotifiers := set.NewStringSet(alert.GetPolicy().GetNotifiers()...)

	p.ns.ForEach(ctx, func(ctx context.Context, notifier notifiers.Notifier, failures pkgNotifier.AlertSet) {
		if alertNotifiers.Contains(notifier.ProtoNotifier().GetId()) {
			go func() {
				err := pkgNotifier.TryToAlert(ctx, notifier, alert)
				if err != nil {
					failures.Add(alert)
				}
			}()
		}
	})
}

// ProcessAuditMessage sends the audit message with all applicable notifiers.
func (p *processorImpl) ProcessAuditMessage(ctx context.Context, msg *v1.Audit_Message) {
	panic("unimplemented")

}

func (p *processorImpl) UpdateNotifierHealthStatus(notifier notifiers.Notifier, healthStatus storage.IntegrationHealth_Status, errMessage string) {
	panic("unimplemented")
}

// Used for testing.
func (p *processorImpl) processAlertSync(ctx context.Context, alert *storage.Alert) {
	alertNotifiers := set.NewStringSet(alert.GetPolicy().GetNotifiers()...)
	p.ns.ForEach(ctx, func(ctx context.Context, notifier notifiers.Notifier, failures pkgNotifier.AlertSet) {
		if alertNotifiers.Contains(notifier.ProtoNotifier().GetId()) {
			err := pkgNotifier.TryToAlert(ctx, notifier, alert)
			if err != nil {
				failures.Add(alert)
			}
		}
	})
}

// New returns a new Processor
func New(ns pkgNotifier.Set) pkgNotifier.Processor {
	return &processorImpl{
		ns: ns,
	}
}
