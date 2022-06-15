package deployments

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Quinn-5/learning-go/ghost/resources"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type ServerConfig struct {

	// Name of user requesting server
	Username string

	// Name of new server
	Servername string

	// Type of server requested
	Type string

	// Number of CPU cores to assign
	CPU resource.Quantity

	// Number of GiB RAM to reserve
	RAM resource.Quantity

	// Number of MiB disk space to reserve
	Disk resource.Quantity

	kubeconfig *kubernetes.Clientset
}

// Cleans inputs to be used for kube api request
func (cfg *ServerConfig) clean() error {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}
	cfg.kubeconfig, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	cfg.Username = strings.ToLower(cfg.Username)
	cfg.Servername = strings.ToLower(cfg.Servername)
	cfg.Type = strings.ToLower(cfg.Type)
	return nil
}

func (cfg *ServerConfig) Create() error {
	cfg.clean()

	Deploy(cfg.kubeconfig, cfg)

	return nil
}

func (cfg *ServerConfig) Delete() error {
	cfg.clean()

	resources.DeleteNodeport(cfg.kubeconfig, cfg.Servername)
	resources.DeletePersistentVolumeClaim(cfg.kubeconfig, cfg.Servername)

	deploymentClient := cfg.kubeconfig.AppsV1().Deployments(apiv1.NamespaceDefault)
	fmt.Println("Deleting Deployment...")
	err := deploymentClient.Delete(context.TODO(), cfg.Servername, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Deleted Deployment %q.\n", cfg.Servername)
	}

	return err
}
