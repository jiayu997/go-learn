package resources

import (
	c2dkv1 "c2dk-operator/api/v1"
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func configmapValidate(configmapList []corev1.ConfigMap) error {
	return nil
}

func GenerateConfigMapByC2app(c2app *c2dkv1.C2app) ([]corev1.ConfigMap, error) {
	var configMapList []corev1.ConfigMap = make([]corev1.ConfigMap, 0)
	for _, application := range c2app.Spec.ApplicationList {
		for _, configInfo := range application.ConfigMapSpec {
			var configmap corev1.ConfigMap = corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ConfigMap",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            configInfo.Name,
					Namespace:       application.NameSpace,
					OwnerReferences: NewOwnerReference(c2app),
				},
				Data: configInfo.Data,
			}
			configMapList = append(configMapList, configmap)
		}
	}
	// todo
	_ = configmapValidate(configMapList)
	return configMapList, nil
}

func CreateConfigmapWithPolicy(cli client.Client, configmap *corev1.ConfigMap) error {
	switch UPDATE_POLICY {
	case CREATE:
		return createConfigMapWithNoUpdate(cli, configmap)
	case UPDATE:
		return createOrUpdateConfigMap(cli, configmap)
	default:
		return errors.New("update policy error")
	}
}

func createOrUpdateConfigMap(cli client.Client, configmap *corev1.ConfigMap) error {
	var objectKey client.ObjectKey = client.ObjectKeyFromObject(configmap)
	var oldConfigmap corev1.ConfigMap

	// get old configmap
	err := cli.Get(context.TODO(), objectKey, &oldConfigmap)
	if err != nil {
		if err := client.IgnoreNotFound(err); err != nil {
			return err
		} else {
			// configmap not exist
			if err := cli.Create(context.TODO(), configmap); err != nil {
				return err
			} else {
				return nil
			}
		}
	}

	configmap.ObjectMeta = oldConfigmap.ObjectMeta
	if err := cli.Update(context.TODO(), configmap); err != nil {
		return err
	} else {
		return nil
	}

	//_, err := controllerutil.CreateOrUpdate(context.TODO(), cli, configmap, func() error {
	//	// mutate configmap
	//	//MutateConfigMap(newConfigMap,oldConfigMap)
	//	//return controllerutil.SetOwnerReference(c2app, configmap, r.Scheme)
	//	return nil
	//})

	//if err != nil {
	//	return err
	//}
	//return nil
}

func createConfigMapWithNoUpdate(cli client.Client, configmap *corev1.ConfigMap) error {
	var objectKey client.ObjectKey = client.ObjectKeyFromObject(configmap)
	var oldConfigMap corev1.ConfigMap
	if err := cli.Get(context.TODO(), objectKey, &oldConfigMap); err != nil {
		if client.IgnoreNotFound(err) == nil {
			// not exist need to create it
			if err := cli.Create(context.TODO(), configmap); err != nil {
				return err
			} else {
				return nil
			}
		} else {
			return errors.New(fmt.Sprintf("namespace: %s configmap: %s get failed", objectKey.Namespace, objectKey.Name))
		}
	}
	return nil
}

func initDefaultConfigMap(c2app *c2dkv1.C2app) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            "c2dk",
			Namespace:       "default",
			OwnerReferences: NewOwnerReference(c2app),
		},
		Data: map[string]string{
			"test1": "test1",
			"test2": "test2",
		},
	}
}
