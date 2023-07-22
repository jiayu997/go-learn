package resources

import (
	c2dkv1 "c2dk-operator/api/v1"
	"context"
	"errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func SecretStatusQueryByObjectKey() {

}

func GenerateSecretByC2app(c2app *c2dkv1.C2app) []*corev1.Secret {
	// todo: need to do it base on c2app resource
	return []*corev1.Secret{
		initDefaultSecret(c2app),
	}
}

func CreateSecretWithPolicy(cli client.Client, secret *corev1.Secret) error {
	switch UPDATE_POLICY {
	case CREATE:
		return createSecretWithNoUpdate(cli, secret)
	case UPDATE:
		return createOrUpdateSecret(cli, secret)
	default:
		return errors.New("update policy error")
	}
}

func createOrUpdateSecret(cli client.Client, secret *corev1.Secret) error {
	_, err := controllerutil.CreateOrUpdate(context.TODO(), cli, secret, func() error {
		return nil
	})

	if err != nil {
		return err
	} else {
		return nil
	}
}

func createSecretWithNoUpdate(cli client.Client, secret *corev1.Secret) error {
	var objectKey client.ObjectKey = client.ObjectKeyFromObject(secret)

	var oldSecret corev1.Secret

	if err := cli.Get(context.TODO(), objectKey, &oldSecret); err != nil {
		if client.IgnoreNotFound(err) == nil {
			// not exist need to create
			if err := cli.Create(context.TODO(), secret); err != nil {
				return err
			} else {
				return nil
			}
		}
	}
	return nil
}

func initDefaultSecret(c2app *c2dkv1.C2app) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            "c2dk",
			Namespace:       "default",
			OwnerReferences: NewOwnerReference(c2app),
		},
		StringData: map[string]string{
			"test1": "test1",
			"test2": "test2",
		},
	}
}
