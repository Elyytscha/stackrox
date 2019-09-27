package compliance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/internalapi/compliance"
	"github.com/stackrox/rox/generated/internalapi/sensor"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/k8sutil"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/orchestrators"
	"github.com/stackrox/rox/pkg/retry"
	"github.com/stackrox/rox/sensor/common/roxmetadata"
)

const (
	scrapeServiceName = "stackrox-compliance"
	scrapeCommand     = "stackrox/compliance"
	scrapeEnvironment = "ROX_SCRAPE_ID"

	benchmarkServiceAccount = "benchmark"
)

var (
	scrapeMounts = []string{
		"/etc:/host/etc:ro",
		"/lib:/host/lib:ro",
		"/usr/bin:/host/usr/bin:ro",
		"/usr/lib:/host/usr/lib:ro",
		"/var/lib:/host/var/lib:ro",
		"/var/log/audit:/host/var/log/audit:ro",
		"/run:/host/run",
		"/var/run/docker.sock:/host/var/run/docker.sock",
	}

	log = logging.LoggerForModule()
)

type commandHandlerImpl struct {
	orchestrator orchestrators.Orchestrator

	roxMetadata roxmetadata.Metadata

	commands chan *central.ScrapeCommand
	updates  chan *central.ScrapeUpdate

	scrapeIDToState map[string]*scrapeState

	stopC    concurrency.ErrorSignal
	stoppedC concurrency.ErrorSignal

	resultsC chan *compliance.ComplianceReturn
}

func (c *commandHandlerImpl) Start() {
	go c.run()
}

func (c *commandHandlerImpl) Stop(err error) {
	c.stopC.SignalWithError(err)
}

func (c *commandHandlerImpl) Stopped() concurrency.ReadOnlyErrorSignal {
	return &c.stoppedC
}

func (c *commandHandlerImpl) SendCommand(command *central.ScrapeCommand) bool {
	select {
	case c.commands <- command:
		return true
	case <-c.stoppedC.Done():
		return false
	}
}

func (c *commandHandlerImpl) Output() <-chan *central.ScrapeUpdate {
	return c.updates
}

func (c *commandHandlerImpl) run() {
	defer c.stoppedC.Signal()

	for {
		select {
		case <-c.stopC.Done():
			c.stoppedC.SignalWithError(c.stopC.Err())

		case command, ok := <-c.commands:
			if !ok {
				c.stoppedC.SignalWithError(errors.New("scrape command input closed"))
				return
			}
			if command.GetScrapeId() == "" {
				log.Errorf("received a command with no id: %s", proto.MarshalTextString(command))
				continue
			}
			if update := c.runCommand(command); update != nil {
				c.sendUpdate(update)
			}

		case result, ok := <-c.resultsC:
			if !ok {
				c.stoppedC.SignalWithError(errors.New("compliance return input closed"))
				return
			}
			if updates := c.commitResult(result); len(updates) > 0 {
				c.sendUpdates(updates)
			}
		}
	}
}

func (c *commandHandlerImpl) resultsChan() chan<- *compliance.ComplianceReturn {
	return c.resultsC
}

func (c *commandHandlerImpl) runCommand(command *central.ScrapeCommand) *central.ScrapeUpdate {
	switch command.Command.(type) {
	case *central.ScrapeCommand_StartScrape:
		return c.startScrape(command.GetScrapeId(), command.GetStartScrape().GetHostnames())
	case *central.ScrapeCommand_KillScrape:
		return c.killScrape(command.GetScrapeId())
	default:
		log.Errorf("unrecognized scrape command: %s", proto.MarshalTextString(command))
	}
	return nil
}

func (c *commandHandlerImpl) handleStartScrapeErr(scrapedID string, err error) *central.ScrapeUpdate {
	errStr := fmt.Sprintf("unable to start scrape %q: %v", scrapedID, err)
	log.Error(errStr)
	return scrapeStarted(scrapedID, errStr)
}

func (c *commandHandlerImpl) startScrape(scrapeID string, expectedHosts []string) *central.ScrapeUpdate {
	// Check that the scrape is not already running.
	if _, running := c.scrapeIDToState[scrapeID]; running {
		return nil
	}

	// Try to launch the scrape.
	// If we fail, return a message to central that we failed.
	svc, err := c.createService(scrapeID)
	if err != nil {
		return c.handleStartScrapeErr(scrapeID, err)
	}

	name, desired, err := c.orchestrator.LaunchDaemonSet(*svc)
	if err != nil {
		return c.handleStartScrapeErr(scrapeID, err)
	}

	// If we succeeded, start tracking the scrape and send a message to central.
	c.scrapeIDToState[scrapeID] = newScrapeState(name, desired, expectedHosts)
	log.Infof("started scrape %q with DaemonSet %q and %d desired nodes", scrapeID, name, desired)
	return scrapeStarted(scrapeID, "")
}

func (c *commandHandlerImpl) killScrape(scrapeID string) *central.ScrapeUpdate {
	// Check that the scrape has not already been killed.
	scrapeState, running := c.scrapeIDToState[scrapeID]
	if !running {
		return nil
	}

	// Try to kill the scrape.
	// If we fail, return a message to central that we failed.
	var errStr string
	if err := c.killScrapeWithRetry(scrapeState.deploymentName, scrapeID); err != nil {
		errStr = fmt.Sprintf("unable to kill scrape %s", err)
		log.Error(errStr)
		return scrapeKilled(scrapeID, errStr)
	}

	// If killed successfully, remove the scrape from tracking.
	delete(c.scrapeIDToState, scrapeID)
	return scrapeKilled(scrapeID, errStr)
}

