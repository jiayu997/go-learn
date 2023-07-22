package resources

import (
	c2cloudv1 "c2dk-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// 创建deployemnt
func NewDeploymentYaml(app c2cloudv1.C2AppSpec, C2app *c2cloudv1.C2app) appsv1.Deployment {
	return appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
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
		Spec: appsv1.DeploymentSpec{
			Replicas: &app.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: app.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: app.Labels},
				Spec:       app.PodSpec,
			},
		},
	}
}
