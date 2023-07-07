package manager

import (
	"fmt"
	"strings"

	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/alert"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/env"
	eventPkg "github.com/stackrox/rox/pkg/sensor/event"
	"github.com/stackrox/rox/pkg/sensor/hash"
	"github.com/stackrox/rox/pkg/stringutils"
	"github.com/stackrox/rox/pkg/sync"
)

// Deduper is an interface to deduping logic used to determine whether an event should be processed
type Deduper interface {
	// GetSuccessfulHashes returns a map of key to hashes that were successfully processed by Central
	// and thus can be persisted in the database as being processed
	GetSuccessfulHashes() map[string]uint64
	// ShouldProcess takes a message and determines if it should be processed
	ShouldProcess(msg *central.MsgFromSensor) bool
	// MarkSuccessful marks the message if necessary as being successfully processed, so it can be committed to the database
	// It will promote the message from the received map to the successfully processed map
	MarkSuccessful(msg *central.MsgFromSensor)
	// StartSync is called once a new Sensor connection is initialized
	StartSync()
	// ProcessSync processes the Sensor sync message and reconciles the successfully processed and received maps
	ProcessSync()
}

var (
	maxHashes = env.MaxEventHashSize.IntegerSetting()

	nodeInventoryResourceKey = eventPkg.GetEventTypeWithoutPrefix((*central.SensorEvent_NodeInventory)(nil))
	alertResourceKey         = eventPkg.GetEventTypeWithoutPrefix((*central.SensorEvent_AlertResults)(nil))
	deploymentResourceKey    = eventPkg.GetEventTypeWithoutPrefix((*central.SensorEvent_Deployment)(nil))
)

// NewDeduper creates a new deduper from the passed existing hashes
func NewDeduper(existingHashes map[string]uint64) Deduper {
	existingEntries := make(map[string]*entry)

	// In 4.1, node inventory hashes would be added to the stored hashes, but they will not be
	// removed by node removal so we should remove them upon start
	for k := range existingHashes {
		if strings.HasPrefix(k, nodeInventoryResourceKey) {
			delete(existingHashes, k)
		}
	}

	if len(existingHashes) < maxHashes {
		for k, v := range existingHashes {
			existingEntries[k] = &entry{
				val: v,
			}
		}
	}

	return &deduperImpl{
		received:              make(map[string]*entry),
		successfullyProcessed: existingEntries,
		hasher:                hash.NewHasher(),
	}
}

type entry struct {
	val       uint64
	processed bool
}

type deduperImpl struct {
	hashLock sync.RWMutex
	// received map contains messages that have been received but not successfully processed
	received map[string]*entry
	// successfully processed map contains hashes of objects that have been successfully processed
	successfullyProcessed map[string]*entry

	hasher *hash.Hasher
}

func isAlwaysTrackedObject(event *central.SensorEvent) bool {
	// This switch specifies the events that should always be tracked unconditionally
	switch event.Resource.(type) {
	case *central.SensorEvent_NetworkPolicy:
	case *central.SensorEvent_Deployment:
	case *central.SensorEvent_Pod:
	case *central.SensorEvent_Namespace:
	case *central.SensorEvent_Node:
	case *central.SensorEvent_ServiceAccount:
	case *central.SensorEvent_Role:
	case *central.SensorEvent_Binding:
	case *central.SensorEvent_ComplianceOperatorResult:
	case *central.SensorEvent_ComplianceOperatorProfile:
	case *central.SensorEvent_ComplianceOperatorScan:
	case *central.SensorEvent_ComplianceOperatorRule:
	case *central.SensorEvent_ComplianceOperatorScanSettingBinding:
	default:
		// If the object is not in the above list, return false
		return false
	}
	return true
}

// skipDedupe signifies that a message from Sensor cannot be deduped and won't be stored
func skipDedupe(msg *central.MsgFromSensor) bool {
	eventMsg, ok := msg.Msg.(*central.MsgFromSensor_Event)
	if !ok {
		return true
	}
	event := eventMsg.Event

	// If the object should always be tracked, then we should never skip deduping
	if isAlwaysTrackedObject(event) {
		return false
	}
	if results := event.GetAlertResults(); results != nil {
		// Do not dedupe deploy time alerts that are unresolved and are not attempted
		if alert.IsDeployTimeAlertResult(event.GetAlertResults()) &&
			!alert.IsAlertResultResolved(results) &&
			!alert.AnyAttemptedAlert(results.GetAlerts()...) {
			return false
		}
	}
	// If we did not match an expected case, then skip deduping which will not have any correctness implications
	// just potential performance implications
	return true
}

func (d *deduperImpl) shouldReprocess(hashKey string, hash uint64) bool {
	d.hashLock.RLock()
	defer d.hashLock.RUnlock()

	prevValue, ok := d.getValueNoLock(hashKey)
	if !ok {
		// This implies that a REMOVE event has been processed before this event
		// Note: we may want to handle alerts specifically because we should insert them as already resolved for completeness
		return false
	}
	// This implies that no new event was processed after the initial processing of the current message
	return prevValue == hash
}

func getIDFromKey(key string) string {
	return stringutils.GetAfter(key, ":")
}

func buildKey(typ, id string) string {
	return fmt.Sprintf("%s:%s", typ, id)
}

