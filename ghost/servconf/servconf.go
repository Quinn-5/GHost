package servconf

import (
	"flag"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"
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
func (cfg *ServerConfig) Init() error {
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

func (cfg *ServerConfig) GetKubeConfig() *kubernetes.Clientset {
	return cfg.kubeconfig
}
