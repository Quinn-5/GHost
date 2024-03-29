package resources

import (
	"context"
	"fmt"

	"github.com/Quinn-5/GHost/ghost/configs/servconf"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Creates a PersistentVolumeClaim on the cluster in the default namespace, and with the given parameters
//
// Size should be the number of bytes requested
func CreatePersistentVolumeClaim(config *servconf.ServerConfig) error {
	storageClient := config.Clientset.CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault)

	storageClass := "longhorn"
	volumeMode := apiv1.PersistentVolumeFilesystem

	pvc := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.ServerName,
			Labels: map[string]string{
				"user": config.Username,
			},
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteOnce,
			},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					"storage": config.Disk,
				},
			},
			StorageClassName: &storageClass,
			VolumeMode:       &volumeMode,
		},
	}

	fmt.Println("Creating PersistentVolumeClaim...")
	result, err := storageClient.Create(context.TODO(), pvc, metav1.CreateOptions{})
	if err != nil {
		if err.Error() == fmt.Sprintf("persistentvolumeclaims \"%s\" already exists", config.ServerName) {
			return fmt.Errorf("volume claim named %s already exists", config.ServerName)
		}
	} else {
		fmt.Printf("Created PersistentVolumeClaim %q.\n", result.GetObjectMeta().GetName())
	}
	return err
}

func DeletePersistentVolumeClaim(config *servconf.ServerConfig) error {
	storageClient := config.Clientset.CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault)

	fmt.Println("Deleting PersistentVolumeClaim...")
	err := storageClient.Delete(context.TODO(), config.ServerName, metav1.DeleteOptions{})
	if err != nil {
		if err.Error() == fmt.Sprintf("persistentvolumeclaims \"%s\" not found", config.ServerName) {
			return fmt.Errorf("volume claim named %s already exists", config.ServerName)
		}
	} else {
		fmt.Printf("Deleted PersistentVolumeClaim %q.\n", config.ServerName)
	}
	return err
}
