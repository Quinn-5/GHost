package ghost

import (
	"errors"
	"io"

	"github.com/Quinn-5/GHost/ghost/configs/configstore"
	"github.com/Quinn-5/GHost/ghost/configs/servconf"
	"github.com/Quinn-5/GHost/ghost/deployments"
	"github.com/Quinn-5/GHost/ghost/resources"
	v1 "k8s.io/api/apps/v1"
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
		resources.DeleteNodeport(config)
		return err
	}
	err = resources.CreateDeployment(config, deployment)
	if err != nil {
		resources.DeleteNodeport(config)
		resources.DeletePersistentVolumeClaim(config)
		return err
	}

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

func NewTerminal(config *servconf.ServerConfig) (io.Reader, io.Writer) {
	podIn, dataIn := io.Pipe()
	dataOut, podOut := io.Pipe()
	resources.ShellPrompt(config, podIn, podOut)
	return dataOut, dataIn
}
