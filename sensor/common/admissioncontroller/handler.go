package admissioncontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/enforcers"
	"github.com/stackrox/rox/pkg/kubernetes"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/protoconv/resources"
	"github.com/stackrox/rox/pkg/stringutils"
	"github.com/stackrox/rox/pkg/templates"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/sensor/common/clusterid"
	"google.golang.org/grpc"
	admission "k8s.io/api/admission/v1beta1"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	// This purposefully leaves a newline at the top for formatting when using kubectl
	kubectlTemplate = `
Policy: {{.Title}}
{{if not .DisableBypass}}In case of emergency, add the annotation {"admission.stackrox.io/break-glass": "ticket-1234"} to your deployment with an updated ticket number
{{- end}}
{{range .Alerts}}
{{.Policy.Name}}
- Description:
    ↳ {{wrap .Policy.Description}}
- Rationale:
    ↳ {{wrap .Policy.Rationale}}
- Remediation:
    ↳ {{wrap .Policy.Remediation}}
- Violations:
    {{- range .Violations}}
    - {{.Message}}
    {{- end}}
{{ end}}
`
)

var (
	log = logging.LoggerForModule()

	msgTemplate = template.Must(template.New("name").Funcs(
		template.FuncMap{
			"wrap": stringutils.Wrap,
		}).Parse(kubectlTemplate))
)

// DynamicConfigProvider abstracts access to the dynamic cluster configuration.
type DynamicConfigProvider interface {
	GetConfig() *storage.DynamicClusterConfig
}

// NewHandler returns a handler that proxies admission controllers to Central
func NewHandler(conn grpc.ClientConnInterface, centralReachable *concurrency.Flag, configProvider DynamicConfigProvider) http.Handler {
	return &handlerImpl{
		client:           v1.NewDetectionServiceClient(conn),
		centralReachable: centralReachable,
		configProvider:   configProvider,
	}
}

type handlerImpl struct {
	client           v1.DetectionServiceClient
	centralReachable *concurrency.Flag
	configProvider   DynamicConfigProvider
}

func admissionPass(w http.ResponseWriter, id types.UID) {
	writeResponse(w, id, true, "")
}

