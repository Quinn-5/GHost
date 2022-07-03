package resources

import (
	"context"

	"github.com/Quinn-5/GHost/ghost/servconf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNodeIP(config *servconf.ServerConfig) {
	nodesClient := config.GetKubeConfig().CoreV1().Nodes()

	result, err := nodesClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	ip := result.Items[0].Status.Addresses[0].Address
	config.SetIP(ip)
}
