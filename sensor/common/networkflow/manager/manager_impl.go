package manager

import (
	"errors"
	"time"

	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/internalapi/sensor"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/centralsensor"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/net"
	"github.com/stackrox/rox/pkg/netutil"
	"github.com/stackrox/rox/pkg/networkgraph"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/timestamp"
	"github.com/stackrox/rox/sensor/common/clusterentities"
	"github.com/stackrox/rox/sensor/common/clusterid"
	"github.com/stackrox/rox/sensor/common/metrics"
)

const (
	// Wait at least this long before determining that an unresolvable IP is "outside of the cluster".
	clusterEntityResolutionWaitPeriod = 10 * time.Second

	connectionDeletionGracePeriod = 5 * time.Minute
)

type hostConnections struct {
	connections        map[connection]*connStatus
	lastKnownTimestamp timestamp.MicroTS

	// connectionsSequenceID is the sequence ID of the current active connections state
	connectionsSequenceID int64
	// currentSequenceID is the sequence ID of the most recent `Register` call
	currentSequenceID int64

	pendingDeletion *time.Timer

	mutex sync.Mutex
}

type connStatus struct {
	firstSeen timestamp.MicroTS
	lastSeen  timestamp.MicroTS
	used      bool
}

type networkConnIndicator struct {
	srcEntity networkgraph.Entity
	dstEntity networkgraph.Entity
	dstPort   uint16
	protocol  storage.L4Protocol
}

func (i networkConnIndicator) toProto(ts timestamp.MicroTS) *storage.NetworkFlow {
	proto := &storage.NetworkFlow{
		Props: &storage.NetworkFlowProperties{
			SrcEntity:  i.srcEntity.ToProto(),
			DstEntity:  i.dstEntity.ToProto(),
			DstPort:    uint32(i.dstPort),
			L4Protocol: i.protocol,
		},
	}

	if ts != timestamp.InfiniteFuture {
		proto.LastSeenTimestamp = ts.GogoProtobuf()
	}
	return proto
}

// connection is an instance of a connection as reported by collector
type connection struct {
	local       net.IPPortPair
	remote      net.NumericEndpoint
	containerID string
	incoming    bool
}

type networkFlowManager struct {
	connectionsByHost      map[string]*hostConnections
	connectionsByHostMutex sync.Mutex

	clusterEntities *clusterentities.Store

	enrichedLastSentState map[networkConnIndicator]timestamp.MicroTS

	done        concurrency.Signal
	flowUpdates chan *central.MsgFromSensor
}

func (m *networkFlowManager) ProcessMessage(msg *central.MsgToSensor) error {
	return nil
}

func (m *networkFlowManager) Start() error {
	go m.enrichConnections()
	return nil
}

func (m *networkFlowManager) Stop(_ error) {
	m.done.Signal()
}

func (m *networkFlowManager) Capabilities() []centralsensor.SensorCapability {
	return nil
}

func (m *networkFlowManager) ResponsesC() <-chan *central.MsgFromSensor {
	return m.flowUpdates
}

func (m *networkFlowManager) enrichConnections() {
	ticker := time.NewTicker(time.Second * 30)

	for {
		select {
		case <-m.done.WaitC():
			return
		case <-ticker.C:
			m.enrichAndSend()
		}
	}
}

func (m *networkFlowManager) enrichAndSend() {
	current := m.currentEnrichedConns()

	protoToSend := computeUpdateMessage(current, m.enrichedLastSentState)
	m.enrichedLastSentState = current

	if protoToSend == nil {
		return
	}

	metrics.IncrementTotalNetworkFlowsSentCounter(clusterid.Get(), len(protoToSend.Updated))
	log.Debugf("Flow update : %v", protoToSend)
	select {
	case <-m.done.Done():
		return
	case m.flowUpdates <- &central.MsgFromSensor{
		Msg: &central.MsgFromSensor_NetworkFlowUpdate{
			NetworkFlowUpdate: protoToSend,
		},
	}:
		return
	}
}

