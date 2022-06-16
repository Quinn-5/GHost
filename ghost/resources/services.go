package resources

import (
	"context"
	"fmt"

	"github.com/Quinn-5/learning-go/ghost/servconf"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNodeport(config *servconf.ServerConfig, internalPort int32, protocol apiv1.Protocol) {
	servicesClient := config.GetKubeConfig().CoreV1().Services(apiv1.NamespaceDefault)

	nodeport := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.Servername,
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": config.Servername,
			},
			Ports: []apiv1.ServicePort{
				{
					Port:     internalPort,
					Protocol: protocol,
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
	err := servicesClient.Delete(context.TODO(), config.Servername, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Deleted NodePort %q.\n", config.Servername)
	}
}
