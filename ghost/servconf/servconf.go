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

// ServerConfig for use in backend
// Initialize with New()
type ServerConfig struct {

	// Name of user requesting server
	username string

	// Name of new server
	serverName string

	// Type of server requested
	serverType string

	// Number of CPU cores to assign
	cpu resource.Quantity

	// Number of GiB RAM to reserve
	ram resource.Quantity

	// Number of MiB disk space to reserve
	disk resource.Quantity

	// IP address to connect to
	ip string

	// Internal port to be exposed
	internalPort int32

	// External port to connect
	externalPort int32

	// Protocol used for communication
	protocol apiv1.Protocol

	// kubeconfig
	clientset *kubernetes.Clientset
}

type PubServConf struct {
	Username     string
	ServerName   string
	ServerType   string
	CPU          string
	RAM          string
	Disk         string
	IP           string
	InternalPort int32
	ExternalPort int32
}

func New(username string, serverName string) *ServerConfig {
	cfg := &ServerConfig{}

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
		panic(err)
	}
	cfg.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	cfg.setUsername(serverName)
	cfg.setServerName(serverName)

	return cfg
}

func (cfg *ServerConfig) GetUsername() string {
	return cfg.username
}

func (cfg *ServerConfig) setUsername(username string) {
	cfg.username = strings.ToLower(username)
}

func (cfg *ServerConfig) GetServerName() string {
	return cfg.serverName
}

func (cfg *ServerConfig) setServerName(serverName string) {
	cfg.serverName = strings.ToLower(serverName)
}

func (cfg *ServerConfig) GetServerType() string {
	return cfg.serverType
}

func (cfg *ServerConfig) SetType(servertype string) {
	cfg.serverType = strings.ToLower(servertype)
}

func (cfg *ServerConfig) GetInternalPort() int32 {
	return cfg.internalPort
}

func (cfg *ServerConfig) SetInternalPort(port int32) {
	cfg.internalPort = port
}

func (cfg *ServerConfig) GetExternalPort() int32 {
	return cfg.externalPort
}

func (cfg *ServerConfig) SetExternalPort(port int32) {
	cfg.externalPort = port
}

func (cfg *ServerConfig) GetIP() string {
	return cfg.ip
}

func (cfg *ServerConfig) SetIP(ip string) {
	cfg.ip = ip
}

func (cfg *ServerConfig) GetProtocol() apiv1.Protocol {
	return cfg.protocol
}

func (cfg *ServerConfig) SetProtocol(protocol apiv1.Protocol) {
	cfg.protocol = protocol
}

func (cfg *ServerConfig) GetCPU() resource.Quantity {
	return cfg.cpu
}

func (cfg *ServerConfig) SetCPU(cpu string) {
	if n, err := resource.ParseQuantity(cpu); err == nil {
		cfg.cpu = n
	}
}

func (cfg *ServerConfig) GetRAM() resource.Quantity {
	return cfg.ram
}

func (cfg *ServerConfig) SetRAM(ram string) {
	if n, err := resource.ParseQuantity(ram + "Gi"); err == nil {
		cfg.ram = n
	}
}

func (cfg *ServerConfig) GetDisk() resource.Quantity {
	return cfg.disk
}

func (cfg *ServerConfig) SetDisk(disk string) {
	if n, err := resource.ParseQuantity(disk + "Gi"); err == nil {
		cfg.disk = n
	}
}

func (cfg *ServerConfig) GetKubeConfig() *kubernetes.Clientset {
	return cfg.clientset
}

func (cfg *ServerConfig) PubConf() *PubServConf {
	pconf := &PubServConf{
		Username:     cfg.GetUsername(),
		ServerName:   cfg.GetServerName(),
		ServerType:   cfg.GetServerType(),
		CPU:          cfg.GetCPU().OpenAPISchemaFormat(),
		RAM:          cfg.GetRAM().OpenAPISchemaFormat(),
		Disk:         cfg.GetRAM().OpenAPISchemaFormat(),
		IP:           cfg.GetIP(),
		InternalPort: cfg.GetInternalPort(),
		ExternalPort: cfg.GetExternalPort(),
	}
	return pconf
}
