package deploytime

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	deploymentIndexer "github.com/stackrox/rox/central/deployment/index"
	"github.com/stackrox/rox/central/detection"
	"github.com/stackrox/rox/central/globalindex"
	imageIndexer "github.com/stackrox/rox/central/image/index"
	"github.com/stackrox/rox/central/searchbasedpolicies"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
)

func newSingleDeploymentExecutor(ctx DetectionContext, searcher search.Searcher, deployment *storage.Deployment) alertCollectingExecutor {
	return &policyExecutor{
		ctx:        ctx,
		searcher:   searcher,
		deployment: deployment,
	}
}

type policyExecutor struct {
	ctx        DetectionContext
	searcher   search.Searcher
	deployment *storage.Deployment
	alerts     []*storage.Alert
}

func (d *policyExecutor) GetAlerts() []*storage.Alert {
	return d.alerts
}

func (d *policyExecutor) ClearAlerts() {
	d.alerts = nil
}

func (d *policyExecutor) Execute(compiled detection.CompiledPolicy) error {
	// Check predicate on deployment.
	if !compiled.AppliesTo(d.deployment) {
		return nil
	}

	// Check enforcement on deployment if we don't want unenforced alerts.
	enforcement, _ := buildEnforcement(compiled.Policy(), d.deployment)
	if enforcement == storage.EnforcementAction_UNSET_ENFORCEMENT && d.ctx.EnforcementOnly {
		return nil
	}

	// Generate violations.
	violations, err := d.getViolations(enforcement, compiled.Matcher())
	if err != nil {
		return errors.Wrapf(err, "evaluating violations for policy %s; deployment %s/%s", compiled.Policy().GetName(), d.deployment.GetNamespace(), d.deployment.GetName())
	}
	if len(violations) > 0 {
		d.alerts = append(d.alerts, policyDeploymentAndViolationsToAlert(compiled.Policy(), d.deployment, violations))
	}
	return nil
}

func (d *policyExecutor) getViolations(enforcement storage.EnforcementAction, matcher searchbasedpolicies.Matcher) ([]*storage.Alert_Violation, error) {
	var err error
	var violations []*storage.Alert_Violation
	if enforcement != storage.EnforcementAction_UNSET_ENFORCEMENT {
		violations, err = matchWithEmptyImageIDs(matcher, d.deployment)
	} else {
		var violationsWrapper searchbasedpolicies.Violations
		// Purposefully, use searcher for deployment check
		violationsWrapper, err = matcher.MatchOne(d.searcher, d.deployment.GetId())
		violations = violationsWrapper.AlertViolations
	}
	return violations, err
}

func matchWithEmptyImageIDs(matcher searchbasedpolicies.Matcher, deployment *storage.Deployment) ([]*storage.Alert_Violation, error) {
	deploymentIndex, deployment, err := singleDeploymentSearcher(deployment)
	if err != nil {
		return nil, err
	}
	violations, err := matcher.MatchOne(deploymentIndex, deployment.GetId())
	if err != nil {
		return nil, err
	}
	return violations.AlertViolations, nil
}

const deploymentID = "deployment-id"

func singleDeploymentSearcher(deployment *storage.Deployment) (search.Searcher, *storage.Deployment, error) {
	clonedDeployment := proto.Clone(deployment).(*storage.Deployment)
	if clonedDeployment.GetId() == "" {
		clonedDeployment.Id = deploymentID
	}

	tempIndex, err := globalindex.MemOnlyIndex()
	if err != nil {
		return nil, nil, errors.Wrap(err, "initializing temp index")
	}

	imageIndex := imageIndexer.New(tempIndex)
	deploymentIndex := deploymentIndexer.New(tempIndex)
	for i, container := range clonedDeployment.GetContainers() {
		if container.GetImage() == nil {
			continue
		}
		if container.GetImage().GetId() == "" {
			container.Image.Id = fmt.Sprintf("image-id-%d", i)
		}
		if err := imageIndex.AddImage(container.GetImage()); err != nil {
			return nil, nil, err
		}
	}
	if err := deploymentIndex.AddDeployment(clonedDeployment); err != nil {
		return nil, nil, err
	}
	return deploymentIndex, clonedDeployment, nil
}
