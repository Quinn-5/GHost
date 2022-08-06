package deployments

import (
	"github.com/Quinn-5/GHost/ghost/configs/servconf"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Minecraft(config *servconf.ServerConfig) *appsv1.Deployment {
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
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name: config.ServerName,
							Env: []apiv1.EnvVar{
								{
									Name:  "EULA",
									Value: "TRUE",
								},
								{
									Name:  "MEMORY",
									Value: "",
								},
								// Allow JVM heap to use 80% of the container's memory
								{
									Name:  "JVM_XX_OPTS",
									Value: "-XX:MaxRAMPercentage=80",
								},
							},
							Image: "itzg/minecraft-server",
							Stdin: true,
							TTY:   true,
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
									MountPath: "/data",
									Name:      config.ServerName,
								},
							},
						},
					},
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

	config.InternalPort = 25565
	config.Protocol = apiv1.ProtocolTCP

	return deployment
}