func writeResponse(w http.ResponseWriter, id types.UID, allowed bool, reason string) {
	var ar *admission.AdmissionReview
	if allowed {
		ar = &admission.AdmissionReview{
			Response: &admission.AdmissionResponse{
				UID:     id,
				Allowed: true,
			},
		}
	} else {
		ar = &admission.AdmissionReview{
			Response: &admission.AdmissionResponse{
				UID:     id,
				Allowed: false,
				Result: &metav1.Status{
					Status:  "Failure",
					Reason:  metav1.StatusReason("Failed currently enforced policies from StackRox"),
					Message: reason,
				},
			},
		}
	}

	data, err := json.Marshal(ar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parseIntoDeployment(ar *admission.AdmissionReview, registryOverride string) (*storage.Deployment, error) {
	var objType interface{}
	if ar.Request == nil {
		return nil, nil
	}
	switch ar.Request.Kind.Kind {
	case kubernetes.Pod:
		objType = &core.Pod{}
	case kubernetes.Deployment:
		objType = &apps.Deployment{}
	case kubernetes.StatefulSet:
		objType = &apps.StatefulSet{}
	case kubernetes.DaemonSet:
		objType = &apps.DaemonSet{}
	case kubernetes.ReplicationController:
		objType = &core.ReplicationController{}
	case kubernetes.ReplicaSet:
		objType = &apps.ReplicaSet{}
	default:
		return nil, errors.Errorf("currently do not recognize kind %q in admission controller", ar.Request.Kind.Kind)
	}

	if err := json.Unmarshal(ar.Request.Object.Raw, &objType); err != nil {
		return nil, err
	}

	return resources.NewDeploymentFromStaticResource(objType, ar.Request.Kind.Kind, registryOverride)
}

// ServeHTTP serves the admission controller endpoint
func (s *handlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_, _ = w.Write([]byte("{}"))
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Endpoint only supports GET and POST requests", http.StatusBadRequest)
		return
	}

	if !s.centralReachable.Get() {
		http.Error(w, "Connection to central has not yet been established. Cannot handle admission controller requests", http.StatusServiceUnavailable)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var admissionReview admission.AdmissionReview
	if err := decoder.Decode(&admissionReview); err != nil {
		log.Errorf("Error decoding admission review: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if admissionReview.Request == nil {
		errMsg := fmt.Sprintf("invalid admission review. nil request: %+v", admissionReview)
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// Can guarantee that there is an admission controller config because we receive it between central reachable is set to true
	conf := s.configProvider.GetConfig().GetAdmissionControllerConfig()

	// If admission controller is not enabled in dynamic config, then always pass
	if !conf.GetEnabled() {
		admissionPass(w, admissionReview.Request.UID)
		return
	}

	deployment, err := parseIntoDeployment(&admissionReview, s.configProvider.GetConfig().GetRegistryOverride())
	if err != nil {
		log.Errorf("error parsing into deployment: %v", err)
		admissionPass(w, admissionReview.Request.UID)
		return
	}

	// A nil deployment implies that it does not need to be processed. For example,
	// if there is an owner reference then we are assuming that we ran the admission controller
	// on the higher level object
	if deployment == nil {
		admissionPass(w, admissionReview.Request.UID)
		return
	}

	// This checks to see if the deployment has a bypass annotation only if the bypass is not disabled
	if !conf.GetDisableBypass() && !enforcers.ShouldEnforce(deployment.GetAnnotations()) {
		log.Warnf("deployment %s/%s of type %s was deployed without being checked due to matching bypass annotation %q",
			deployment.GetNamespace(), deployment.GetName(), deployment.GetType(), enforcers.EnforcementBypassAnnotationKey)
		admissionPass(w, admissionReview.Request.UID)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(conf.GetTimeoutSeconds())*time.Second)
	defer cancel()

	resp, err := s.client.DetectDeployTime(ctx, &v1.DeployDetectionRequest{
		Resource: &v1.DeployDetectionRequest_Deployment{
			Deployment: deployment,
		},
		NoExternalMetadata: !conf.GetScanInline(),
		EnforcementOnly:    true,
		ClusterId:          clusterid.Get(),
	})
	if err != nil {
		log.Errorf("Deployment %s/%s of type %s was deployed without being checked due to detection error: %v", deployment.GetNamespace(), deployment.GetName(), deployment.GetType(), err)
		return
	}

	var enforcedAlerts []*storage.Alert
	var totalPolicies int
	// There will only ever be one run in this call
	for _, r := range resp.GetRuns() {
		for _, a := range r.GetAlerts() {
			totalPolicies++
			for _, e := range a.GetPolicy().GetEnforcementActions() {
				if e == storage.EnforcementAction_SCALE_TO_ZERO_ENFORCEMENT ||
					e == storage.EnforcementAction_UNSATISFIABLE_NODE_CONSTRAINT_ENFORCEMENT {
					enforcedAlerts = append(enforcedAlerts, a)
					break
				}
			}
		}
	}

	var topMsg string
	if len(enforcedAlerts) == 0 {
		admissionPass(w, admissionReview.Request.UID)
		return
	}
	if len(enforcedAlerts) == 1 {
		topMsg = fmt.Sprintf("Violated %d policies total. 1 enforced policy is described below:", totalPolicies)
	} else {
		topMsg = fmt.Sprintf("Violated %d policies total. %d enforced policies are described below:", totalPolicies, len(enforcedAlerts))
	}

	msg, err := templates.ExecuteToString(msgTemplate, map[string]interface{}{
		"Title":         topMsg,
		"Alerts":        enforcedAlerts,
		"DisableBypass": conf.GetDisableBypass(),
	})
	if err != nil {
		msg = fmt.Sprintf("internal failure executing admission controller msg template: %v", err)
		utils.Should(errors.New(msg))
	}
	writeResponse(w, admissionReview.Request.UID, false, msg)
}
