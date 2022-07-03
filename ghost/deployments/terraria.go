package deployments

import (
	"github.com/Quinn-5/GHost/ghost/servconf"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// An empty example of a functional deployment. For use in testing or creating new deployments
func Terraria(config *servconf.ServerConfig) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: stdMeta(config),
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": config.Servername,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": config.Servername,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  config.Servername,
							Image: "ryshe/terraria",
							Args: []string{
								"-world",
								"/root/.local/share/Terraria/Worlds/" + config.Servername + ".wld",
								"-autocreate",
								"2",
							},
							Stdin: true,
							TTY:   true,
							Resources: apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									apiv1.ResourceCPU:    config.CPU,
									apiv1.ResourceMemory: config.RAM,
								},
							},

							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/root/.local/share/Terraria/Worlds",
									Name:      config.Servername,
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: config.Servername,
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: config.Servername,
								},
							},
						},
					},
				},
			},
		},
	}

	config.SetInternalPort(7777)
	config.SetProtocol(apiv1.ProtocolTCP)

	return deployment
}