func (m *networkFlowManager) enrichConnection(conn *connection, status *connStatus, enrichedConnections map[networkConnIndicator]timestamp.MicroTS) {
	isFresh := timestamp.Now().ElapsedSince(status.firstSeen) < clusterEntityResolutionWaitPeriod
	if !isFresh {
		status.used = true
	}

	container, ok := m.clusterEntities.LookupByContainerID(conn.containerID)
	if !ok {
		log.Debugf("Unable to fetch deployment information for container %s: no deployment found", conn.containerID)
		return
	}

	lookupResults := m.clusterEntities.LookupByEndpoint(conn.remote)
	if len(lookupResults) == 0 {
		if isFresh {
			return
		}

		var port uint16
		if conn.incoming {
			port = conn.local.Port
		} else {
			port = conn.remote.IPAndPort.Port
		}

		// Fake a lookup result with an empty deployment ID.
		lookupResults = []clusterentities.LookupResult{
			{
				Entity: networkgraph.Entity{
					Type: storage.NetworkEntityInfo_INTERNET,
				},
				ContainerPorts: []uint16{port},
			},
		}
	} else {
		status.used = true
		if conn.incoming {
			// Only report incoming connections from outside of the cluster. These are already taken care of by the
			// corresponding outgoing connection from the other end.
			return
		}
	}

	for _, lookupResult := range lookupResults {
		for _, port := range lookupResult.ContainerPorts {
			indicator := networkConnIndicator{
				dstPort:  port,
				protocol: conn.remote.L4Proto.ToProtobuf(),
			}

			if conn.incoming {
				indicator.srcEntity = lookupResult.Entity
				indicator.dstEntity = networkgraph.EntityForDeployment(container.DeploymentID)
			} else {
				indicator.srcEntity = networkgraph.EntityForDeployment(container.DeploymentID)
				indicator.dstEntity = lookupResult.Entity
			}

			// Multiple connections from a collector can result in a single enriched connection
			// hence update the timestamp only if we have a more recent connection than the one we have already enriched.
			if oldTS, found := enrichedConnections[indicator]; !found || oldTS < status.lastSeen {
				enrichedConnections[indicator] = status.lastSeen
			}
		}
	}
}

func (m *networkFlowManager) enrichHostConnections(hostConns *hostConnections, enrichedConnections map[networkConnIndicator]timestamp.MicroTS) {
	hostConns.mutex.Lock()
	defer hostConns.mutex.Unlock()

	for conn, status := range hostConns.connections {
		m.enrichConnection(&conn, status, enrichedConnections)
		if status.used && status.lastSeen != timestamp.InfiniteFuture {
			// connections that are no longer active and have already been used can be deleted.
			delete(hostConns.connections, conn)
		}
	}
}

func (m *networkFlowManager) currentEnrichedConns() map[networkConnIndicator]timestamp.MicroTS {
	allHostConns := m.getAllHostConnections()

	enrichedConnections := make(map[networkConnIndicator]timestamp.MicroTS)
	for _, hostConns := range allHostConns {
		m.enrichHostConnections(hostConns, enrichedConnections)
	}

	return enrichedConnections
}

func computeUpdateMessage(current map[networkConnIndicator]timestamp.MicroTS, previous map[networkConnIndicator]timestamp.MicroTS) *central.NetworkFlowUpdate {
	var updates []*storage.NetworkFlow

	for conn, currTS := range current {
		prevTS, ok := previous[conn]
		if !ok || currTS > prevTS {
			updates = append(updates, conn.toProto(currTS))
		}
	}

	for conn, prevTS := range previous {
		if _, ok := current[conn]; !ok {
			updates = append(updates, conn.toProto(prevTS))
		}
	}

	if len(updates) == 0 {
		return nil
	}

	return &central.NetworkFlowUpdate{
		Updated: updates,
		Time:    timestamp.Now().GogoProtobuf(),
	}
}

func (m *networkFlowManager) getAllHostConnections() []*hostConnections {
	// Get a snapshot of all *hostConnections. This allows us to lock the individual mutexes without having to hold
	// two locks simultaneously.
	m.connectionsByHostMutex.Lock()
	defer m.connectionsByHostMutex.Unlock()

	allHostConns := make([]*hostConnections, 0, len(m.connectionsByHost))
	for _, hostConns := range m.connectionsByHost {
		allHostConns = append(allHostConns, hostConns)
	}

	return allHostConns
}

func (m *networkFlowManager) RegisterCollector(hostname string) (HostNetworkInfo, int64) {
	m.connectionsByHostMutex.Lock()
	defer m.connectionsByHostMutex.Unlock()

	conns := m.connectionsByHost[hostname]

	if conns == nil {
		conns = &hostConnections{
			connections: make(map[connection]*connStatus),
		}
		m.connectionsByHost[hostname] = conns
	}

	conns.mutex.Lock()
	defer conns.mutex.Unlock()

	if conns.pendingDeletion != nil {
		// Note that we don't need to check the return value, since `deleteHostConnections` needs to acquire
		// m.connectionsByHostMutex. It can therefore only proceed once this function returns, in which case it will be
		// a no-op due to `pendingDeletion` being `nil`.
		conns.pendingDeletion.Stop()
		conns.pendingDeletion = nil
	}

	conns.currentSequenceID++

	return conns, conns.currentSequenceID
}

