package resources

import (
	c2dkv1 "c2dk-operator/api/v1"
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func pvcValidate(pvcList []corev1.PersistentVolumeClaim) error {
	return nil
}

func GenerateStorageByC2app(c2app *c2dkv1.C2app) ([]corev1.PersistentVolumeClaim, error) {
	var pvcList []corev1.PersistentVolumeClaim = make([]corev1.PersistentVolumeClaim, 0)
	defaultRequest, _ := resource.ParseQuantity("1Gi")
	defaultLimit, _ := resource.ParseQuantity("10Gi")
	defaultStorageClass := STORAGE_CLASS_NFS

	for _, application := range c2app.Spec.ApplicationList {
		for _, storage := range application.StorageSpec {
			var pvc corev1.PersistentVolumeClaim = corev1.PersistentVolumeClaim{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "PersistentVolumeClaim",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      storage.PvcName,
					Namespace: application.NameSpace,
					Labels:    application.Labels,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.PersistentVolumeAccessMode(storage.AccessMode)},
					Resources: corev1.ResourceRequirements{
						Requests: map[corev1.ResourceName]resource.Quantity{
							corev1.ResourceStorage: defaultRequest,
						},
						Limits: map[corev1.ResourceName]resource.Quantity{
							corev1.ResourceStorage: defaultLimit,
						},
					},
					StorageClassName: &defaultStorageClass,
				},
			}
			pvcList = append(pvcList, pvc)
		}
	}

	return pvcList, nil
}

func CreatePvcWithNoPolicy(cli client.Client, pvc *corev1.PersistentVolumeClaim) error {
	var objectKey client.ObjectKey = client.ObjectKeyFromObject(pvc)
	var oldPvc corev1.PersistentVolumeClaim
	if err := cli.Get(context.TODO(), objectKey, &oldPvc); err != nil {
		if client.IgnoreNotFound(err) == nil {
			// not exist
			if err := cli.Create(context.TODO(), pvc); err != nil {
				return err
			} else {
				return nil
			}
		} else {
			return errors.New(fmt.Sprintf("namespace: %s pvc: %s get failed", pvc.Namespace, pvc.Name))
		}
	}

	// exist, we can't change it
	return nil
}
