package resources

import (
	"context"
	"fmt"

	"github.com/Quinn-5/learning-go/ghost/servconf"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateDeployment(config *servconf.ServerConfig, deployment *appsv1.Deployment) {
	deploymentsClient := config.GetKubeConfig().AppsV1().Deployments(apiv1.NamespaceDefault)

	fmt.Println("Creating Deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created Deployment %q.\n", result.GetObjectMeta().GetName())
}

func DeleteDeployment(config *servconf.ServerConfig) {
	deploymentClient := config.GetKubeConfig().AppsV1().Deployments(apiv1.NamespaceDefault)
	fmt.Println("Deleting Deployment...")
	err := deploymentClient.Delete(context.TODO(), config.Servername, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Deleted Deployment %q.\n", config.Servername)
	}
}
