package resources

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Creates a PersistentVolumeClaim on the cluster in the default namespace, and with the given parameters
//
// Size should be the number of bytes requested
func CreatePersistentVolumeClaim(clientset *kubernetes.Clientset, name string, size int64) {
	storageClient := clientset.CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault)

	storageClass := "csi-rbd-sc"
	volumeMode := apiv1.PersistentVolumeFilesystem

	pvc := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteOnce,
			},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					"storage": *resource.NewQuantity(size, resource.DecimalSI),
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
	}
	fmt.Printf("Created PersistentVolumeClaim %q.\n", result.GetObjectMeta().GetName())

}
