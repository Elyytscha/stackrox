package resources

import (
	"encoding/json"
	"fmt"
	"reflect"

	ptypes "github.com/gogo/protobuf/types"
	openshift_appsv1 "github.com/openshift/api/apps/v1"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	imageUtils "github.com/stackrox/rox/pkg/images/utils"
	"github.com/stackrox/rox/pkg/kubernetes"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/protoconv/k8s"
	"github.com/stackrox/rox/pkg/protoconv/resources/volumes"
	"k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	openshiftEncodedDeploymentConfigAnnotation = `openshift.io/encoded-deployment-config`
)

var (
	log = logging.LoggerForModule()
)

// DeploymentWrap is a wrapper around a deployment to help convert the static fields
type DeploymentWrap struct {
	*storage.Deployment
}

// NewDeploymentWrap creates a deployment wrap from an existing deployment
func NewDeploymentWrap(d *storage.Deployment) *DeploymentWrap {
	return &DeploymentWrap{
		Deployment: d,
	}
}

// This checks if a reflect value is a Zero value, which means the field did not exist
func doesFieldExist(value reflect.Value) bool {
	return !reflect.DeepEqual(value, reflect.Value{})
}

func isTrackedOwnerReference(reference metav1.OwnerReference) bool {
	// Validate that the object is one that we are tracking as a Deployment
	if !kubernetes.IsDeploymentResource(reference.Kind) {
		return false
	}
	return kubernetes.IsNativeAPI(reference.APIVersion)
}

// NewDeploymentFromStaticResource returns a storage.Deployment from a k8s object
func NewDeploymentFromStaticResource(obj interface{}, deploymentType string) (*storage.Deployment, error) {
	objMeta, err := meta.Accessor(obj)
	if err != nil {
		return nil, errors.Wrapf(err, "could not access metadata of object of type %T", obj)
	}
	kind := deploymentType

	// Ignore resources that are owned by another tracked resource.
	for _, ref := range objMeta.GetOwnerReferences() {
		if isTrackedOwnerReference(ref) {
			return nil, nil
		}
	}

	// This only applies to OpenShift
	if encDeploymentConfig, ok := objMeta.GetLabels()[openshiftEncodedDeploymentConfigAnnotation]; ok {
		newMeta, newKind, err := extractDeploymentConfig(encDeploymentConfig)
		if err != nil {
			log.Error(err)
		} else {
			objMeta, kind = newMeta, newKind
		}
	}

	wrap := newWrap(objMeta, kind)
	wrap.populateFields(obj)
	return wrap.Deployment, nil

}

func extractDeploymentConfig(encodedDeploymentConfig string) (metav1.Object, string, error) {
	// Anonymous struct that only contains the fields we are interested in (note: json.Unmarshal silently ignores
	// fields that are not in the destination object).
	dc := struct {
		metav1.TypeMeta
		MetaData metav1.ObjectMeta `json:"metadata"`
	}{}
	err := json.Unmarshal([]byte(encodedDeploymentConfig), &dc)
	return &dc.MetaData, dc.Kind, err
}

func newWrap(meta metav1.Object, kind string) *DeploymentWrap {
	createdTime, err := ptypes.TimestampProto(meta.GetCreationTimestamp().Time)
	if err != nil {
		log.Error(err)
	}
	return &DeploymentWrap{
		Deployment: &storage.Deployment{
			Id:          string(meta.GetUID()),
			Name:        meta.GetName(),
			Type:        kind,
			Namespace:   meta.GetNamespace(),
			Labels:      meta.GetLabels(),
			Annotations: meta.GetAnnotations(),
			Created:     createdTime,
		},
	}
}

