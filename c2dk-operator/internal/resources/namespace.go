package resources

import (
	c2dkv1 "c2dk-operator/api/v1"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func namespaceValidate(namespaceList []corev1.Namespace) error {
	return nil
}

func GenerateNamespaceByC2app(c2app *c2dkv1.C2app) ([]corev1.Namespace, error) {
	namespaceList := make([]corev1.Namespace, 0)

	for _, application := range c2app.Spec.ApplicationList {
		application := application
		namespace := corev1.Namespace{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Namespace",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:   application.NameSpace,
				Labels: application.Labels,
				//OwnerReferences: NewOwnerReference(c2app), // don't use this
			},
		}
		namespaceList = append(namespaceList, namespace)
	}
	_ = namespaceValidate(namespaceList)
	return namespaceList, nil
}

func CreateNamespaceWithNoPolicy(cli client.Client, namespace *corev1.Namespace) error {
	objectKey := client.ObjectKeyFromObject(namespace)
	var oldNamespace corev1.Namespace

	err := cli.Get(context.TODO(), objectKey, &oldNamespace)
	if err != nil {
		if err := client.IgnoreNotFound(err); err != nil {
			return err
		} else {
			// namespace not exist
			if err := cli.Create(context.TODO(), namespace); err != nil {
				return err
			} else {
				return nil
			}
		}
	}

	// namespace exist
	return nil
}

func NamespaceExist(cli client.Client, namespace *corev1.Namespace) bool {
	return true
}
