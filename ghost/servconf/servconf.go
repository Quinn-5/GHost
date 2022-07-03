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

var kubeconfig *string

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

	// IP address to connect to
	IP string

	// Internal port to be exposed
	internalPort int32

	// External port to connect
	ExternalPort int32

	// Protocol used for communication
	protocol apiv1.Protocol

	// kubeconfig
	clientset *kubernetes.Clientset
}

// Generates config and cleans inputs
func (cfg *ServerConfig) Init() error {
	if kubeconfig == nil {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}
	cfg.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	cfg.Username = strings.ToLower(cfg.Username)
	cfg.Servername = strings.ToLower(cfg.Servername)
	cfg.Type = strings.ToLower(cfg.Type)
	return nil
}

func (cfg *ServerConfig) GetKubeConfig() *kubernetes.Clientset {
	return cfg.clientset
}

func (cfg *ServerConfig) GetPort() int32 {
	return cfg.internalPort
}

func (cfg *ServerConfig) SetIP(IP string) {
	cfg.IP = IP
}

func (cfg *ServerConfig) SetInternalPort(port int32) {
	cfg.internalPort = port
}

func (cfg *ServerConfig) SetExternalPort(port int32) {
	cfg.ExternalPort = port
}

func (cfg *ServerConfig) GetProtocol() apiv1.Protocol {
	return cfg.protocol
}

func (cfg *ServerConfig) SetProtocol(protocol apiv1.Protocol) {
	cfg.protocol = protocol
}
