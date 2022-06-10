package deployments

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/Quinn-5/learning-go/ghost/resources"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type ServerConfig struct {
	// Name of user requesting server
	Username string

	// Name of new server
	Servername string

	// Type of server requested
	Type string

	// Number of CPU cores to assign
	CPU int64

	// Number of GiB RAM to reserve
	RAM int64

	// Number of MiB disk space to reserve
	Disk int64
}

func (cfg *ServerConfig) Create() error {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	cfg.Disk = cfg.Disk * 1024 * 1024
	cfg.RAM = cfg.RAM * 1024 * 1024 * 1024

	Deploy(clientset, cfg)

	return nil
}

func Deploy(clientset *kubernetes.Clientset, config *ServerConfig) {

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.Servername,
			Labels: map[string]string{
				"user": config.Username,
			},
		},
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
							Name: config.Servername,
							Env: []apiv1.EnvVar{
								{
									Name:  "EULA",
									Value: "TRUE",
								},
							},
							Image: "itzg/minecraft-server",
							Resources: apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									apiv1.ResourceCPU:    *resource.NewQuantity(config.CPU, resource.DecimalSI),
									apiv1.ResourceMemory: *resource.NewQuantity(config.RAM, resource.DecimalSI),
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/data",
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

	resources.CreateNodeport(clientset, config.Servername, 25565, apiv1.ProtocolTCP)
	resources.CreatePersistentVolumeClaim(clientset, config.Servername, config.Disk)

	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func int32Ptr(i int32) *int32 { return &i }
