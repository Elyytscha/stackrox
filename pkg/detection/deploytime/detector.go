package deploytime

import (
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/detection"
)

// DetectionContext is the context for detection
type DetectionContext struct {
	EnforcementOnly bool
}

type signatureVerifier = func(image *storage.Image) bool

// Detector provides an interface for getting and managing alerts and enforcements on deployments.
type Detector interface {
	PolicySet() detection.PolicySet

	Detect(ctx DetectionContext, deployment *storage.Deployment, images []*storage.Image, filters ...detection.FilterOption) ([]*storage.Alert, error)
}

// NewDetector returns a new instance of a Detector.
func NewDetector(policySet detection.PolicySet) Detector {
	return &detectorImpl{
		policySet: policySet,
	}
}