func (w *DeploymentWrap) populateFields(obj interface{}) {
	objValue := reflect.Indirect(reflect.ValueOf(obj))
	spec := objValue.FieldByName("Spec")
	if !doesFieldExist(spec) {
		log.Errorf("Obj %+v does not have a Spec field", objValue)
		return
	}

	w.populateReplicas(spec)

	var podSpec v1.PodSpec

	switch o := obj.(type) {
	case *openshift_appsv1.DeploymentConfig:
		if o.Spec.Template == nil {
			log.Errorf("Spec obj %+v does not have a Template field or is not a pointer pod spec", spec)
			return
		}
		podSpec = o.Spec.Template.Spec
		// Pods don't have the abstractions that higher level objects have so maintain it's lifecycle independently
	case *v1.Pod:
		// Standalone Pods do not have a PodTemplate, like the other deployment
		// types do. So, we need to directly access the Pod's Spec field,
		// instead of looking for it inside a PodTemplate.
		podSpec = o.Spec
	case *v1beta1.CronJob:
		podSpec = o.Spec.JobTemplate.Spec.Template.Spec
	default:
		podTemplate, ok := spec.FieldByName("Template").Interface().(v1.PodTemplateSpec)
		if !ok {
			log.Errorf("Spec obj %+v does not have a Template field", spec)
			return
		}
		podSpec = podTemplate.Spec
	}

	w.PopulateDeploymentFromPodSpec(podSpec)
}

// PopulateDeploymentFromPodSpec fills in the initialized wrap with data from the passed pod spec
func (w *DeploymentWrap) PopulateDeploymentFromPodSpec(podSpec v1.PodSpec) {
	w.HostNetwork = podSpec.HostNetwork
	w.populateTolerations(podSpec)
	w.populateServiceAccount(podSpec)
	w.populateImagePullSecrets(podSpec)

	w.populateContainers(podSpec)
}

func (w *DeploymentWrap) populateTolerations(podSpec v1.PodSpec) {
	w.Tolerations = make([]*storage.Toleration, 0, len(podSpec.Tolerations))
	for _, toleration := range podSpec.Tolerations {
		w.Tolerations = append(w.Tolerations, &storage.Toleration{
			Key:         toleration.Key,
			Value:       toleration.Value,
			Operator:    k8s.ToRoxTolerationOperator(toleration.Operator),
			TaintEffect: k8s.ToRoxTaintEffect(toleration.Effect),
		})
	}
}

func (w *DeploymentWrap) populateContainers(podSpec v1.PodSpec) {
	w.Deployment.Containers = make([]*storage.Container, 0, len(podSpec.Containers))
	for _, c := range podSpec.Containers {
		w.Deployment.Containers = append(w.Deployment.Containers, &storage.Container{
			Id:   fmt.Sprintf("%s:%s", w.Deployment.Id, c.Name),
			Name: c.Name,
		})
	}

	w.populateContainerConfigs(podSpec)
	w.populateImages(podSpec)
	w.populateSecurityContext(podSpec)
	w.populateVolumesAndSecrets(podSpec)
	w.populatePorts(podSpec)
	w.populateResources(podSpec)
}

func (w *DeploymentWrap) populateServiceAccount(podSpec v1.PodSpec) {
	if podSpec.ServiceAccountName == "" {
		w.ServiceAccount = "default"
	} else {
		w.ServiceAccount = podSpec.ServiceAccountName
	}
}

func (w *DeploymentWrap) populateImagePullSecrets(podSpec v1.PodSpec) {
	secrets := make([]string, 0, len(podSpec.ImagePullSecrets))
	for _, s := range podSpec.ImagePullSecrets {
		secrets = append(secrets, s.Name)
	}
	w.ImagePullSecrets = secrets
}

func (w *DeploymentWrap) populateReplicas(spec reflect.Value) {
	replicaField := spec.FieldByName("Replicas")
	if !doesFieldExist(replicaField) {
		return
	}

	replicasPointer, ok := replicaField.Interface().(*int32)
	if ok && replicasPointer != nil {
		w.Deployment.Replicas = int64(*replicasPointer)
	}

	replicas, ok := replicaField.Interface().(int32)
	if ok {
		w.Deployment.Replicas = int64(replicas)
	}
}

