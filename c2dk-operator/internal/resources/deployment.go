package resources

import (
	c2dkv1 "c2dk-operator/api/v1"
	"context"
	"errors"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type C2DeploymentStatus struct {
	Name      string        `json:"name"`
	NameSpace string        `json:"namespace"`
	Status    bool          `json:"status"`
	PodStatus []C2PodStatus `json:"podstatus"`
}

func removeDuplicateDeployment(deploymentList []appsv1.Deployment) []appsv1.Deployment {
	result := make([]appsv1.Deployment, 0)
	for i := range deploymentList {
		flag := true
		for j := range result {
			if deploymentList[i].Name == result[j].Name && deploymentList[i].Namespace == result[j].Namespace {
				flag = false
				break
			}
			if flag {
				result = append(result, deploymentList[i])
			}
		}
	}
	return result
}

func deploymentValidate(deploymentList []appsv1.Deployment) error {
	for _, deployment := range deploymentList {
		if *deployment.Spec.Replicas >= 10 {
			return errors.New(fmt.Sprintf("namespace: %s deployment: %s replicas is too much", deployment.Namespace, deployment.Name))
		}
	}
	return nil
}

func GenerateDeploymentByC2app(c2app *c2dkv1.C2app) ([]appsv1.Deployment, error) {
	// data validata
	deploymentList := make([]appsv1.Deployment, 0)

	for _, application := range c2app.Spec.ApplicationList {
		application := application
		replicas := application.Replicas // application address not change, can't use it direct
		deployment := appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            application.Name,
				Namespace:       application.NameSpace,
				Labels:          application.Labels,
				OwnerReferences: NewOwnerReference(c2app),
				Annotations:     application.Annotations,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: application.Labels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: application.Labels},
					Spec:       application.PodSpec,
				},
			},
		}
		deploymentList = append(deploymentList, deployment)
	}
	_ = deploymentValidate(deploymentList)
	return deploymentList, nil
}

func CreateDeploymentWithPolicy(cli client.Client, deployment *appsv1.Deployment) error {
	switch UPDATE_POLICY {
	case CREATE:
		return createDeploymentWithNoUpdate(cli, deployment)
	case UPDATE:
		return createOrUpdateDeployment(cli, deployment)
	default:
		return errors.New("update policy error")
	}
}

// create deployment
func createDeploymentWithNoUpdate(cli client.Client, deployment *appsv1.Deployment) error {
	objectKey := client.ObjectKeyFromObject(deployment)
	var oldDeployment appsv1.Deployment

	if err := cli.Get(context.TODO(), objectKey, &oldDeployment); err != nil {
		if client.IgnoreNotFound(err) == nil {
			// not exist
			if err := cli.Create(context.TODO(), deployment); err != nil {
				return err
			} else {
				return nil
			}
		} else {
			return errors.New(fmt.Sprintf("namespace: %s deployment: %s get failed", deployment.Namespace, deployment.Name))
		}
	}
	return nil
}

// create deployment or update deployment
func createOrUpdateDeployment(cli client.Client, deployment *appsv1.Deployment) error {
	objectKey := client.ObjectKeyFromObject(deployment)
	var oldDeployment appsv1.Deployment

	// get old deployment
	err := cli.Get(context.TODO(), objectKey, &oldDeployment)
	if err != nil {
		if err := client.IgnoreNotFound(err); err != nil {
			return err
		} else {
			// deployment not exist
			if err := cli.Create(context.TODO(), deployment); err != nil {
				return err
			} else {
				return nil
			}
		}
	}
	//_, err := controllerutil.CreateOrUpdate(context.TODO(), cli, deployment, func() error {
	//	return nil
	//})
	deployment.ObjectMeta = oldDeployment.ObjectMeta

	if err := cli.Update(context.TODO(), deployment); err != nil {
		return err
	} else {
		return nil
	}
}

// deployment and pod status query
func DeploymentStatusQuery(cli client.Client, objectKey client.ObjectKey) (*C2DeploymentStatus, error) {
	var deployment appsv1.Deployment
	var deploymentStatus C2DeploymentStatus

	// init deployment info
	deploymentStatus.Name = objectKey.Name
	deploymentStatus.NameSpace = objectKey.Namespace
	deploymentStatus.Status = false

	err := cli.Get(context.TODO(), objectKey, &deployment)
	if err != nil {
		// deployment not exist or get failed
		return nil, err
	}

	// pods of deployment status check
	podStatusList, err := PodListStatusQuery(cli, objectKey, deployment.Spec.Selector.MatchLabels)
	if err != nil {
		deploymentStatus.PodStatus = podStatusList
		//deploymentStatus.Status = false
		return &deploymentStatus, err
	} else {
		deploymentStatus.PodStatus = podStatusList
	}

	// deployment status check
	if err := deploymentStatusCheck(&deployment); err != nil {
		//deploymentStatus.Status = false
		return nil, err
	} else {
		deploymentStatus.Status = true
		return &deploymentStatus, nil
	}
}

// deployment status check
func deploymentStatusCheck(deployment *appsv1.Deployment) error {
	if deployment.Status.ReadyReplicas == deployment.Status.Replicas {
		return nil
	} else {
		return errors.New(fmt.Sprintf("namespace/%s  -  deployment/%s Replicas Not Ready", deployment.Namespace, deployment.Name))
	}
}

// just used for test
func initDefaultDeployment(c2app *c2dkv1.C2app) *appsv1.Deployment {
	var replicas int32 = 3

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "c2dk-mysql",
			Namespace: "default",
			Labels: map[string]string{
				"application": "c2dk-mysql",
			},
			OwnerReferences: NewOwnerReference(c2app),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"application": "c2dk-mysql",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"application": "c2dk-mysql",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "c2dk-mysql",
							Image: "mysql:8.0",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3306,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: "c2dk-mysql",
								},
							},
							LivenessProbe: &corev1.Probe{
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
								TimeoutSeconds:      5,
								SuccessThreshold:    1,
								FailureThreshold:    3,
								ProbeHandler: corev1.ProbeHandler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"mysqladmin",
											"-uroot",
											"-p${MYSQL_ROOT_PASSWORD}",
											"ping",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// query deployemnt exist
func DeploymentExistQuery(cli client.Client, objectKey client.ObjectKey) (bool, error) {
	var deployment appsv1.Deployment

	err := cli.Get(context.TODO(), objectKey, &deployment)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			return false, err
		} else {
			return false, err
		}
	}
	return true, nil
}
