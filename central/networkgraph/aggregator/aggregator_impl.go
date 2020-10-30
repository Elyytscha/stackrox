package aggregator

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/networkgraph"
	"github.com/stackrox/rox/pkg/networkgraph/tree"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/utils"
)

const (
	unnamedExtSrcPrefix = "unnamed external source"
)

var (
	log = logging.LoggerForModule()
)

type aggregateDefaultToCustomExtSrcsImpl struct {
	networkTree *tree.NetworkTreeWrapper
}

// Aggregate aggregates multiple external network connections with same external endpoint,
// as determined by name, into a single connection.
func (a *aggregateDefaultToCustomExtSrcsImpl) Aggregate(conns []*storage.NetworkFlow) []*storage.NetworkFlow {
	normalizedConns := make(map[networkgraph.NetworkConnIndicator]*storage.NetworkFlow)
	ret := make([]*storage.NetworkFlow, 0, len(conns)/4)
	supernetCache := make(map[string]*storage.NetworkEntityInfo)

	for _, conn := range conns {
		conn = conn.Clone()
		srcEntity, dstEntity := conn.GetProps().GetSrcEntity(), conn.GetProps().GetDstEntity()
		// This is essentially an invalid connection.
		if srcEntity == nil || dstEntity == nil {
			continue
		}

		if networkgraph.IsExternal(srcEntity) && networkgraph.IsExternal(dstEntity) {
			utils.Should(errors.Errorf("network conn %s with all external endpoints is unexcepted", networkgraph.GetNetworkConnIndicator(conn).String()))
		}

		// If both endpoints are not external (including INTERNET), skip processing.
		if !networkgraph.IsExternal(srcEntity) && !networkgraph.IsExternal(dstEntity) {
			ret = append(ret, conn)
			continue
		}

		// Move the connection from default external network to non-default supernet. If none is found, it gets mapped to INTERNET.
		if networkgraph.IsKnownDefaultExternal(srcEntity) {
			conn.Props.SrcEntity = a.getSupernet(srcEntity.GetId(), supernetCache)
		} else if networkgraph.IsKnownDefaultExternal(dstEntity) {
			conn.Props.DstEntity = a.getSupernet(dstEntity.GetId(), supernetCache)
		}

		connID := networkgraph.GetNetworkConnIndicator(conn)
		if storedFlow := normalizedConns[connID]; storedFlow != nil {
			if storedFlow.GetLastSeenTimestamp().Compare(conn.GetLastSeenTimestamp()) < 0 {
				storedFlow.LastSeenTimestamp = conn.GetLastSeenTimestamp()
			}
		} else {
			normalizedConns[connID] = conn
		}
	}

	for _, conn := range normalizedConns {
		ret = append(ret, conn)
	}
	return ret
}

func (a *aggregateDefaultToCustomExtSrcsImpl) getSupernet(id string, cache map[string]*storage.NetworkEntityInfo) *storage.NetworkEntityInfo {
	supernet := cache[id]
	if supernet == nil {
		supernet = a.networkTree.GetMatchingSupernet(id, func(e *storage.NetworkEntityInfo) bool { return !e.GetExternalSource().GetDefault() })
		cache[id] = supernet
	}
	return supernet
}

type aggregateExternalConnByNameImpl struct{}

// Aggregate aggregates multiple external network connections with same external endpoint, as determined by name,
// into a single connection.
func (a *aggregateExternalConnByNameImpl) Aggregate(flows []*storage.NetworkFlow) []*storage.NetworkFlow {
	conns := make(map[networkgraph.NetworkConnIndicator]*storage.NetworkFlow)
	// dupNameExtSrcTracker stores external source name to ID mapping. This tracks whether an external source name is
	// duplicated by multiple external sources. When an external source name is duplicated, we set the ID value to empty string.
	dupNameExtSrcTracker := make(map[string]string)
	// unnamedExtSrcCounter stores occurrence count of unnamed entities by ID. This helps us generate appropriate names
	// for unnamed external entities such as, "unnamed external source #1/2/3...".
	unnamedExtSrcCounter := make(map[string]int)
	ret := make([]*storage.NetworkFlow, 0, len(flows))

	for _, flow := range flows {
		flow = flow.Clone()
		srcEntity, dstEntity := flow.GetProps().GetSrcEntity(), flow.GetProps().GetDstEntity()
		// This is essentially an invalid connection.
		if srcEntity == nil || dstEntity == nil {
			continue
		}

		if networkgraph.IsExternal(srcEntity) && networkgraph.IsExternal(dstEntity) {
			utils.Should(errors.Errorf("network conn %s with all external endpoints is unexcepted", networkgraph.GetNetworkConnIndicator(flow).String()))
		}

		// If both endpoints are not known external sources, skip processing.
		if !networkgraph.IsKnownExternalSrc(srcEntity) && !networkgraph.IsKnownExternalSrc(dstEntity) {
			ret = append(ret, flow)
			continue
		}

		updateDupNameExtSrcTracker(srcEntity, dupNameExtSrcTracker)
		updateDupNameExtSrcTracker(dstEntity, dupNameExtSrcTracker)

		connIndicator := getNormalizedConnIndicator(flow, unnamedExtSrcCounter)
		// If multiple connections collapse into one, use the latest connection's timestamp to correctly indicate the
		// liveliness of the connection.
		if storedFlow := conns[connIndicator]; storedFlow != nil {
			if storedFlow.GetLastSeenTimestamp().Compare(flow.GetLastSeenTimestamp()) < 0 {
				storedFlow.LastSeenTimestamp = flow.GetLastSeenTimestamp()
			}
		} else {
			conns[connIndicator] = flow
		}
	}

	for connIndicator, conn := range conns {
		// Since entity IDs in conn indicator are normalized to respective entity names, hence we can use them as keys.
		if id, ok := dupNameExtSrcTracker[connIndicator.SrcEntity.ID]; ok && id == "" {
			normalizeDupNameExtSrcs(conn.GetProps().GetSrcEntity())
		}

		if id, ok := dupNameExtSrcTracker[connIndicator.DstEntity.ID]; ok && id == "" {
			normalizeDupNameExtSrcs(conn.GetProps().GetDstEntity())
		}

		ret = append(ret, conn)
	}

	return ret
}

