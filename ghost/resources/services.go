package resources

import (
	"context"
	"fmt"

	"github.com/Quinn-5/GHost/ghost/servconf"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNodeport(config *servconf.ServerConfig) {
	servicesClient := config.GetKubeConfig().CoreV1().Services(apiv1.NamespaceDefault)

	nodeport := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.GetServerName(),
			Labels: map[string]string{
				"user": config.GetUsername(),
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": config.GetServerName(),
			},
			Ports: []apiv1.ServicePort{
				{
					Port:     config.GetInternalPort(),
					Protocol: config.GetProtocol(),
				},
			},
			Type: apiv1.ServiceTypeNodePort,
		},
	}

	fmt.Println("Creating NodePort...")
	result, err := servicesClient.Create(context.TODO(), nodeport, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Created NodePort %q.\n", result.GetObjectMeta().GetName())
	}

}

func DeleteNodeport(config *servconf.ServerConfig) {
	servicesClient := config.GetKubeConfig().CoreV1().Services(apiv1.NamespaceDefault)

	fmt.Println("Deleting NodePort...")
	err := servicesClient.Delete(context.TODO(), config.GetServerName(), metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Deleted NodePort %q.\n", config.GetServerName())
	}
}

func GetExternalPort(config *servconf.ServerConfig) {
	servicesClient := config.GetKubeConfig().CoreV1().Services(apiv1.NamespaceDefault)

	fmt.Println("Getting NodePort...")
	result, err := servicesClient.Get(context.TODO(), config.GetServerName(), metav1.GetOptions{})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Found NodePort %q.\n", config.GetServerName())
	}
	port := result.Spec.Ports[0].NodePort
	config.SetExternalPort(port)
}
