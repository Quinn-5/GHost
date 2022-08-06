package resources

import (
	"context"
	"fmt"
	"io"

	"github.com/Quinn-5/GHost/ghost/configs/servconf"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func GetPod(config *servconf.ServerConfig) *apiv1.Pod {
	podsClient := config.Clientset.CoreV1().Pods(apiv1.NamespaceDefault)

	fmt.Println("Getting Pod...")
	pods, err := podsClient.List(context.TODO(), metav1.ListOptions{LabelSelector: "app=" + config.ServerName})
	pod := pods.Items[0]
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Found Pod %q.\n", config.ServerName)
	}

	return &pod
}

func ShellPrompt(config *servconf.ServerConfig, stdin io.Reader, stdout io.Writer) error {
	podname := GetPod(config).Name

	req := config.Clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podname).
		Namespace(apiv1.NamespaceDefault).
		SubResource("exec")

	req.VersionedParams(&apiv1.PodExecOptions{
		Command: []string{"bash"},
		Stdin:   true,
		Stdout:  true,
		Stderr:  false,
		TTY:     true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config.Config, "POST", req.URL())
	if err != nil {
		return err
	}

	go exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: nil,
		Tty:    true,
	})

	return nil
}
