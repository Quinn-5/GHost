package deployments

import (
	"github.com/Quinn-5/GHost/ghost/servconf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func int32Ptr(i int32) *int32 { return &i }

func stdMeta(config *servconf.ServerConfig) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name: config.GetServerName(),
		Labels: map[string]string{
			"user": config.GetUsername(),
			"type": config.GetServerType(),
		},
	}
}
