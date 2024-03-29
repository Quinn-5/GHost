package deployments

import (
	"github.com/Quinn-5/GHost/ghost/configs/servconf"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
					"app": config.ServerName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": config.ServerName,
					},
				},
				// Pod spec should be the only thing that needs to change between deployments
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name: config.ServerName,

							// The container image for this game server
							Image: "nginx",

							Resources: apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									apiv1.ResourceCPU:    config.CPU,
									apiv1.ResourceMemory: config.RAM,
								},
								Requests: apiv1.ResourceList{
									apiv1.ResourceCPU:    resource.MustParse("0"),
									apiv1.ResourceMemory: resource.MustParse("0"),
								},
							},

							VolumeMounts: []apiv1.VolumeMount{
								{
									// Container's internal mount point for persistent data
									MountPath: "/app/config",
									Name:      config.ServerName,
								},
							},
						},
					},

					// Only necessary if you need persistent storage, but you probably do
					Volumes: []apiv1.Volume{
						{
							Name: config.ServerName,
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: config.ServerName,
								},
							},
						},
					},
				},
			},
		},
	}

	// Set accordingly
	config.InternalPort = 80
	// Game protocol is almost always TCP, but some implementations differ.
	config.Protocol = apiv1.ProtocolTCP

	return deployment
}
