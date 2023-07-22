package resources

import (
	c2cloudv1 "c2dk-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// 创建service
func NewServiceYaml(app c2cloudv1.C2AppSpec, C2app *c2cloudv1.C2app) corev1.Service {
	return corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			Labels:    app.Labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(C2app, schema.GroupVersionKind{
					Group:   c2cloudv1.GroupVersion.Group,
					Version: c2cloudv1.GroupVersion.Version,
					Kind:    C2app.Kind,
				}),
			},
		},
		Spec: app.ServiceSpec,
	}
}
