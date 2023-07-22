package resources

import (
	c2cloudv1 "c2dk-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewPersistentVolumeClaimYaml(storage c2cloudv1.C2Storage, namespace string, C2app *c2cloudv1.C2app) corev1.PersistentVolumeClaim {
	resourceStorage, _ := resource.ParseQuantity("1Gi")
	storageClassName := storage.StorageClassName
	return corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "PersistentVolumeClaim",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      storage.PvcName,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(C2app, schema.GroupVersionKind{
					Group:   c2cloudv1.GroupVersion.Group,
					Version: c2cloudv1.GroupVersion.Version,
					Kind:    C2app.Kind,
				}),
			},
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: resourceStorage,
				},
				Limits: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: resourceStorage,
				},
			},
			StorageClassName: &storageClassName,
		},
	}
}
