package aggregator

import (
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/networkgraph"
	"github.com/stackrox/rox/pkg/networkgraph/externalsrcs"
	"github.com/stackrox/rox/pkg/networkgraph/tree"
	"github.com/stackrox/rox/pkg/utils"
)

func mapToSupernet(networkTree tree.ReadOnlyNetworkTree,
	supernetCache map[string]*storage.NetworkEntityInfo,
	supernetPred func(e *storage.NetworkEntityInfo) bool, entityPtrs ...**storage.NetworkEntityInfo) {
	for _, entityPtr := range entityPtrs {
		entity := *entityPtr
		if !networkgraph.IsKnownExternalSrc(entity) {
			continue
		}

		cidr, err := externalsrcs.NetworkFromID(entity.GetId())
		if err != nil {
			utils.Should(errors.Wrapf(err, "getting CIDR from external source ID %s", entity.GetId()))
			*entityPtr = networkgraph.InternetEntity().ToProto()
			continue
		}
		*entityPtr = getSupernet(networkTree, supernetCache, cidr, supernetPred)
	}
}

func getSupernet(networkTree tree.ReadOnlyNetworkTree,
	cache map[string]*storage.NetworkEntityInfo,
	cidr string,
	supernetPred func(e *storage.NetworkEntityInfo) bool) *storage.NetworkEntityInfo {
	supernet := cache[cidr]
	if supernet == nil {
		supernet = networkTree.GetMatchingSupernetForCIDR(cidr, supernetPred)
		cache[cidr] = supernet
	}
	return supernet
}
