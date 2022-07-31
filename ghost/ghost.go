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
	"os"

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
		return errors.New("invalid server type")
	}

	err := resources.CreateNodeport(config)
	if err != nil {
		return err
	}
	err = resources.CreatePersistentVolumeClaim(config)
	if err != nil {
		return err
	}
	err = resources.CreateDeployment(config, deployment)

	return err
}

func Delete(config *servconf.ServerConfig) error {
	err := resources.DeleteNodeport(config)
	if err != nil {
		return err
	}
	err = resources.DeletePersistentVolumeClaim(config)
	if err != nil {
		return err
	}
	err = resources.DeleteDeployment(config)

	return err
}

// take a look at this later
func GetAllDeploymentsForUser(config *servconf.ServerConfig) []*servconf.ServerConfig {
	deploymentList, _ := resources.ListUserDeployments(config)
	var deployments []*servconf.ServerConfig
	for _, element := range deploymentList.Items {
		username := element.ObjectMeta.Labels["user"]
		servername := element.Name

		store := configstore.New(username, servername)
		deployments = append(deployments, store.Get())
	}
	return deployments
}

func NewTerminal(config *servconf.ServerConfig) {
	stdin := os.Stdin
	stdout := os.Stdout
	resources.Exec(config, stdin, stdout)
}
