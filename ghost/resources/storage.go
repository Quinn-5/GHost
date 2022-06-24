package resources

import (
	"context"
	"fmt"

	"github.com/Quinn-5/GHost/ghost/servconf"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Creates a PersistentVolumeClaim on the cluster in the default namespace, and with the given parameters
//
// Size should be the number of bytes requested
func CreatePersistentVolumeClaim(config *servconf.ServerConfig) {
	storageClient := config.GetKubeConfig().CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault)

	storageClass := "csi-rbd-sc"
	volumeMode := apiv1.PersistentVolumeFilesystem

	pvc := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.Servername,
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
		panic(err)
	} else {
		fmt.Printf("Created PersistentVolumeClaim %q.\n", result.GetObjectMeta().GetName())
	}
}

func DeletePersistentVolumeClaim(config *servconf.ServerConfig) {
	storageClient := config.GetKubeConfig().CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault)

	fmt.Println("Deleting PersistentVolumeClaim...")
	err := storageClient.Delete(context.TODO(), config.Servername, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Deleted PersistentVolumeClaim %q.\n", config.Servername)
	}
}
