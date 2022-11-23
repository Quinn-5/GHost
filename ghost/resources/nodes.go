package resources

import (
	"context"

	"github.com/Quinn-5/GHost/ghost/configs/servconf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Returns the IP address of the node currently running the server described in the ServerConfig
func GetNodeIP(config *servconf.ServerConfig) string {
	nodesClient := config.Clientset.CoreV1().Nodes()

	result, err := nodesClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	ip := result.Items[0].Status.Addresses[0].Address
	return ip
}