func (w *DeploymentWrap) populateContainerConfigs(podSpec v1.PodSpec) {
	for i, c := range podSpec.Containers {
		config := &storage.ContainerConfig{
			Command:   c.Command,
			Args:      c.Args,
			Directory: c.WorkingDir,
		}

		envSlice := make([]*storage.ContainerConfig_EnvironmentConfig, len(c.Env))
		for i, env := range c.Env {
			if env.ValueFrom == nil {
				envSlice[i] = &storage.ContainerConfig_EnvironmentConfig{
					Key:          env.Name,
					Value:        env.Value,
					EnvVarSource: storage.ContainerConfig_EnvironmentConfig_RAW,
				}
			} else {
				var value string
				var envVarSrc storage.ContainerConfig_EnvironmentConfig_EnvVarSource
				switch {
				case env.ValueFrom.SecretKeyRef != nil:
					envVarSrc = storage.ContainerConfig_EnvironmentConfig_SECRET_KEY
					ref := env.ValueFrom.SecretKeyRef
					value = fmt.Sprintf("Refers to secret %q with key %q", ref.Name, ref.Key)
				case env.ValueFrom.ConfigMapKeyRef != nil:
					envVarSrc = storage.ContainerConfig_EnvironmentConfig_CONFIG_MAP_KEY
					ref := env.ValueFrom.ConfigMapKeyRef
					value = fmt.Sprintf("Refers to config map %q with key %q", ref.Key, ref.Name)
				case env.ValueFrom.FieldRef != nil:
					envVarSrc = storage.ContainerConfig_EnvironmentConfig_FIELD
					ref := env.ValueFrom.FieldRef
					value = fmt.Sprintf("Refers to field %q", ref.FieldPath)
				case env.ValueFrom.ResourceFieldRef != nil:
					envVarSrc = storage.ContainerConfig_EnvironmentConfig_RESOURCE_FIELD
					ref := env.ValueFrom.ResourceFieldRef
					value = fmt.Sprintf("Refers to resource %q from container %q", ref.Resource, ref.ContainerName)
				default:
					envVarSrc = storage.ContainerConfig_EnvironmentConfig_UNKNOWN
					value = "Unknown environment value reference"
				}
				envSlice[i] = &storage.ContainerConfig_EnvironmentConfig{
					Key:          env.Name,
					Value:        value,
					EnvVarSource: envVarSrc,
				}
			}
		}

		config.Env = envSlice

		if s := c.SecurityContext; s != nil {
			if uid := s.RunAsUser; uid != nil {
				config.Uid = *uid
			}
		}

		w.Deployment.Containers[i].Config = config
	}
}

func (w *DeploymentWrap) populateImages(podSpec v1.PodSpec) {
	for i, c := range podSpec.Containers {
		parsedImage, err := imageUtils.GenerateImageFromString(c.Image)
		if err != nil {
			log.Error(err)
			parsedImage = &storage.ContainerImage{
				Name: &storage.ImageName{
					FullName: fmt.Sprintf("%s is an invalid image", c.Image),
				},
			}
		}
		w.Deployment.Containers[i].Image = parsedImage
	}
}

func (w *DeploymentWrap) populateSecurityContext(podSpec v1.PodSpec) {
	for i, c := range podSpec.Containers {
		sc := &storage.SecurityContext{}
		if s := c.SecurityContext; s != nil {
			if p := s.Privileged; p != nil {
				sc.Privileged = *p
			}

			if p := s.ReadOnlyRootFilesystem; p != nil {
				sc.ReadOnlyRootFilesystem = *p
			}

			if SELinux := s.SELinuxOptions; SELinux != nil {
				sc.Selinux = &storage.SecurityContext_SELinux{
					User:  SELinux.User,
					Role:  SELinux.Role,
					Type:  SELinux.Type,
					Level: SELinux.Level,
				}
			}

			if capabilities := s.Capabilities; capabilities != nil {
				for _, add := range capabilities.Add {
					sc.AddCapabilities = append(sc.AddCapabilities, string(add))
				}

				for _, drop := range capabilities.Drop {
					sc.DropCapabilities = append(sc.DropCapabilities, string(drop))
				}
			}

		}
		w.Deployment.Containers[i].SecurityContext = sc
	}
}

