package servconf

import (
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
)

type ServerConfig struct {
	Username     string
	ServerName   string
	ServerType   string
	CPU          resource.Quantity
	RAM          resource.Quantity
	Disk         resource.Quantity
	IP           string
	InternalPort int32
	ExternalPort int32
	Protocol     apiv1.Protocol
	Clientset    *kubernetes.Clientset
}
