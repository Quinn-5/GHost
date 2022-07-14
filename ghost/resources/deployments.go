package resources

import (
	"context"
	"fmt"

	"github.com/Quinn-5/GHost/ghost/configs/servconf"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateDeployment(config *servconf.ServerConfig, deployment *appsv1.Deployment) {
	deploymentsClient := config.Clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	fmt.Println("Creating Deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created Deployment %q.\n", result.GetObjectMeta().GetName())
}

func DeleteDeployment(config *servconf.ServerConfig) {
	deploymentClient := config.Clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	fmt.Println("Deleting Deployment...")
	err := deploymentClient.Delete(context.TODO(), config.ServerName, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Deleted Deployment %q.\n", config.ServerName)
	}
}

func GetDeployment(config *servconf.ServerConfig) (*appsv1.Deployment, error) {
	deploymentClient := config.Clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	fmt.Printf("Searching for Deployment %s...\n", config.ServerName)
	deployment, err := deploymentClient.Get(context.TODO(), config.ServerName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	} else {
		fmt.Printf("Found Deployment %s.\n", config.ServerName)
		return deployment, nil
	}
}

func ListUserDeployments(config *servconf.ServerConfig) *appsv1.DeploymentList {
	deploymentClient := config.Clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	fmt.Println("Getting Deployments...")
	deploymentList, err := deploymentClient.List(context.TODO(), metav1.ListOptions{LabelSelector: "user=" + config.Username})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Found %d Deployments for %q.\n", len(deploymentList.Items), config.Username)
	}

	return deploymentList
}
