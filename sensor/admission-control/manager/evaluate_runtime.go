package manager

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/internalapi/sensor"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/kubernetes"
	"github.com/stackrox/rox/pkg/stringutils"
	admission "k8s.io/api/admission/v1beta1"
)

func (m *manager) shouldBypassRuntimeDetection(s *state, req *admission.AdmissionRequest) bool {
	if s.allRuntimePoliciesDetector == nil {
		log.Debugf("Runtime policy matcher not found, bypassing %s request on %s/%s [%s]", req.Operation, req.Namespace, req.Name, req.Kind)
		return true
	}

	// Allow the request if it comes from a blessed user
	if s.bypassForUsers.Contains(req.UserInfo.Username) {
		log.Debugf("Request comes from privileged user %s, bypassing %s request on %s/%s [%s]", req.UserInfo.Username, req.Operation, req.Namespace, req.Name, req.Kind)
		return true
	}

	// Allow the request if it comes from a service account in a system namespace
	if strings.HasPrefix(req.UserInfo.Username, "system:serviceaccount:") {
		saNamespace, _ := stringutils.Split2(req.UserInfo.Username[len("system:serviceaccount:"):], ":")
		if kubernetes.SystemNamespaceSet.Contains(saNamespace) {
			log.Debugf("Request comes from a system service account %s, bypassing %s request on %s/%s [%s]", req.UserInfo.Username, req.Operation, req.Namespace, req.Name, req.Kind)
			return true
		}
	}

	// Allow the request if it comes from a blessed group
	for _, group := range req.UserInfo.Groups {
		if s.bypassForGroups.Contains(group) {
			log.Debugf("Request comes from privileged group %s, bypassing %s request on %s/%s [%s]", group, req.Operation, req.Namespace, req.Name, req.Kind)
			return true
		}
	}

	return false
}

func (m *manager) evaluateRuntimeAdmissionRequest(s *state, req *admission.AdmissionRequest) (*admission.AdmissionResponse, error) {
	log.Debugf("Evaluating request %+v", req)
	if !features.K8sEventDetection.Enabled() {
		return pass(req.UID), nil
	}

	if m.shouldBypassRuntimeDetection(s, req) {
		return pass(req.UID), nil
	}

	log.Debugf("Not bypassing %s request on %s/%s [%s]", req.Operation, req.Namespace, req.Name, req.Kind)

	event, err := kubernetes.AdmissionRequestToKubeEventObj(req)
	if err != nil {
		return nil, errors.Wrap(err, "translating admission request object from request")
	}

	log.Debugf("Evaluating policies on kubernetes request %s", kubernetes.EventAsString(event))

	alerts, enrichedWithDeployment, err := m.evaluatePodEvent(s, event)
	if err != nil {
		return nil, errors.Wrap(err, "running StackRox detection")
	}

	if len(alerts) == 0 {
		log.Debugf("No policies violated, allowing %s request on %s/%s [%s]", req.Operation, req.Namespace, req.Name, req.Kind)
		return pass(req.UID), nil
	}

	if failReviewRequest(alerts...) {
		// TODO: Mark enforced violations as attempted and send to Sensor.
		return fail(req.UID, message(alerts, false)), nil
	}

	if enrichedWithDeployment {
		go m.putAlertsOnChan(alerts)
	}
	return pass(req.UID), nil
}

func (m *manager) evaluatePodEvent(s *state, event *storage.KubernetesEvent) ([]*storage.Alert, bool, error) {
	deployment := m.getDeploymentForPod(event.GetObject().GetNamespace(), event.GetObject().GetName())
	if deployment != nil {
		log.Debugf("Found deployment %s (id=%s) for %s/%s", deployment.GetName(), deployment.GetId(),
			event.GetObject().GetNamespace(), event.GetObject().GetName())

		var fetchImgCtx context.Context
		if timeoutSecs := s.GetClusterConfig().GetAdmissionControllerConfig().GetTimeoutSeconds(); timeoutSecs > 1 {
			var cancel context.CancelFunc
			fetchImgCtx, cancel = context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
			defer cancel()
		}

		getAlertsFunc := func(dep *storage.Deployment, imgs []*storage.Image) ([]*storage.Alert, error) {
			return s.allRuntimePoliciesDetector.DetectForDeployment(dep, imgs, nil, false, event)
		}

		alerts, err := m.kickOffImgScansAndDetect(fetchImgCtx, s, getAlertsFunc, deployment)
		if err != nil {
			return nil, false, err
		}
		return alerts, true, nil
	}

	// If deployment is not available , detect without deployment to respond to admission review. Run detection with
	// deployment enrichment in the background and record it.
	log.Warnf("Deployment for %s/%s not found. "+
		"Policies with deploy-time fields for kubernetes event %s will be detected in background",
		event.GetObject().GetNamespace(), event.GetObject().GetName(), kubernetes.EventAsString(event))

	go m.waitForDeploymentAndDetect(s, event)

	alerts, err := s.runtimeDetectorForPoliciesWithoutDeployFields.DetectForDeployment(&storage.Deployment{}, nil, nil, false, event)
	if err != nil {
		return nil, false, err
	}
	return alerts, false, nil
}

func (m *manager) waitForDeploymentAndDetect(s *state, event *storage.KubernetesEvent) {
	select {
	case <-m.stopSig.Done():
		return
	case <-m.initialSyncSig.Done():
		deployment := m.getDeploymentForPod(event.GetObject().GetNamespace(), event.GetObject().GetName())
		if deployment == nil {
			dep, err := m.depClient.GetDeploymentForPod(context.Background(), &sensor.GetDeploymentForPodRequest{
				PodName:   event.GetObject().GetName(),
				Namespace: event.GetObject().GetNamespace(),
			})
			if err != nil {
				log.Errorf("Could not fetch deployment for namespace/%s/pod/%s from Sensor. ",
					event.GetObject().GetNamespace(), event.GetObject().GetName())
				return
			}
			if dep == nil {
				return
			}
		}

		log.Debugf("Found deployment %s (id=%s) for %s/%s", deployment.GetName(), deployment.GetId(),
			event.GetObject().GetNamespace(), event.GetObject().GetName())

		var fetchImgCtx context.Context
		if timeoutSecs := s.GetClusterConfig().GetAdmissionControllerConfig().GetTimeoutSeconds(); timeoutSecs > 1 {
			var cancel context.CancelFunc
			fetchImgCtx, cancel = context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
			defer cancel()
		}

		getAlertsFunc := func(dep *storage.Deployment, imgs []*storage.Image) ([]*storage.Alert, error) {
			return s.runtimeDetectorForPoliciesWithDeployFields.DetectForDeployment(dep, imgs, nil, false, event)
		}

		alerts, err := m.kickOffImgScansAndDetect(fetchImgCtx, s, getAlertsFunc, deployment)
		if err != nil {
			log.Errorf("Failed to run StackRox detection: %v", err)
			return
		}
		if len(alerts) == 0 {
			return
		}

		go m.putAlertsOnChan(alerts)

		return
	}
}

func failReviewRequest(alerts ...*storage.Alert) bool {
	for _, alert := range alerts {
		if alert.GetEnforcement().GetAction() == storage.EnforcementAction_FAIL_KUBE_REQUEST_ENFORCEMENT {
			return true
		}
	}
	return false
}