// Helper function to kill a scrape with retry.
func (c *commandHandlerImpl) killScrapeWithRetry(name, scrapeID string) error {
	return retry.WithRetry(
		func() error {
			return c.orchestrator.Kill(name)
		},
		retry.Tries(5),
		retry.BetweenAttempts(func(_ int) {
			time.Sleep(time.Second)
		}),
		retry.OnFailedAttempts(func(err error) {
			log.Errorf("failed to kill scrape %s: %s", scrapeID, err)
		}),
	)
}

// Helper function that converts a scrape command into a SystemService that can be launched.
func (c *commandHandlerImpl) createService(scrapeID string) (*orchestrators.SystemService, error) {
	image := c.roxMetadata.GetSensorImage()
	if image == "" {
		return nil, errors.New("couldn't find sensor image to use")
	}
	zeroInt64 := int64(0)
	return &orchestrators.SystemService{
		GenerateName: scrapeServiceName,
		ExtraPodLabels: map[string]string{
			"compliance.stackrox.io/scrape-id": scrapeID,
			"com.stackrox.io/service":          "compliance",
		},
		Command: []string{scrapeCommand},
		Mounts:  scrapeMounts,
		HostPID: true,
		Envs: []string{
			env.CombineSetting(env.AdvertisedEndpoint),
			env.Combine(scrapeEnvironment, scrapeID),
		},
		SpecialEnvs: []orchestrators.SpecialEnvVar{
			orchestrators.NodeName,
		},
		Image:          image,
		ServiceAccount: benchmarkServiceAccount,
		Resources: &storage.Resources{
			// We want to request very little because this is a daemonset and we want it to be scheduled
			// on every node unless the node is literally full.
			CpuCoresRequest: 0.01,
			CpuCoresLimit:   1,
			MemoryMbRequest: 10,
			MemoryMbLimit:   2048,
		},
		Secrets: []orchestrators.Secret{
			{
				Name: "benchmark-tls",
				Items: map[string]string{
					"ca.pem":             "ca.pem",
					"benchmark-cert.pem": "cert.pem",
					"benchmark-key.pem":  "key.pem",
				},
				TargetPath: "/run/secrets/stackrox.io/certs/",
			},
		},
		RunAsUser: &zeroInt64,
	}, nil
}

func (c *commandHandlerImpl) commitResult(result *compliance.ComplianceReturn) (ret []*central.ScrapeUpdate) {
	// Check that the scrape has not already been killed.
	scrapeState, running := c.scrapeIDToState[result.GetScrapeId()]
	if !running {
		log.Errorf("received result for scrape not tracked: %q", result.GetScrapeId())
		return
	}

	// Check that we have not already received a result for the host.
	if scrapeState.foundNodes.Contains(result.GetNodeName()) {
		log.Errorf("received duplicate result in scrape %s for node %s", result.GetScrapeId(), result.GetNodeName())
		return
	}
	scrapeState.desiredNodes--

	// Check if the node did not exist in the scrape request
	if !scrapeState.remainingNodes.Contains(result.GetNodeName()) {
		log.Errorf("found node %s not requested by Central for scrape %s", result.GetNodeName(), result.GetScrapeId())
		return
	}

	// Pass the update back to central.
	scrapeState.remainingNodes.Remove(result.GetNodeName())
	ret = append(ret, scrapeUpdate(result))
	// If that was the last expected update, kill the scrape.
	if scrapeState.desiredNodes == 0 || scrapeState.remainingNodes.Cardinality() == 0 {
		if scrapeState.remainingNodes.Cardinality() != 0 {
			log.Warnf("compliance data for the following nodes was not collected: %+v", scrapeState.remainingNodes.AsSlice())
		}
		if update := c.killScrape(result.GetScrapeId()); update != nil {
			ret = append(ret, update)
		}
	}

	return
}

func (c *commandHandlerImpl) sendUpdates(updates []*central.ScrapeUpdate) {
	if len(updates) > 0 {
		for _, update := range updates {
			c.sendUpdate(update)
		}
	}
}

func (c *commandHandlerImpl) sendUpdate(update *central.ScrapeUpdate) {
	select {
	case <-c.stoppedC.Done():
		log.Errorf("failed to send update: %s", proto.MarshalTextString(update))
		return
	case c.updates <- update:
		return
	}
}

// GetScrapeConfig returns the scrape configuration for the given node name and scrape ID.
func (c *commandHandlerImpl) GetScrapeConfig(ctx context.Context, nodeName, scrapeID string) (*sensor.ScrapeConfig, error) {
	nodeInfo, err := c.orchestrator.GetNode(nodeName)
	if err != nil {
		return nil, err
	}

	rt, _ := k8sutil.ParseContainerRuntimeString(nodeInfo.Status.NodeInfo.ContainerRuntimeVersion)

	return &sensor.ScrapeConfig{
		ContainerRuntime: rt,
	}, nil
}

// Helper functions.
///////////////////

func scrapeStarted(scrapeID, err string) *central.ScrapeUpdate {
	return &central.ScrapeUpdate{
		ScrapeId: scrapeID,
		Update: &central.ScrapeUpdate_ScrapeStarted{
			ScrapeStarted: &central.ScrapeStarted{
				ErrorMessage: err,
			},
		},
	}
}

func scrapeKilled(scrapeID, err string) *central.ScrapeUpdate {
	return &central.ScrapeUpdate{
		ScrapeId: scrapeID,
		Update: &central.ScrapeUpdate_ScrapeKilled{
			ScrapeKilled: &central.ScrapeKilled{
				ErrorMessage: err,
			},
		},
	}
}

func scrapeUpdate(result *compliance.ComplianceReturn) *central.ScrapeUpdate {
	return &central.ScrapeUpdate{
		ScrapeId: result.GetScrapeId(),
		Update: &central.ScrapeUpdate_ComplianceReturn{
			ComplianceReturn: result,
		},
	}
}
