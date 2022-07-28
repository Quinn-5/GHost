package configstore

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Quinn-5/GHost/ghost/configs/servconf"
	"github.com/Quinn-5/GHost/ghost/resources"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var kubeconfig *string

// Server config with private values
// for use in backend functions.
// Initialize with servconf.New()
type ConfigStore struct {

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

	// rest config for direct connections
	kubeconfig *rest.Config

	// Kubernetes clientset
	clientset *kubernetes.Clientset
}

func (cfg *ConfigStore) setAllValues() error {
	deployment, err := resources.GetDeployment(cfg.Get())
	if err == nil {
		cfg.SetType(deployment.ObjectMeta.Labels["type"])
		cfg.SetCPU(deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String())
		cfg.SetRAM(deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String())
		cfg.SetIP(resources.GetNodeIP(cfg.Get()))

		port, err := resources.GetExternalPort(cfg.Get())
		if err == nil {
			cfg.SetExternalPort(port)
		}

		return err
	} else {
		return fmt.Errorf("deployment %s does not exist", cfg.serverName)
	}
}

func New(username string, serverName string) *ConfigStore {
	cfg := &ConfigStore{}
	var err error

	if kubeconfig == nil {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
	}
	cfg.kubeconfig, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	cfg.clientset, err = kubernetes.NewForConfig(cfg.kubeconfig)
	if err != nil {
		panic(err)
	}

	cfg.setUsername(username)
	cfg.setServerName(serverName)
	cfg.setAllValues()

	return cfg
}

func (cfg *ConfigStore) setUsername(username string) error {
	exp := regexp.MustCompile(`[a-z]([-a-z0-9]*[a-z0-9])?`)
	if !bytes.Equal(exp.Find([]byte(username)), []byte(username)) {
		return errors.New("username must contain only alphanumeric, lowercase characters")
	} else {
		cfg.username = strings.ToLower(username)
		return nil
	}
}

func (cfg *ConfigStore) setServerName(serverName string) error {
	exp := regexp.MustCompile(`[a-z]([-a-z0-9]*[a-z0-9])?`)
	if !bytes.Equal(exp.Find([]byte(serverName)), []byte(serverName)) {
		return errors.New("servername must contain only alphanumeric, lowercase characters")
	} else {
		cfg.serverName = strings.ToLower(serverName)
		return nil
	}
}

func (cfg *ConfigStore) SetType(serverType string) error {
	exp := regexp.MustCompile(`[a-z]([-a-z0-9]*[a-z0-9])?`)
	if !bytes.Equal(exp.Find([]byte(serverType)), []byte(serverType)) {
		return errors.New("servertype must contain only alphanumeric, lowercase characters")
	} else {
		cfg.serverType = strings.ToLower(serverType)
		return nil
	}
}

func (cfg *ConfigStore) SetInternalPort(port int32) {
	cfg.internalPort = port
}

func (cfg *ConfigStore) SetExternalPort(port int32) {
	cfg.externalPort = port
}

func (cfg *ConfigStore) SetIP(ip string) {
	cfg.ip = ip
}

func (cfg *ConfigStore) SetProtocol(protocol apiv1.Protocol) {
	cfg.protocol = protocol
}

func (cfg *ConfigStore) SetCPU(cpu string) {
	if n, err := resource.ParseQuantity(cpu); err == nil {
		cfg.cpu = n
	}
}

func (cfg *ConfigStore) SetRAM(ram string) {
	if n, err := resource.ParseQuantity(ram + "Gi"); err == nil {
		cfg.ram = n
	}
}

func (cfg *ConfigStore) SetDisk(disk string) {
	if n, err := resource.ParseQuantity(disk + "Gi"); err == nil {
		cfg.disk = n
	}
}

func (cfg *ConfigStore) Get() *servconf.ServerConfig {
	return &servconf.ServerConfig{
		Username:     cfg.username,
		ServerName:   cfg.serverName,
		ServerType:   cfg.serverType,
		CPU:          cfg.cpu,
		RAM:          cfg.ram,
		Disk:         cfg.disk,
		IP:           cfg.ip,
		InternalPort: cfg.internalPort,
		ExternalPort: cfg.externalPort,
		Protocol:     cfg.protocol,
		Config:       cfg.kubeconfig,
		Clientset:    cfg.clientset,
	}
}
