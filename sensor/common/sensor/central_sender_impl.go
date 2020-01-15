package sensor

import (
	"errors"

	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/listeners"
	"github.com/stackrox/rox/sensor/common"
	"github.com/stackrox/rox/sensor/common/clusterstatus"
	"github.com/stackrox/rox/sensor/common/compliance"
	"github.com/stackrox/rox/sensor/common/config"
	"github.com/stackrox/rox/sensor/common/deduper"
	"github.com/stackrox/rox/sensor/common/metrics"
	networkConnManager "github.com/stackrox/rox/sensor/common/networkflow/manager"
	"github.com/stackrox/rox/sensor/common/networkpolicies"
	"github.com/stackrox/rox/sensor/common/signal"
)

type centralSenderImpl struct {
	// Generate messages to be sent to central.
	listener                      listeners.Listener
	signalService                 signal.Service
	networkConnManager            networkConnManager.Manager
	scrapeCommandHandler          compliance.CommandHandler
	configCommandHandler          config.Handler
	networkPoliciesCommandHandler networkpolicies.CommandHandler
	clusterStatusUpdater          clusterstatus.Updater
	components                    []common.SensorComponent

	stopC    concurrency.ErrorSignal
	stoppedC concurrency.ErrorSignal
}

func (s *centralSenderImpl) Start(stream central.SensorService_CommunicateClient, onStops ...func(error)) {
	go s.send(stream, onStops...)
}

func (s *centralSenderImpl) Stop(err error) {
	s.stopC.SignalWithError(err)
}

func (s *centralSenderImpl) Stopped() concurrency.ReadOnlyErrorSignal {
	return &s.stoppedC
}

func (s *centralSenderImpl) forwardResponses(from <-chan *central.MsgFromSensor, to chan<- *central.MsgFromSensor) {
	for !s.stopC.IsDone() {
		select {
		case msg, ok := <-from:
			if !ok {
				return
			}
			select {
			case to <- msg:
			case <-s.stopC.Done():
				return
			}
		case <-s.stopC.Done():
			return
		}
	}
}

func (s *centralSenderImpl) send(stream central.SensorService_CommunicateClient, onStops ...func(error)) {
	defer func() {
		s.stoppedC.SignalWithError(s.stopC.Err())
		runAll(s.stopC.Err(), onStops...)
	}()

	wrappedStream := metrics.NewCountingEventStream(stream, "unique")
	wrappedStream = deduper.NewDedupingMessageStream(wrappedStream)
	wrappedStream = metrics.NewCountingEventStream(wrappedStream, "total")

	componentMsgsC := make(chan *central.MsgFromSensor)
	for _, component := range s.components {
		if responsesC := component.ResponsesC(); responsesC != nil {
			go s.forwardResponses(responsesC, componentMsgsC)
		}
	}

	// NB: The centralSenderImpl reserves the right to perform arbitrary reads and writes on the returned objects.
	// The providers that send the messages below are responsible for making sure that once they send events here,
	// they do not use the objects again in any way, since that can result in a race condition caused by concurrent
	// reads and writes.
	// Ideally, if you're going to continue to hold a reference to the object, you want to proto.Clone it before
	// sending it to this function.
	for {
		var msg *central.MsgFromSensor

		select {
		case sig, ok := <-s.signalService.Indicators():
			if !ok {
				s.stopC.SignalWithError(errors.New("signals channel closed"))
				return
			}
			msg = &central.MsgFromSensor{
				Msg: &central.MsgFromSensor_Event{
					Event: sig,
				},
			}
		case evt, ok := <-s.listener.Events():
			if !ok {
				s.stopC.SignalWithError(errors.New("orchestrator events channel closed"))
				return
			}
			msg = &central.MsgFromSensor{
				Msg: &central.MsgFromSensor_Event{
					Event: evt,
				},
			}
		case flowUpdate, ok := <-s.networkConnManager.FlowUpdates():
			if !ok {
				s.stopC.SignalWithError(errors.New("flow updates channel closed"))
				return
			}
			msg = &central.MsgFromSensor{
				Msg: &central.MsgFromSensor_NetworkFlowUpdate{
					NetworkFlowUpdate: flowUpdate,
				},
			}
		case scrapeUpdate, ok := <-s.scrapeCommandHandler.Output():
			if !ok {
				s.stopC.SignalWithError(errors.New("scrape command handler channel closed"))
				return
			}
			msg = &central.MsgFromSensor{
				Msg: &central.MsgFromSensor_ScrapeUpdate{
					ScrapeUpdate: scrapeUpdate,
				},
			}
		case networkPolicyCommandResponse, ok := <-s.networkPoliciesCommandHandler.Responses():
			if !ok {
				s.stopC.SignalWithError(errors.New("network policies command handler channel closed"))
				return
			}
			msg = &central.MsgFromSensor{
				Msg: &central.MsgFromSensor_NetworkPoliciesResponse{
					NetworkPoliciesResponse: networkPolicyCommandResponse,
				},
			}
		case clusterStatusUpdate, ok := <-s.clusterStatusUpdater.Updates():
			if !ok {
				s.stopC.SignalWithError(errors.New("cluster status update handler closed"))
				return
			}
			msg = &central.MsgFromSensor{
				Msg: &central.MsgFromSensor_ClusterStatusUpdate{
					ClusterStatusUpdate: clusterStatusUpdate,
				},
			}
		case componentMsg := <-componentMsgsC:
			msg = componentMsg
		case <-s.stopC.Done():
			return
		case <-stream.Context().Done():
			s.stopC.SignalWithError(stream.Context().Err())
			return
		}

		if msg != nil {
			if err := wrappedStream.Send(msg); err != nil {
				s.stopC.SignalWithError(err)
				return
			}
		}
	}
}
