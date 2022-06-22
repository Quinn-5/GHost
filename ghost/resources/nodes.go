package resources

import (
	"context"
	"fmt"

	"github.com/Quinn-5/learning-go/ghost/servconf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNodeIP(config *servconf.ServerConfig) string {
	nodesClient := config.GetKubeConfig().CoreV1().Nodes()

	node := GetPod(config).Spec.NodeName

	result, err := nodesClient.Get(context.TODO(), node, metav1.GetOptions{})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Found NodePort %q.\n", config.Servername)
	}
	ip := result.Status.Addresses[0].Address

	return ip
}
