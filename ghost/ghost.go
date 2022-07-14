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

	"github.com/Quinn-5/GHost/ghost/configs/configstore"
	"github.com/Quinn-5/GHost/ghost/configs/servconf"
	"github.com/Quinn-5/GHost/ghost/deployments"
	"github.com/Quinn-5/GHost/ghost/resources"
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

// checks if the defined config already exists on the cluster
func exists(cfg *servconf.ServerConfig) bool {
	return true
}

func SetAllValues(cfg *servconf.ServerConfig) error {
	if exists(cfg) {
		return nil
	}
	return nil
}

func Create(config *servconf.ServerConfig) error {
	var deployment *v1.Deployment

	switch config.ServerType {
	case "factorio":
		deployment = deployments.Factorio(config)
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
	resources.DeleteNodeport(config)
	resources.DeletePersistentVolumeClaim(config)
	resources.DeleteDeployment(config)

	return nil
}

func GetAddress(config *configstore.ConfigStore) {
	config.SetIP(resources.GetNodeIP(config.Get()))
	config.SetExternalPort(resources.GetExternalPort(config.Get()))
}

// take a look at this later
func GetAllDeploymentsForUser(config *servconf.ServerConfig) []*servconf.ServerConfig {
	deploymentList := resources.ListUserDeployments(config)
	var deployments []*servconf.ServerConfig
	for _, element := range deploymentList.Items {
		username := element.ObjectMeta.Labels["user"]
		servername := element.Name
		serverType := element.ObjectMeta.Labels["type"]
		store := configstore.New(username, servername)
		store.SetType(serverType)
		deployments = append(deployments, store.Get())
	}
	return deployments
}
