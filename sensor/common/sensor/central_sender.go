package sensor

import (
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/listeners"
	"github.com/stackrox/rox/sensor/common"
	"github.com/stackrox/rox/sensor/common/clusterstatus"
	"github.com/stackrox/rox/sensor/common/compliance"
	"github.com/stackrox/rox/sensor/common/config"
	networkConnManager "github.com/stackrox/rox/sensor/common/networkflow/manager"
	"github.com/stackrox/rox/sensor/common/networkpolicies"
	"github.com/stackrox/rox/sensor/common/signal"
)

// CentralSender handles sending from sensor to central.
type CentralSender interface {
	Start(stream central.SensorService_CommunicateClient, onStops ...func(error))

	Stop(err error)
	Stopped() concurrency.ReadOnlyErrorSignal
}

// NewCentralSender returns a new instance of a CentralSender.
func NewCentralSender(listener listeners.Listener,
	signalService signal.Service,
	networkConnManager networkConnManager.Manager,
	scrapeCommandHandler compliance.CommandHandler,
	networkPoliciesCommandHandler networkpolicies.CommandHandler,
	clusterStatusUpdater clusterstatus.Updater,
	configCommandHandler config.Handler,
	otherComponents ...common.SensorComponent) CentralSender {
	return &centralSenderImpl{
		listener:                      listener,
		signalService:                 signalService,
		networkConnManager:            networkConnManager,
		scrapeCommandHandler:          scrapeCommandHandler,
		networkPoliciesCommandHandler: networkPoliciesCommandHandler,
		clusterStatusUpdater:          clusterStatusUpdater,
		configCommandHandler:          configCommandHandler,
		components:                    otherComponents,

		stopC:    concurrency.NewErrorSignal(),
		stoppedC: concurrency.NewErrorSignal(),
	}
}
