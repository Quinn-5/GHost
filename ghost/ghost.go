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
	"context"
	"fmt"

	"github.com/Quinn-5/learning-go/ghost/deployments"
	"github.com/Quinn-5/learning-go/ghost/resources"
	"github.com/Quinn-5/learning-go/ghost/servconf"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	config.Init()

	deployments.Deploy(config)
	resources.CreateNodeport(config)
	resources.CreatePersistentVolumeClaim(config)

	return nil
}

func Delete(config *servconf.ServerConfig) error {
	config.Init()

	resources.DeleteNodeport(config)
	resources.DeletePersistentVolumeClaim(config)

	deploymentClient := config.GetKubeConfig().AppsV1().Deployments(apiv1.NamespaceDefault)
	fmt.Println("Deleting Deployment...")
	err := deploymentClient.Delete(context.TODO(), config.Servername, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Deleted Deployment %q.\n", config.Servername)
	}

	return err
}
