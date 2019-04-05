package resources

import (
	"testing"
	"time"

	timestamp "github.com/gogo/protobuf/types"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/kubernetes"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestConvertDifferentContainerNumbers(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name               string
		inputObj           interface{}
		deploymentType     string
		action             central.ResourceAction
		podLister          *mockPodLister
		expectedDeployment *storage.Deployment
	}{
		{
			name: "Deployment",
			inputObj: &v1beta1.Deployment{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1beta1",
					Kind:       "Deployment",
				},
				ObjectMeta: metav1.ObjectMeta{
					UID:               types.UID("FooID"),
					Name:              "deployment",
					Namespace:         "namespace",
					CreationTimestamp: metav1.NewTime(time.Unix(1000, 0)),
				},
				Spec: v1beta1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: make(map[string]string),
					},
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  "container1",
									Image: "docker.io/stackrox/kafka:latest",
								},
							},
						},
					},
				},
			},
			deploymentType: kubernetes.Deployment,
			action:         central.ResourceAction_UPDATE_RESOURCE,
			podLister: &mockPodLister{
				pods: []*v1.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							UID:       types.UID("ebf487f0-a7c3-11e8-8600-42010a8a0066"),
							Name:      "deployment-blah-blah",
							Namespace: "myns",
							OwnerReferences: []metav1.OwnerReference{
								{
									Kind: kubernetes.Deployment,
								},
							},
						},
						Status: v1.PodStatus{
							ContainerStatuses: []v1.ContainerStatus{
								{
									Name:  "container1",
									Image: "docker.io/stackrox/kafka:1.3",
								},
								{
									Name:        "container2",
									Image:       "docker.io/stackrox/policy-engine:1.3",
									ImageID:     "docker-pullable://docker.io/stackrox/policy-engine@sha256:6b561c3bb9fed1b028520cce3852e6c9a6a91161df9b92ca0c3a20ebecc0581a",
									ContainerID: "docker://35669191c32a9cfb532e5d79b09f2b0926c0faf27e7543f1fbe433bd94ae78d7",
								},
							},
						},
						Spec: v1.PodSpec{
							NodeName: "mynode",
							Containers: []v1.Container{
								{
									Name:  "container1",
									Image: "docker.io/stackrox/kafka:latest",
								},
								{
									Name:  "container2",
									Image: "docker.io/stackrox/policy-engine:1.3",
								},
							},
						},
					},
				},
			},
			expectedDeployment: &storage.Deployment{
				Id:          "FooID",
				Name:        "deployment",
				Namespace:   "namespace",
				NamespaceId: "FAKENSID",
				Type:        kubernetes.Deployment,
				LabelSelector: &storage.LabelSelector{
					MatchLabels: map[string]string{},
				},
				UpdatedAt:                    &timestamp.Timestamp{Seconds: 1000},
				ImagePullSecrets:             []string{},
				Tolerations:                  []*storage.Toleration{},
				AutomountServiceAccountToken: true,
				Containers: []*storage.Container{
					{
						Id:   "FooID:container1",
						Name: "container1",
						Image: &storage.Image{
							Name: &storage.ImageName{
								Registry: "docker.io",
								Remote:   "stackrox/kafka",
								Tag:      "latest",
								FullName: "docker.io/stackrox/kafka:latest",
							},
						},
						Config: &storage.ContainerConfig{
							Env: []*storage.ContainerConfig_EnvironmentConfig{},
						},
						SecurityContext: &storage.SecurityContext{},
						Resources:       &storage.Resources{},
						Instances: []*storage.ContainerInstance{
							{
								InstanceId: &storage.ContainerInstanceID{
									Node: "mynode",
								},
								ContainingPodId: "deployment-blah-blah.myns@ebf487f0-a7c3-11e8-8600-42010a8a0066",
							},
						},
					},
					{
						Id:   "FooID:container2",
						Name: "container2",
						Image: &storage.Image{
							Id: "sha256:6b561c3bb9fed1b028520cce3852e6c9a6a91161df9b92ca0c3a20ebecc0581a",
							Name: &storage.ImageName{
								Registry: "docker.io",
								Remote:   "stackrox/policy-engine",
								Tag:      "1.3",
								FullName: "docker.io/stackrox/policy-engine:1.3",
							},
						},
						Config: &storage.ContainerConfig{
							Env: []*storage.ContainerConfig_EnvironmentConfig{},
						},
						SecurityContext: &storage.SecurityContext{},
						Resources:       &storage.Resources{},
						Instances: []*storage.ContainerInstance{
							{
								InstanceId: &storage.ContainerInstanceID{
									ContainerRuntime: storage.ContainerRuntime_DOCKER_CONTAINER_RUNTIME,
									Id:               "35669191c32a9cfb532e5d79b09f2b0926c0faf27e7543f1fbe433bd94ae78d7",
									Node:             "mynode",
								},
								ContainingPodId: "deployment-blah-blah.myns@ebf487f0-a7c3-11e8-8600-42010a8a0066",
							},
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := newDeploymentEventFromResource(c.inputObj, c.action, c.deploymentType, c.podLister, mockNamespaceStore).GetDeployment()
			assert.Equal(t, c.expectedDeployment, actual)
		})
	}
}