func (w *DeploymentWrap) getVolumeSourceMap(podSpec v1.PodSpec) map[string]volumes.VolumeSource {
	volumeSourceMap := make(map[string]volumes.VolumeSource)
	for _, v := range podSpec.Volumes {
		val := reflect.ValueOf(v.VolumeSource)
		for i := 0; i < val.NumField(); i++ {
			f := val.Field(i)
			if !f.IsNil() {
				sourceCreator, ok := volumes.VolumeRegistry[val.Type().Field(i).Name]
				if !ok {
					volumeSourceMap[v.Name] = &volumes.Unimplemented{}
				} else {
					volumeSourceMap[v.Name] = sourceCreator(f.Interface())
				}
			}
		}
	}
	return volumeSourceMap
}

func (w *DeploymentWrap) populateResources(podSpec v1.PodSpec) {
	for i, c := range podSpec.Containers {
		w.Deployment.Containers[i].Resources = &storage.Resources{
			CpuCoresRequest: k8s.ConvertQuantityToCores(c.Resources.Requests.Cpu()),
			CpuCoresLimit:   k8s.ConvertQuantityToCores(c.Resources.Limits.Cpu()),
			MemoryMbRequest: k8s.ConvertQuantityToMB(c.Resources.Requests.Memory()),
			MemoryMbLimit:   k8s.ConvertQuantityToMB(c.Resources.Limits.Memory()),
		}
	}
}

func (w *DeploymentWrap) populateVolumesAndSecrets(podSpec v1.PodSpec) {
	volumeSourceMap := w.getVolumeSourceMap(podSpec)
	for i, c := range podSpec.Containers {
		for _, v := range c.VolumeMounts {
			sourceVolume, ok := volumeSourceMap[v.Name]
			if !ok {
				sourceVolume = &volumes.Unimplemented{}
			}
			if sourceVolume.Type() == "Secret" {
				w.Deployment.Containers[i].Secrets = append(w.Deployment.Containers[i].Secrets, &storage.EmbeddedSecret{
					Name: sourceVolume.Source(),
					Path: v.MountPath,
				})
				continue
			}

			w.Deployment.Containers[i].Volumes = append(w.Deployment.Containers[i].Volumes, &storage.Volume{
				Name:        v.Name,
				Source:      sourceVolume.Source(),
				Destination: v.MountPath,
				ReadOnly:    v.ReadOnly,
				Type:        sourceVolume.Type(),
			})
		}

		for _, s := range podSpec.ImagePullSecrets {
			w.Deployment.Containers[i].Secrets = append(w.Deployment.Containers[i].Secrets, &storage.EmbeddedSecret{
				Name: s.Name,
			})
		}
	}
}

func (w *DeploymentWrap) populatePorts(podSpec v1.PodSpec) {
	w.Ports = nil
	for i, c := range podSpec.Containers {
		for _, p := range c.Ports {
			var exposures []*storage.PortConfig_ExposureInfo
			exposureLevel := storage.PortConfig_UNSET
			if p.HostPort != 0 {
				hostPortExposure := &storage.PortConfig_ExposureInfo{
					Level:    storage.PortConfig_HOST,
					NodePort: p.HostPort,
				}
				exposures = []*storage.PortConfig_ExposureInfo{hostPortExposure}
				exposureLevel = storage.PortConfig_HOST
			}

			protocolStr := string(p.Protocol)
			if protocolStr == "" {
				protocolStr = string(v1.ProtocolTCP)
			}

			portConfig := &storage.PortConfig{
				Name:          p.Name,
				ContainerPort: p.ContainerPort,
				Protocol:      protocolStr,
				Exposure:      exposureLevel,
				ExposureInfos: exposures,
			}
			w.Deployment.Containers[i].Ports = append(w.Deployment.Containers[i].Ports, portConfig)
			w.Ports = append(w.Ports, portConfig)
		}
	}
}
