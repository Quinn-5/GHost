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
	Username   string
	Servername string
	Type       string
	CPU        resource.Format
	RAM        resource.Format
	Disk       resource.Format
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
									apiv1.ResourceCPU:    resource.Quantity{Format: resource.Format(config.CPU)},
									apiv1.ResourceMemory: resource.Quantity{Format: resource.Format(config.RAM)},
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
