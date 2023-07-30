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

func serviceValidate(serviceList []corev1.Service) error {
	return nil
}

func GenerateServiceByC2app(c2app *c2dkv1.C2app) ([]corev1.Service, error) {
	serviceList := make([]corev1.Service, 0)
	for _, application := range c2app.Spec.ApplicationList {
		application := application
		for serviceType, serviceSpec := range application.ServiceSpec {
			if len(serviceSpec.Ports) > 0 {
				service := corev1.Service{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "ConfigMap",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      serviceSpec.Name,
						Namespace: application.NameSpace,
					},
					Spec: corev1.ServiceSpec{
						Ports:    application.ServiceSpec[serviceType].Ports,
						Selector: application.ServiceSpec[serviceType].Selector,
						Type:     application.ServiceSpec[serviceType].Type,
					},
				}
				serviceList = append(serviceList, service)
			}
		}
	}

	_ = serviceValidate(serviceList)
	return serviceList, nil
}

func CreateServiceWithPolicy(cli client.Client, service *corev1.Service) error {
	switch UPDATE_POLICY {
	case CREATE:
		return createServiceWithNoUpdate(cli, service)
	case UPDATE:
		return createOrUpdateService(cli, service)
	default:
		return errors.New("update policy error")
	}
}

func createOrUpdateService(cli client.Client, service *corev1.Service) error {
	objectKey := client.ObjectKeyFromObject(service)
	var oldService corev1.Service

	// get old service
	err := cli.Get(context.TODO(), objectKey, &oldService)
	if err != nil {
		if err := client.IgnoreNotFound(err); err != nil {
			return err
		} else {
			// service not exist
			if err := cli.Create(context.TODO(), service); err != nil {
				return err
			} else {
				return nil
			}
		}
	}

	service.ObjectMeta = oldService.ObjectMeta
	if err := cli.Update(context.TODO(), service); err != nil {
		return err
	} else {
		return nil
	}
}

func createServiceWithNoUpdate(cli client.Client, service *corev1.Service) error {
	objectKey := client.ObjectKeyFromObject(service)
	var oldService corev1.Service
	if err := cli.Get(context.TODO(), objectKey, &oldService); err != nil {
		if client.IgnoreNotFound(err) == nil {
			// not exist
			if err := cli.Create(context.TODO(), service); err != nil {
				return err
			} else {
				return nil
			}
		} else {
			return errors.New(fmt.Sprintf("namespace: %s service: %s get failed", objectKey.Namespace, objectKey.Name))
		}
	}
	return nil
}

func initDefaultService(c2app *c2dkv1.C2app) *corev1.Service {
	return &corev1.Service{}
}
