package multipliers

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/risk"
	"github.com/stackrox/rox/pkg/set"
)

const (
	configSaturation = 8
	configMaxScore   = 2
)

// serviceConfigMultiplier is a scorer for the service configuration
type serviceConfigMultiplier struct{}

// NewServiceConfig scores the data based on the service configuration
func NewServiceConfig() Multiplier {
	return &serviceConfigMultiplier{}
}

// Score takes a deployment and evaluates its risk based on the service configuration
func (s *serviceConfigMultiplier) Score(_ context.Context, msg proto.Message) *storage.Risk_Result {
	deployment, ok := msg.(*storage.Deployment)
	if !ok {
		return nil
	}
	riskResult := &storage.Risk_Result{
		Name: risk.ServiceConfiguration.DisplayTitle,
	}
	var overallScore float32
	if volumeFactor := s.scoreVolumes(deployment); volumeFactor != "" {
		overallScore++
		riskResult.Factors = append(riskResult.Factors, &storage.Risk_Result_Factor{Message: volumeFactor})
	}
	if secretFactor := s.scoreSecrets(deployment); secretFactor != "" {
		overallScore++
		riskResult.Factors = append(riskResult.Factors, &storage.Risk_Result_Factor{Message: secretFactor})
	}
	capAddFactor, capDropFactor := s.scoreCapabilities(deployment)
	if capAddFactor != "" {
		overallScore++
		riskResult.Factors = append(riskResult.Factors, &storage.Risk_Result_Factor{Message: capAddFactor})
	}
	if capDropFactor != "" {
		overallScore++
		riskResult.Factors = append(riskResult.Factors, &storage.Risk_Result_Factor{Message: capDropFactor})
	}
	if factor := s.scorePrivilege(deployment); factor != "" {
		overallScore *= 2
		riskResult.Factors = append(riskResult.Factors, &storage.Risk_Result_Factor{Message: factor})
	}
	if len(riskResult.GetFactors()) == 0 {
		return nil
	}
	// riskResult.Score is the normalized [1.0,2.0] score
	riskResult.Score = normalizeScore(overallScore, configSaturation, configMaxScore)
	return riskResult
}

func (s *serviceConfigMultiplier) scoreVolumes(deployment *storage.Deployment) string {
	var volumeNames []string
	for _, container := range deployment.GetContainers() {
		for _, volume := range container.GetVolumes() {
			if !volume.GetReadOnly() {
				volumeNames = append(volumeNames, volume.GetName())
			}
		}
	}
	if len(volumeNames) != 0 {
		return fmt.Sprintf("Volumes %s were mounted RW", strings.Join(volumeNames, ", "))
	}
	return ""
}

func (s *serviceConfigMultiplier) scoreSecrets(deployment *storage.Deployment) string {
	var secrets []string
	for _, container := range deployment.GetContainers() {
		for _, secret := range container.GetSecrets() {
			secrets = append(secrets, secret.GetName())
		}
	}
	if len(secrets) != 0 {
		return fmt.Sprintf("Secrets %s are used inside the deployment", strings.Join(secrets, ", "))
	}
	return ""
}

var relevantCapAdds = set.NewFrozenStringSet(
	"ALL",
	"CAP_SYS_ADMIN",
	"CAP_NET_ADMIN",
	"CAP_SYS_MODULE",
)

func (s *serviceConfigMultiplier) scoreCapabilities(deployment *storage.Deployment) (capAddFactor, capDropFactor string) {
	capsAdded := set.NewStringSet()
	capsDropped := set.NewStringSet()
	for _, container := range deployment.GetContainers() {
		context := container.GetSecurityContext()
		for _, cap := range context.GetAddCapabilities() {
			if relevantCapAdds.Contains(cap) {
				capsAdded.Add(cap)
			}
		}
		for _, cap := range context.GetDropCapabilities() {
			capsDropped.Add(cap)
		}
	}
	if capsAdded.Cardinality() != 0 {
		addedSlice := capsAdded.AsSlice()
		capAddFactor = fmt.Sprintf("Capabilities %s were added", strings.Join(addedSlice, ", "))
	}
	if capsDropped.Cardinality() == 0 {
		capDropFactor = fmt.Sprintf("No capabilities were dropped")
	}
	return
}

func (s *serviceConfigMultiplier) scorePrivilege(deployment *storage.Deployment) string {
	for _, container := range deployment.GetContainers() {
		if container.GetSecurityContext().GetPrivileged() {
			return fmt.Sprintf("Container %q in the deployment is privileged", container.GetName())
		}
	}
	return ""
}