// updateDupNameExtSrcTracker updates dupNameExtSrcTracker which tracks whether an external source name is duplicated
// by multiple external sources. When an external source name is duplicated, we set the ID value to empty string.
func updateDupNameExtSrcTracker(entity *storage.NetworkEntityInfo, dupNameExtSrcTracker map[string]string) {
	if !networkgraph.IsKnownExternalSrc(entity) {
		return
	}

	val, ok := dupNameExtSrcTracker[entity.GetExternalSource().GetName()]
	// If the name is already marked as duplicate, nothing to do.
	if ok && val == "" {
		return
	}

	if !ok {
		val = entity.GetId()
	} else if val != entity.GetId() {
		val = ""
	}
	dupNameExtSrcTracker[entity.GetExternalSource().GetName()] = val
}

// getNormalizedConnIndicator returns indicator for network connections where entity IDs are replaced by their name.
func getNormalizedConnIndicator(conn *storage.NetworkFlow, unnamedExtSrcCounter map[string]int) networkgraph.NetworkConnIndicator {
	srcEntity, dstEntity := conn.GetProps().GetSrcEntity(), conn.GetProps().GetDstEntity()
	connIndicator := networkgraph.GetNetworkConnIndicator(conn)

	// Use entity name as ID for known external sources so that many networks with same name are mapped to
	// one indicator given that other connection indicator properties are the same. External entities created via API
	// always have associated name, therefore, following normalization is unexpected.
	if networkgraph.IsKnownExternalSrc(srcEntity) {
		normalizeUnnamedExternalEntities(srcEntity, unnamedExtSrcCounter)
		connIndicator.SrcEntity.ID = srcEntity.GetExternalSource().GetName()
	} else if networkgraph.IsKnownExternalSrc(dstEntity) {
		normalizeUnnamedExternalEntities(dstEntity, unnamedExtSrcCounter)
		connIndicator.DstEntity.ID = dstEntity.GetExternalSource().GetName()
	}

	return connIndicator
}

func normalizeUnnamedExternalEntities(entity *storage.NetworkEntityInfo, unnamedExtSrcCounter map[string]int) bool {
	if !networkgraph.IsKnownExternalSrc(entity) {
		return false
	}

	if entity.GetExternalSource() == nil {
		entity.Desc = &storage.NetworkEntityInfo_ExternalSource_{
			ExternalSource: &storage.NetworkEntityInfo_ExternalSource{},
		}
	}

	if entity.GetExternalSource().GetName() != "" {
		return false
	}

	if _, ok := unnamedExtSrcCounter[entity.GetId()]; !ok {
		unnamedExtSrcCounter[entity.GetId()] = len(unnamedExtSrcCounter) + 1
	}

	entity.GetExternalSource().Name = fmt.Sprintf("%s #%d", unnamedExtSrcPrefix, unnamedExtSrcCounter[entity.GetId()])
	return true
}

// Note: Update storage.NetworkEntityInfo.ExternalSource comment if this function is refactored, if necessary.
func normalizeDupNameExtSrcs(entity *storage.NetworkEntityInfo) {
	if entity.GetExternalSource() == nil || !networkgraph.IsKnownExternalSrc(entity) {
		return
	}

	// In case of error, we skip normalization. External entities created via API always have correct resource ID,
	// hence, the following errors are unexpected.
	decodedID, err := sac.ParseResourceID(entity.GetId())
	if err != nil {
		log.Errorf("failed to normalize external sources: %v", err)
		return
	}

	id, err := sac.NewClusterScopeResourceID(decodedID.ClusterID(), entity.GetExternalSource().GetName())
	if err != nil {
		log.Errorf("failed to normalize external sources: %v", err)
		return
	}

	*entity = storage.NetworkEntityInfo{
		Id:   id.String(),
		Type: storage.NetworkEntityInfo_EXTERNAL_SOURCE,
		Desc: &storage.NetworkEntityInfo_ExternalSource_{
			ExternalSource: &storage.NetworkEntityInfo_ExternalSource{
				Name:    entity.GetExternalSource().GetName(),
				Default: entity.GetExternalSource().GetDefault(),
				// Since many CIDRs are mapped to one endpoint, we clear the CIDR field.
			},
		},
	}
}