func getKey(msg *central.MsgFromSensor) string {
	event := msg.GetEvent()
	return buildKey(eventPkg.GetEventTypeWithoutPrefix(event.GetResource()), event.GetId())
}

// StartSync is called when Sensor starts a new connection
func (d *deduperImpl) StartSync() {
	d.hashLock.Lock()
	defer d.hashLock.Unlock()

	for _, v := range d.received {
		v.processed = false
	}

	// Mark all hashes as unseen when a new sensor connection is created
	for _, v := range d.successfullyProcessed {
		v.processed = false
	}
}

// MarkSuccessful marks a message as successfully processed
func (d *deduperImpl) MarkSuccessful(msg *central.MsgFromSensor) {
	// If the object isn't eligible for deduping then do not mark it as being successfully processed
	// because it does not exist in the received map
	if skipDedupe(msg) {
		return
	}
	key := getKey(msg)

	d.hashLock.Lock()
	defer d.hashLock.Unlock()
	// If we are removing, then we do not need to mark it as successful as there is nothing more
	// to potentially dedupe. We do need to remove it from successfully processed as ShouldProcess is
	// evaluated as objects come into the queue and an object may have been successfully processed after
	if msg.GetEvent().GetAction() == central.ResourceAction_REMOVE_RESOURCE {
		delete(d.successfullyProcessed, key)
		return
	}

	val, ok := d.received[key]
	// Only remove from this map if the hash matches as received could contain a more recent event
	if ok && val.val == msg.GetEvent().GetSensorHash() {
		delete(d.received, key)
	}
	d.successfullyProcessed[key] = &entry{
		val:       msg.GetEvent().GetSensorHash(),
		processed: true,
	}
}

func (d *deduperImpl) getValueNoLock(key string) (uint64, bool) {
	if prevValue, ok := d.received[key]; ok {
		return prevValue.val, ok
	}
	prevValue, ok := d.successfullyProcessed[key]
	if !ok {
		return 0, false
	}
	return prevValue.val, true
}

// ProcessSync is triggered by the sync message sent from Sensor
func (d *deduperImpl) ProcessSync() {
	// Reconcile successfully processed map with received map. Any keys that exist in successfully processed
	// but do not exist in received, can be dropped from successfully processed
	d.hashLock.Lock()
	defer d.hashLock.Unlock()

	for k, v := range d.successfullyProcessed {
		// Ignore alerts in the first pass because they are not reconciled at the same time
		if strings.HasPrefix(k, alertResourceKey) {
			continue
		}
		if !v.processed {
			if val, ok := d.received[k]; ok && val.processed {
				continue
			}
			delete(d.successfullyProcessed, k)
			// If a deployment is being removed due to reconciliation, then we will need to remove the alerts too
			if strings.HasPrefix(k, deploymentResourceKey) {
				alertKey := buildKey(alertResourceKey, getIDFromKey(k))
				delete(d.successfullyProcessed, alertKey)
			}
		}
	}
}

// ShouldProcess determines if a message should be processed or if it should be deduped and dropped
func (d *deduperImpl) ShouldProcess(msg *central.MsgFromSensor) bool {
	if skipDedupe(msg) {
		return true
	}
	event := msg.GetEvent()
	key := getKey(msg)
	switch event.GetAction() {
	case central.ResourceAction_REMOVE_RESOURCE:
		d.hashLock.Lock()
		defer d.hashLock.Unlock()

		delete(d.received, key)
		delete(d.successfullyProcessed, key)
		return true
	case central.ResourceAction_SYNC_RESOURCE:
		// check if element is in successfully processed and mark as processed for syncs so that these are not reconciled away
		concurrency.WithLock(&d.hashLock, func() {
			if val, ok := d.successfullyProcessed[key]; ok {
				val.processed = true
			}
		})
	}
	// Backwards compatibility with a previous Sensor
	if event.GetSensorHashOneof() == nil {
		// Compute the sensor hash
		hashValue, ok := d.hasher.HashEvent(msg.GetEvent())
		if !ok {
			return true
		}
		event.SensorHashOneof = &central.SensorEvent_SensorHash{
			SensorHash: hashValue,
		}
	}
	// In the reprocessing case, the above will never evaluate to not nil, but it makes testing easier
	if msg.GetProcessingAttempt() > 0 {
		return d.shouldReprocess(key, event.GetSensorHash())
	}

	d.hashLock.Lock()
	defer d.hashLock.Unlock()

	prevValue, ok := d.getValueNoLock(key)
	if ok && prevValue == event.GetSensorHash() {
		return false
	}
	d.received[key] = &entry{
		val:       event.GetSensorHash(),
		processed: true,
	}
	return true
}

// GetSuccessfulHashes returns a copy of the successfullyProcessed map
func (d *deduperImpl) GetSuccessfulHashes() map[string]uint64 {
	d.hashLock.RLock()
	defer d.hashLock.RUnlock()

	copied := make(map[string]uint64, len(d.successfullyProcessed))
	for k, v := range d.successfullyProcessed {
		copied[k] = v.val
	}

	// Do not persist copied map if it has more than the max number of hashes. This could occur in two cases
	// 1. The environment is much larger than anticipated (default max hashes is 1 million).
	// 2. There is an event that does not have a proper lifecycle and it is causing unbounded growth.
	if len(copied) > maxHashes {
		return make(map[string]uint64)
	}
	return copied
}
