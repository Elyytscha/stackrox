package alertmanager

import (
	"context"

	"github.com/gogo/protobuf/proto"
	ptypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	alertDataStore "github.com/stackrox/rox/central/alert/datastore"
	"github.com/stackrox/rox/central/detection/runtime"
	notifierProcessor "github.com/stackrox/rox/central/notifier/processor"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/booleanpolicy/violationmessages/printer"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/protoutils"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
)

const maxProcessViolationsPerAlert = 40

var (
	log = logging.LoggerForModule()
)

type alertManagerImpl struct {
	notifier        notifierProcessor.Processor
	alerts          alertDataStore.DataStore
	runtimeDetector runtime.Detector
}

func getDeploymentIDsFromAlerts(alertSlices ...[]*storage.Alert) set.StringSet {
	s := set.NewStringSet()
	for _, slice := range alertSlices {
		for _, alert := range slice {
			s.Add(alert.GetDeployment().GetId())
		}
	}
	return s
}

func (d *alertManagerImpl) AlertAndNotify(ctx context.Context, currentAlerts []*storage.Alert, oldAlertFilters ...AlertFilterOption) (set.StringSet, error) {
	// Merge the old and the new alerts.
	newAlerts, updatedAlerts, staleAlerts, err := d.mergeManyAlerts(ctx, currentAlerts, oldAlertFilters...)
	if err != nil {
		return nil, err
	}

	modifiedDeployments := getDeploymentIDsFromAlerts(newAlerts, updatedAlerts, staleAlerts)

	// Mark any old alerts no longer generated as stale, and insert new alerts.
	err = d.notifyAndUpdateBatch(ctx, newAlerts)
	if err != nil {
		return nil, err
	}
	err = d.updateBatch(ctx, updatedAlerts)
	if err != nil {
		return nil, err
	}
	err = d.markAlertsStale(ctx, staleAlerts)
	if err != nil {
		return nil, err
	}
	return modifiedDeployments, nil
}

// UpdateBatch updates all of the alerts in the datastore.
func (d *alertManagerImpl) updateBatch(ctx context.Context, alertsToMark []*storage.Alert) error {
	errList := errorhelpers.NewErrorList("Error updating alerts: ")
	for _, existingAlert := range alertsToMark {
		errList.AddError(d.alerts.UpsertAlert(ctx, existingAlert))
	}
	return errList.ToError()
}

// MarkAlertsStale marks all of the input alerts stale in the input datastore.
func (d *alertManagerImpl) markAlertsStale(ctx context.Context, alertsToMark []*storage.Alert) error {
	errList := errorhelpers.NewErrorList("Error marking alerts as stale: ")
	for _, existingAlert := range alertsToMark {
		err := d.alerts.MarkAlertStale(ctx, existingAlert.GetId())
		if err == nil {
			// run notifier for all the resolved alerts
			existingAlert.State = storage.ViolationState_RESOLVED
			d.notifier.ProcessAlert(ctx, existingAlert)
		}
		errList.AddError(err)
	}
	return errList.ToError()
}

// NotifyAndUpdateBatch runs the notifier on the input alerts then stores them.
func (d *alertManagerImpl) notifyAndUpdateBatch(ctx context.Context, alertsToMark []*storage.Alert) error {
	for _, existingAlert := range alertsToMark {
		d.notifier.ProcessAlert(ctx, existingAlert)
	}
	return d.updateBatch(ctx, alertsToMark)
}

// It is the caller's responsibility to not call this with an empty slice,
// else this function will panic.
func lastTimestamp(processes []*storage.ProcessIndicator) *ptypes.Timestamp {
	return processes[len(processes)-1].GetSignal().GetTime()
}

// Some processes in the old alert might have been deleted from the process store because of our pruning,
// which means they only exist in the old alert, and will not be in the new generated alert.
// We don't want to lose them, though, so we keep all the processes from the old alert, and add ones from the new, if any.
// Note that the old alert _was_ active which means that all the processes in it are guaranteed to violate the policy.
func mergeProcessesFromOldIntoNew(old, newAlert *storage.Alert) (newAlertHasNewProcesses bool) {
	oldProcessViolation := old.GetProcessViolation()

	if len(oldProcessViolation.GetProcesses()) == 0 {
		log.Errorf("UNEXPECTED: found no old violation with processes for runtime alert %s", proto.MarshalTextString(old))
		newAlertHasNewProcesses = true
		return
	}

	if len(newAlert.GetProcessViolation().GetProcesses()) == 0 {
		log.Errorf("UNEXPECTED: found no new violation with processes for runtime alert %s", proto.MarshalTextString(newAlert))
		return
	}

	if len(oldProcessViolation.GetProcesses()) >= maxProcessViolationsPerAlert {
		return
	}

	newProcessesSlice := oldProcessViolation.GetProcesses()
	// De-dupe processes using timestamps.
	timestamp := lastTimestamp(oldProcessViolation.GetProcesses())
	for _, process := range newAlert.GetProcessViolation().GetProcesses() {
		if process.GetSignal().GetTime().Compare(timestamp) > 0 {
			newAlertHasNewProcesses = true
			newProcessesSlice = append(newProcessesSlice, process)
		}
	}
	// If there are no new processes, we'll just use the old alert.
	if !newAlertHasNewProcesses {
		return
	}
	if len(newProcessesSlice) > maxProcessViolationsPerAlert {
		newProcessesSlice = newProcessesSlice[:maxProcessViolationsPerAlert]
	}
	newAlert.ProcessViolation.Processes = newProcessesSlice
	printer.UpdateProcessAlertViolationMessage(newAlert.ProcessViolation)
	return
}

