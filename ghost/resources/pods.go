package resources

import (
	"context"
	"fmt"

	"github.com/Quinn-5/GHost/ghost/servconf"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPod(config *servconf.ServerConfig) *apiv1.Pod {
	podsClient := config.GetKubeConfig().CoreV1().Pods(apiv1.NamespaceDefault)

	fmt.Println("Getting Pod...")
	pod, err := podsClient.Get(context.TODO(), config.GetServerName(), metav1.GetOptions{})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Found Pod %q.\n", config.GetServerName())
	}

	return pod
}
