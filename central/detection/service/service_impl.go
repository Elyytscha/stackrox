package service

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	clusterDatastore "github.com/stackrox/rox/central/cluster/datastore"
	"github.com/stackrox/rox/central/detection"
	"github.com/stackrox/rox/central/detection/buildtime"
	"github.com/stackrox/rox/central/detection/deploytime"
	"github.com/stackrox/rox/central/enrichment"
	imageDatastore "github.com/stackrox/rox/central/image/datastore"
	"github.com/stackrox/rox/central/role/resources"
	apiV1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/grpc/authz"
	"github.com/stackrox/rox/pkg/grpc/authz/or"
	"github.com/stackrox/rox/pkg/grpc/authz/perrpc"
	"github.com/stackrox/rox/pkg/grpc/authz/user"
	"github.com/stackrox/rox/pkg/images/enricher"
	"github.com/stackrox/rox/pkg/images/types"
	"github.com/stackrox/rox/pkg/images/utils"
	"github.com/stackrox/rox/pkg/kubernetes"
	"github.com/stackrox/rox/pkg/logging"
	resourcesConv "github.com/stackrox/rox/pkg/protoconv/resources"
	"github.com/stackrox/rox/pkg/stringutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	coreV1 "k8s.io/api/core/v1"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

var (
	authorizer = perrpc.FromMap(map[authz.Authorizer][]string{
		user.With(permissions.View(resources.Detection)): {
			"/v1.DetectionService/DetectBuildTime",
			"/v1.DetectionService/DetectDeployTimeFromYAML",
		},
		or.SensorOrAuthorizer(user.With(permissions.Modify(resources.Detection))): {
			"/v1.DetectionService/DetectDeployTime",
		},
	})

	log = logging.LoggerForModule()
)

// serviceImpl provides APIs for alerts.
type serviceImpl struct {
	policySet          detection.PolicySet
	imageEnricher      enricher.ImageEnricher
	imageDatastore     imageDatastore.DataStore
	deploymentEnricher enrichment.Enricher
	buildTimeDetector  buildtime.Detector
	clusters           clusterDatastore.DataStore

	detector deploytime.Detector
}

// RegisterServiceServer registers this service with the given gRPC Server.
func (s *serviceImpl) RegisterServiceServer(grpcServer *grpc.Server) {
	apiV1.RegisterDetectionServiceServer(grpcServer, s)
}

// RegisterServiceHandlerFromEndpoint registers this service with the given gRPC Gateway endpoint.
func (s *serviceImpl) RegisterServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return apiV1.RegisterDetectionServiceHandler(ctx, mux, conn)
}

// AuthFuncOverride specifies the auth criteria for this API.
func (s *serviceImpl) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, authorizer.Authorized(ctx, fullMethodName)
}

// DetectBuildTime runs detection on a built image.
func (s *serviceImpl) DetectBuildTime(ctx context.Context, req *apiV1.BuildDetectionRequest) (*apiV1.BuildDetectionResponse, error) {
	image := req.GetImage()
	if req.GetImageName() != "" {
		var err error
		image, err = utils.GenerateImageFromString(req.GetImageName())
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if image.GetName() == nil {
		return nil, errors.New("image or image_name must be specified")
	}
	// This is a workaround for those who post the full image, but don't fill in fullname
	if name := image.GetName(); name != nil && name.GetFullName() == "" {
		name.FullName = types.Wrapper{GenericImage: image}.FullName()
	}

	img := types.ToImage(image)

	eCtx := enricher.EnrichmentContext{}
	if req.GetNoExternalMetadata() {
		eCtx.FetchOpt = enricher.NoExternalMetadata
	}
	enrichResult, err := s.imageEnricher.EnrichImage(eCtx, img)
	if err != nil {
		return nil, err
	}
	if enrichResult.ImageUpdated {
		img.Id = stringutils.FirstNonEmpty(img.GetId(), img.GetMetadata().GetV2().GetDigest(), img.GetMetadata().GetV1().GetDigest())
		if img.GetId() != "" {
			if err := s.imageDatastore.UpsertImage(ctx, img); err != nil {
				return nil, err
			}
		}
	}

	alerts, err := s.buildTimeDetector.Detect(img)
	if err != nil {
		return nil, err
	}
	return &apiV1.BuildDetectionResponse{
		Alerts: alerts,
	}, nil
}

// getImageID looks for any possible IDs from the image including from the fetched metadata
func getImageID(image *storage.Image) string {
	return stringutils.FirstNonEmpty(image.GetId(), image.GetMetadata().GetV2().GetDigest(), image.GetMetadata().GetV1().GetDigest())
}

func (s *serviceImpl) enrichAndDetect(ctx context.Context, enrichmentContext enricher.EnrichmentContext, deployment *storage.Deployment) (*apiV1.DeployDetectionResponse_Run, error) {
	images, updatedIndices, _, err := s.deploymentEnricher.EnrichDeployment(enrichmentContext, deployment)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	for _, idx := range updatedIndices {
		img := images[idx]
		img.Id = getImageID(img)
		if err := s.imageDatastore.UpsertImage(ctx, images[idx]); err != nil {
			return nil, err
		}
	}

	detectionCtx := deploytime.DetectionContext{
		EnforcementOnly: enrichmentContext.EnforcementOnly,
	}

	alerts, err := s.detector.Detect(detectionCtx, deployment, images)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &apiV1.DeployDetectionResponse_Run{
		Name:   deployment.GetName(),
		Type:   deployment.GetType(),
		Alerts: alerts,
	}, nil
}

func (s *serviceImpl) runDeployTimeDetect(ctx context.Context, eCtx enricher.EnrichmentContext, obj k8sRuntime.Object) (*apiV1.DeployDetectionResponse_Run, error) {
	if !kubernetes.IsDeploymentResource(obj.GetObjectKind().GroupVersionKind().Kind) {
		return nil, nil
	}

	deployment, err := resourcesConv.NewDeploymentFromStaticResource(obj, obj.GetObjectKind().GroupVersionKind().Kind, "")
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Could not convert to deployment from resource: %v", err)
	}
	return s.enrichAndDetect(ctx, eCtx, deployment)
}

