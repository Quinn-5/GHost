package deployments

import (
	"github.com/Quinn-5/GHost/ghost/servconf"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// An empty example of a functional deployment. For use in testing or creating new deployments
func EmptyDeployment(config *servconf.ServerConfig) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: stdMeta(config),
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": config.GetServerName(),
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": config.GetServerName(),
					},
				},
				// Pod spec should be the only thing that needs to change between deployments
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name: config.GetServerName(),

							// The container image for this game server
							Image: "nginx",

							Resources: apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									apiv1.ResourceCPU:    config.GetCPU(),
									apiv1.ResourceMemory: config.GetRAM(),
								},
							},

							VolumeMounts: []apiv1.VolumeMount{
								{
									// Container's internal mount point for persistent data
									MountPath: "/app/config",
									Name:      config.GetServerName(),
								},
							},
						},
					},

					// Only necessary if you need persistent storage, but you probably do
					Volumes: []apiv1.Volume{
						{
							Name: config.GetServerName(),
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: config.GetServerName(),
								},
							},
						},
					},
				},
			},
		},
	}

	// Set accordingly
	config.SetInternalPort(80)
	// Game protocol is almost always TCP, but some implementations differ.
	config.SetProtocol(apiv1.ProtocolTCP)

	return deployment
}
