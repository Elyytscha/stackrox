package orchestrator

import (
	"testing"

	"github.com/stackrox/rox/pkg/orchestrators"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCovertDaemonSet(t *testing.T) {
	t.Parallel()

	zeroInt64 := int64(0)

	cases := []struct {
		input    *serviceWrap
		expected *v1beta1.DaemonSet
	}{
		{
			input: &serviceWrap{
				SystemService: orchestrators.SystemService{
					Mounts: []string{"/var/run/docker.sock:/var/run/docker.sock", "/tmp:/var/lib"},
					ExtraPodLabels: map[string]string{
						"extra-pod-label": "blah-daemon",
					},
					Envs:      []string{"hello=world", "foo=bar"},
					Name:      `daemon`,
					Image:     `stackrox/daemon:1.0`,
					RunAsUser: &zeroInt64,
				},
				namespace: "prevent",
			},
			expected: &v1beta1.DaemonSet{
				TypeMeta: metav1.TypeMeta{
					Kind:       `DaemonSet`,
					APIVersion: `extensions/v1beta1`,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      `daemon`,
					Namespace: `prevent`,
					Labels: map[string]string{
						"app": "daemon",
					},
				},
				Spec: v1beta1.DaemonSetSpec{
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: `prevent`,
							Labels: map[string]string{
								"app":             "daemon",
								"extra-pod-label": "blah-daemon",
							},
						},
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name: `daemon`,
									Env: []v1.EnvVar{
										{
											Name:  `hello`,
											Value: `world`,
										},
										{
											Name:  `foo`,
											Value: `bar`,
										},
									},
									Image: `stackrox/daemon:1.0`,
									VolumeMounts: []v1.VolumeMount{
										{
											Name:      `var-run-docker-sock`,
											MountPath: `/var/run/docker.sock`,
										},
										{
											Name:      `tmp`,
											MountPath: `/var/lib`,
										},
									},
									SecurityContext: &v1.SecurityContext{
										RunAsUser: &zeroInt64,
									},
								},
							},
							Tolerations: []v1.Toleration{
								{
									Effect:   v1.TaintEffectNoSchedule,
									Operator: v1.TolerationOpExists,
								},
							},
							RestartPolicy: v1.RestartPolicyAlways,
							Volumes: []v1.Volume{
								{
									Name: `var-run-docker-sock`,
									VolumeSource: v1.VolumeSource{
										HostPath: &v1.HostPathVolumeSource{
											Path: `/var/run/docker.sock`,
										},
									},
								},
								{
									Name: `tmp`,
									VolumeSource: v1.VolumeSource{
										HostPath: &v1.HostPathVolumeSource{
											Path: `/tmp`,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		actual := asDaemonSet(c.input)

		assert.Equal(t, c.expected, actual)
	}

}