func (m *networkFlowManager) deleteHostConnections(hostname string) {
	m.connectionsByHostMutex.Lock()
	defer m.connectionsByHostMutex.Unlock()

	conns := m.connectionsByHost[hostname]
	if conns == nil {
		return
	}

	conns.mutex.Lock()
	defer conns.mutex.Unlock()

	if conns.pendingDeletion == nil {
		return
	}

	delete(m.connectionsByHost, hostname)
}

func (m *networkFlowManager) UnregisterCollector(hostname string, sequenceID int64) {
	m.connectionsByHostMutex.Lock()
	defer m.connectionsByHostMutex.Unlock()

	conns := m.connectionsByHost[hostname]
	if conns == nil {
		return
	}

	conns.mutex.Lock()
	defer conns.mutex.Unlock()
	if conns.currentSequenceID != sequenceID {
		// Skip deletion if there has been a more recent Register call than the corresponding Unregister call
		return
	}

	if conns.pendingDeletion != nil {
		// Cancel any pending deletions there might be. See `RegisterCollector` on why we do not need to check for the
		// return value of Stop.
		conns.pendingDeletion.Stop()
	}
	conns.pendingDeletion = time.AfterFunc(connectionDeletionGracePeriod, func() {
		m.deleteHostConnections(hostname)
	})
}

func (h *hostConnections) Process(networkInfo *sensor.NetworkConnectionInfo, nowTimestamp timestamp.MicroTS, sequenceID int64) error {
	updatedConnections := getUpdatedConnections(networkInfo)

	collectorTS := timestamp.FromProtobuf(networkInfo.GetTime())
	tsOffset := nowTimestamp - collectorTS

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if sequenceID != h.currentSequenceID {
		return errors.New("replaced by newer connection")
	} else if sequenceID != h.connectionsSequenceID {
		// This is the first message of the new connection.
		for _, status := range h.connections {
			// Mark all connections as closed this is the first update
			// after a connection went down and came back up again.
			status.lastSeen = h.lastKnownTimestamp
		}
		h.connectionsSequenceID = sequenceID
	}

	for c, t := range updatedConnections {
		// timestamp = zero implies the connection is newly added. Add new connections, update existing ones to mark them closed
		if t != timestamp.InfiniteFuture { // adjust timestamp if not zero.
			t += tsOffset
		}
		status := h.connections[c]
		if status == nil {
			status = &connStatus{
				firstSeen: timestamp.Now(),
			}
			if t < status.firstSeen {
				status.firstSeen = t
			}
			h.connections[c] = status
		}
		status.lastSeen = t
	}

	h.lastKnownTimestamp = nowTimestamp

	return nil
}

func getIPAndPort(address *sensor.NetworkAddress) net.IPPortPair {
	return net.IPPortPair{
		Address: net.IPFromBytes(address.GetAddressData()),
		Port:    uint16(address.GetPort()),
	}
}

func getUpdatedConnections(networkInfo *sensor.NetworkConnectionInfo) map[connection]timestamp.MicroTS {
	updatedConnections := make(map[connection]timestamp.MicroTS)

	for _, conn := range networkInfo.GetUpdatedConnections() {
		var incoming bool
		switch conn.Role {
		case sensor.ClientServerRole_ROLE_SERVER:
			incoming = true
		case sensor.ClientServerRole_ROLE_CLIENT:
			incoming = false
		default:
			continue
		}

		remote := net.NumericEndpoint{
			IPAndPort: getIPAndPort(conn.GetRemoteAddress()),
			L4Proto:   net.L4ProtoFromProtobuf(conn.GetProtocol()),
		}
		local := getIPAndPort(conn.GetLocalAddress())

		// Special handling for UDP ports - role reported by collector may be unreliable, so look at which port is more
		// likely to be ephemeral.
		if remote.L4Proto == net.UDP {
			incoming = netutil.IsEphemeralPort(remote.IPAndPort.Port) > netutil.IsEphemeralPort(local.Port)
		}

		c := connection{
			local:       local,
			remote:      remote,
			containerID: conn.GetContainerId(),
			incoming:    incoming,
		}

		// timestamp will be set to close timestamp for closed connections, and zero for newly added connection.
		ts := timestamp.FromProtobuf(conn.CloseTimestamp)
		if ts == 0 {
			ts = timestamp.InfiniteFuture
		}
		updatedConnections[c] = ts
	}

	return updatedConnections
}
