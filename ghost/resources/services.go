package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Quinn-5/GHost/ghost/configs/servconf"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNodeport(config *servconf.ServerConfig) error {
	servicesClient := config.Clientset.CoreV1().Services(apiv1.NamespaceDefault)

	nodeport := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.ServerName,
			Labels: map[string]string{
				"user": config.Username,
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": config.ServerName,
			},
			Ports: []apiv1.ServicePort{
				{
					Port:     config.InternalPort,
					Protocol: config.Protocol,
				},
			},
			Type: apiv1.ServiceTypeNodePort,
		},
	}

	fmt.Println("Creating NodePort...")
	result, err := servicesClient.Create(context.TODO(), nodeport, metav1.CreateOptions{})
	if err != nil {
		if err.Error() == fmt.Sprintf("services \"%s\" already exists", config.ServerName) {
			return errors.New(fmt.Sprintf("nodeport named %s already exists.", config.ServerName))
		}
		panic(err)
	} else {
		fmt.Printf("Created NodePort %q.\n", result.GetObjectMeta().GetName())
		return err
	}

}

func DeleteNodeport(config *servconf.ServerConfig) error {
	servicesClient := config.Clientset.CoreV1().Services(apiv1.NamespaceDefault)

	fmt.Println("Deleting NodePort...")
	err := servicesClient.Delete(context.TODO(), config.ServerName, metav1.DeleteOptions{})
	if err != nil {
		if err.Error() == fmt.Sprintf("services \"%s\" not found", config.ServerName) {
			return errors.New(fmt.Sprintf("nodeport named %s doesn't exist.", config.ServerName))
		}
		panic(err)
	} else {
		fmt.Printf("Deleted NodePort %q.\n", config.ServerName)
		return err
	}
}

func GetExternalPort(config *servconf.ServerConfig) (int32, error) {
	servicesClient := config.Clientset.CoreV1().Services(apiv1.NamespaceDefault)

	fmt.Println("Getting NodePort...")
	result, err := servicesClient.Get(context.TODO(), config.ServerName, metav1.GetOptions{})
	if err != nil {
		if err.Error() == fmt.Sprintf("services \"%s\" not found", config.ServerName) {
			return 1, errors.New(fmt.Sprintf("nodeport named %s doesn't exist.", config.ServerName))
		}
		panic(err)
	} else {
		fmt.Printf("Found NodePort %q.\n", config.ServerName)
		port := result.Spec.Ports[0].NodePort
		return port, err
	}
}