func getObjectsFromYAML(yamlString string) ([]k8sRuntime.Object, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	reader := yaml.NewYAMLReader(bufio.NewReader(bytes.NewBufferString(yamlString)))
	var objects []k8sRuntime.Object
	var err error
	for err == nil {
		yamlBytes, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Failed to read YAML with err: %v", err)
		}
		obj, _, err := decode(yamlBytes, nil, nil)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "could not parse YAML: %v", err)
		}
		if list, ok := obj.(*coreV1.List); ok {
			listResources, err := getObjectsFromList(list)
			if err != nil {
				return nil, err
			}
			objects = append(objects, listResources...)
		} else {
			objects = append(objects, obj)
		}
	}
	return objects, nil
}

func getObjectsFromList(list *coreV1.List) ([]k8sRuntime.Object, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	objects := make([]k8sRuntime.Object, 0, len(list.Items))
	for i, item := range list.Items {
		obj, _, err := decode(item.Raw, nil, nil)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Could not decode item %d in the list: %v", i, err)
		}
		objects = append(objects, obj)
	}
	return objects, nil
}

// DetectDeployTime runs detection on a deployment
func (s *serviceImpl) DetectDeployTimeFromYAML(ctx context.Context, req *apiV1.DeployYAMLDetectionRequest) (*apiV1.DeployDetectionResponse, error) {
	if req.GetYaml() == "" {
		return nil, status.Error(codes.InvalidArgument, "yaml field must be specified in detection request")
	}

	resources, err := getObjectsFromYAML(req.GetYaml())
	if err != nil {
		return nil, err
	}

	eCtx := enricher.EnrichmentContext{
		EnforcementOnly: req.GetEnforcementOnly(),
	}
	if req.GetNoExternalMetadata() {
		eCtx.FetchOpt = enricher.NoExternalMetadata
	}

	var runs []*apiV1.DeployDetectionResponse_Run
	for _, r := range resources {
		run, err := s.runDeployTimeDetect(ctx, eCtx, r)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Unable to convert object: %v", err)
		}
		if run != nil {
			runs = append(runs, run)
		}
	}
	return &apiV1.DeployDetectionResponse{
		Runs: runs,
	}, nil
}

func isDeployTimeEnforcement(actions []storage.EnforcementAction) bool {
	for _, a := range actions {
		if a == storage.EnforcementAction_SCALE_TO_ZERO_ENFORCEMENT || a == storage.EnforcementAction_UNSATISFIABLE_NODE_CONSTRAINT_ENFORCEMENT {
			return true
		}
	}
	return false
}

func (s *serviceImpl) populateDeploymentWithClusterInfo(ctx context.Context, clusterID string, deployment *storage.Deployment) error {
	if clusterID == "" {
		return nil
	}
	clusterName, exists, err := s.clusters.GetClusterName(ctx, clusterID)
	if err != nil {
		return err
	}
	if !exists {
		return status.Errorf(codes.InvalidArgument, "cluster with ID %q does not exist", clusterID)
	}
	deployment.ClusterId = clusterID
	deployment.ClusterName = clusterName
	return nil
}

func (s *serviceImpl) DetectDeployTime(ctx context.Context, req *apiV1.DeployDetectionRequest) (*apiV1.DeployDetectionResponse, error) {
	if req.GetDeployment() == nil {
		return nil, status.Error(codes.InvalidArgument, "Deployment must be passed to deploy time detection")
	}
	if err := s.populateDeploymentWithClusterInfo(ctx, req.GetClusterId(), req.GetDeployment()); err != nil {
		return nil, err
	}

	// If we have enforcement only, then check if any of the policies need enforcement. If not, then just exit with no alerts generated
	if req.GetEnforcementOnly() {
		var evaluationRequired bool
		_ = s.policySet.ForEach(detection.FunctionAsExecutor(func(compiled detection.CompiledPolicy) error {
			if isDeployTimeEnforcement(compiled.Policy().GetEnforcementActions()) {
				evaluationRequired = true
				return errors.New("not a real error, just early exits this foreach")
			}
			return nil
		}))
		if !evaluationRequired {
			return &apiV1.DeployDetectionResponse{
				Runs: []*apiV1.DeployDetectionResponse_Run{
					{
						Name:   req.GetDeployment().GetName(),
						Type:   req.GetDeployment().GetType(),
						Alerts: nil,
					},
				},
			}, nil
		}
	}

	enrichmentCtx := enricher.EnrichmentContext{
		EnforcementOnly: req.GetEnforcementOnly(),
	}
	if req.GetNoExternalMetadata() {
		enrichmentCtx.FetchOpt = enricher.NoExternalMetadata
	}

	run, err := s.enrichAndDetect(ctx, enrichmentCtx, req.GetDeployment())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &apiV1.DeployDetectionResponse{
		Runs: []*apiV1.DeployDetectionResponse_Run{
			run,
		},
	}, nil
}
