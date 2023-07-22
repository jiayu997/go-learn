package resources

import (
	c2cloudv1 "c2dk-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewSecretYaml(secret c2cloudv1.C2Secret, namespace string, C2app *c2cloudv1.C2app) corev1.Secret {
	return corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.Name,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(C2app, schema.GroupVersionKind{
					Group:   c2cloudv1.GroupVersion.Group,
					Version: c2cloudv1.GroupVersion.Version,
					Kind:    C2app.Kind,
				}),
			},
		},
		Data:       secret.Data,
		StringData: secret.StringData,
	}
}
