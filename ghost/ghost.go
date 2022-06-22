/*
Copyright 2017 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package ghost

import (
	"errors"

	"github.com/Quinn-5/learning-go/ghost/deployments"
	"github.com/Quinn-5/learning-go/ghost/resources"
	"github.com/Quinn-5/learning-go/ghost/servconf"
	v1 "k8s.io/api/apps/v1"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func Create(config *servconf.ServerConfig) error {
	// This is probably BAD!!! but I'll fix it whenever I figure out best practices
	config.Init()

	var deployment *v1.Deployment

	switch config.Type {
	case "minecraft":
		deployment = deployments.Minecraft(config)
	case "terraria":
		deployment = deployments.Terraria(config)
	default:
		return errors.New("Invalid server type")
	}

	resources.CreateNodeport(config)
	resources.CreatePersistentVolumeClaim(config)
	resources.CreateDeployment(config, deployment)

	return nil
}

func Delete(config *servconf.ServerConfig) error {
	config.Init()

	resources.DeleteNodeport(config)
	resources.DeletePersistentVolumeClaim(config)
	resources.DeleteDeployment(config)

	return nil
}

func GetAddress(user string, server string) *servconf.Address {
	config := &servconf.ServerConfig{
		Username:   user,
		Servername: server,
	}
	address := &servconf.Address{
		Port: resources.GetExternalPort(config),
		IP:   resources.GetNodeIP(config),
	}
	println(address)
	return address
}