// MergeAlerts merges two alerts.
func mergeAlerts(old, newAlert *storage.Alert) *storage.Alert {
	if old.GetLifecycleStage() == storage.LifecycleStage_RUNTIME && newAlert.GetLifecycleStage() == storage.LifecycleStage_RUNTIME {
		newAlertHasNewProcesses := mergeProcessesFromOldIntoNew(old, newAlert)
		// This ensures that we don't keep updating an old runtime alert, so that we have idempotent checks.
		if !newAlertHasNewProcesses {
			return old
		}
	}

	newAlert.Id = old.GetId()
	// Updated deploy-time alerts continue to have the same enforcement action.
	if newAlert.GetLifecycleStage() == storage.LifecycleStage_DEPLOY && old.GetLifecycleStage() == storage.LifecycleStage_DEPLOY {
		newAlert.Enforcement = old.GetEnforcement()
		// Don't keep updating the timestamp of the violation _unless_ the violations are actually different.
		if protoutils.EqualStorageAlert_ViolationSlices(newAlert.GetViolations(), old.GetViolations()) {
			newAlert.Time = old.GetTime()
		}
	}

	newAlert.FirstOccurred = old.GetFirstOccurred()
	newAlert.Tags = old.GetTags()
	return newAlert
}

// MergeManyAlerts merges two alerts.
func (d *alertManagerImpl) mergeManyAlerts(ctx context.Context, presentAlerts []*storage.Alert, oldAlertFilters ...AlertFilterOption) (newAlerts, updatedAlerts, staleAlerts []*storage.Alert, err error) {
	qb := search.NewQueryBuilder().AddStrings(search.ViolationState, storage.ViolationState_ACTIVE.String())
	for _, filter := range oldAlertFilters {
		filter.apply(qb)
	}
	previousAlerts, err := d.alerts.SearchRawAlerts(ctx, qb.ProtoQuery())
	if err != nil {
		err = errors.Wrapf(err, "couldn't load previous alerts (query was %s)", qb.Query())
		return
	}

	// Merge any alerts that have new and old alerts.
	for _, alert := range presentAlerts {
		if matchingOld := findAlert(alert, previousAlerts); matchingOld != nil {
			mergedAlert := mergeAlerts(matchingOld, alert)
			if mergedAlert != matchingOld && !proto.Equal(mergedAlert, matchingOld) {
				updatedAlerts = append(updatedAlerts, mergedAlert)
			}
			continue
		}

		alert.FirstOccurred = ptypes.TimestampNow()
		newAlerts = append(newAlerts, alert)
	}

	// Find any old alerts no longer being produced.
	for _, alert := range previousAlerts {
		if d.shouldMarkAlertStale(alert, presentAlerts, oldAlertFilters...) {
			staleAlerts = append(staleAlerts, alert)
		}

		if alert.GetLifecycleStage() == storage.LifecycleStage_RUNTIME && d.inactiveDeploymentAlert(alert) {
			if !alert.Deployment.Inactive {
				alert.Deployment.Inactive = true
				updatedAlerts = append(updatedAlerts, alert)
			}
		}
	}
	return
}

func (d *alertManagerImpl) shouldMarkAlertStale(alert *storage.Alert, presentAlerts []*storage.Alert, oldAlertFilters ...AlertFilterOption) bool {
	// If the alert is still being produced, don't mark it stale.
	if matchingNew := findAlert(alert, presentAlerts); matchingNew != nil {
		return false
	}

	// Only runtime alerts should not be marked stale when they are no longer produced.
	// (Deploy time alerts should disappear along with deployments, for example.)
	if alert.GetLifecycleStage() != storage.LifecycleStage_RUNTIME {
		return true
	}

	if !d.runtimeDetector.PolicySet().Exists(alert.GetPolicy().GetId()) {
		return true
	}

	// We only want to mark runtime alerts as stale if a policy update causes them to no longer be produced.
	// To determine if this is a policy update, we check if there is a filter on policy ids here.
	specifiedPolicyIDs := set.NewStringSet()
	for _, filter := range oldAlertFilters {
		if filterSpecified := filter.specifiedPolicyID(); filterSpecified != "" {
			specifiedPolicyIDs.Add(filterSpecified)
		}
	}
	if specifiedPolicyIDs.Cardinality() == 0 {
		return false
	}

	// Some other policies were updated, we don't want to mark this alert stale in response.
	if !specifiedPolicyIDs.Contains(alert.GetPolicy().GetId()) {
		return false
	}

	// If the deployment is excluded from the scope of the policy now, we should mark the alert stale, otherwise we will keep it around.
	return d.runtimeDetector.DeploymentWhitelistedForPolicy(alert.GetDeployment().GetId(), alert.GetPolicy().GetId())
}

func (d *alertManagerImpl) inactiveDeploymentAlert(alert *storage.Alert) bool {
	return d.runtimeDetector.DeploymentInactive(alert.GetDeployment().GetId())
}

func findAlert(toFind *storage.Alert, alerts []*storage.Alert) *storage.Alert {
	for _, alert := range alerts {
		if alertsAreForSamePolicyAndDeployment(alert, toFind) {
			return alert
		}
	}
	return nil
}

func alertsAreForSamePolicyAndDeployment(a1, a2 *storage.Alert) bool {
	return a1.GetPolicy().GetId() == a2.GetPolicy().GetId() && a1.GetDeployment().GetId() == a2.GetDeployment().GetId()
}
