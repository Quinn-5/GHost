package servconf

import (
	"flag"
	"path/filepath"
	"strings"

	apiv1 "k8s.io/api/core/v1"
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

	port int32

	protocol apiv1.Protocol

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

func (cfg *ServerConfig) SetPort(port int32) {
	cfg.port = port
}

func (cfg *ServerConfig) GetPort() int32 {
	return cfg.port
}

func (cfg *ServerConfig) SetProtocol(protocol apiv1.Protocol) {
	cfg.protocol = protocol
}

func (cfg *ServerConfig) GetProtocol() apiv1.Protocol {
	return cfg.protocol
}
