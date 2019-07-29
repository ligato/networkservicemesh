package pods

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VppTestCommonPod creates a new vpp-based testing pod
func VppTestCommonPod(app, name, container string, node *v1.Node, env map[string]string, sa string) *v1.Pod {
	envVars := []v1.EnvVar{{Name: "TEST_APPLICATION", Value: app}}
	for k, v := range env {
		envVars = append(envVars,
			v1.EnvVar{
				Name:  k,
				Value: v,
			})
	}

	pod := &v1.Pod{
		ObjectMeta: v12.ObjectMeta{
			Name: name,
		},
		TypeMeta: v12.TypeMeta{
			Kind: "Deployment",
		},
		Spec: v1.PodSpec{
			ServiceAccountName: sa,
			Volumes: []v1.Volume{
				spireVolume(),
			},
			Containers: []v1.Container{
				containerMod(&v1.Container{
					Name:            container,
					Image:           "networkservicemesh/vpp-test-common:latest",
					ImagePullPolicy: v1.PullIfNotPresent,
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							"networkservicemesh.io/socket": resource.NewQuantity(1, resource.DecimalSI).DeepCopy(),
						},
					},
					VolumeMounts: []v1.VolumeMount{
						spireVolumeMount(),
					},
					Env: envVars,
				}),
			},
			TerminationGracePeriodSeconds: &ZeroGraceTimeout,
		},
	}

	if node != nil {
		pod.Spec.NodeSelector = map[string]string{
			"kubernetes.io/hostname": node.Labels["kubernetes.io/hostname"],
		}
	}

	return pod
}
