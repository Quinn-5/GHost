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

func ListUserDeployments(config *servconf.ServerConfig) *appsv1.DeploymentList {
	deploymentClient := config.Clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	fmt.Println("Getting Deployments...")
	deploymentList, err := deploymentClient.List(context.TODO(), metav1.ListOptions{LabelSelector: "user=" + config.Username})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Found %d Deployments for %q.\n", len(deploymentList.Items), config.Username)
	}

	// var deployments []*servconf.ServerConfig

	// for _, element := range deploymentList.Items {
	// 	username := element.ObjectMeta.Labels["user"]
	// 	servername := element.Name
	// 	serverType := element.ObjectMeta.Labels["type"]

	// 	conf := servconf.New(username, servername)
	// 	conf.SetType(serverType)
	// 	deployments = append(deployments, conf)
	// }

	return deploymentList
}
