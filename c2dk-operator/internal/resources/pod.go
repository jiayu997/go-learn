package resources

import (
	"context"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type C2PodStatus struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    bool   `json:"status"`
	Reason    string `json:"reason"`
}

// query all pods status of deployment
func PodListStatusQuery(cli client.Client, objectKey client.ObjectKey, label map[string]string) ([]C2PodStatus, error) {
	var podStatusList []C2PodStatus
	var podStatusFlag bool = true

	// get all pods by deployment
	podList := &corev1.PodList{}

	// method 1
	podLabel := labels.SelectorFromSet(label)
	err := cli.List(context.TODO(), podList, &client.ListOptions{
		Namespace:     objectKey.Namespace,
		LabelSelector: podLabel,
	})
	if err != nil {
		return nil, err
	}

	// pod replicas = 0
	if len(podList.Items) == 0 {
		return nil, errors.New(fmt.Sprintf("namespace: %s deployment: %s no avaliable replicas", objectKey.Namespace, objectKey.Name))
	}

	// pod status check
	for _, pod := range podList.Items {
		// pod not health(default health)
		var podStatus C2PodStatus = C2PodStatus{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    true,
			Reason:    "success",
		}
		if err := podStatusPhaseQuery(&pod); err != nil {
			podStatusFlag = false
			podStatus.Status = false
			podStatus.Reason = err.Error()
			podStatusList = append(podStatusList, podStatus)
		} else {
			podStatusList = append(podStatusList, podStatus)
		}
	}

	if podStatusFlag == true {
		return podStatusList, nil
	} else {
		return podStatusList, errors.New("pod status failed")
	}
}

func PodStatusQueryByObjectKey(cli client.Client, name, namespace string) (*C2PodStatus, error) {
	var objectKey = client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}

	var podStatus C2PodStatus = C2PodStatus{
		Name:      name,
		Namespace: namespace,
		Status:    true,
		Reason:    "success",
	}

	pod := &corev1.Pod{}
	if err := cli.Get(context.TODO(), objectKey, pod); err != nil {
		return nil, err
	}

	if err := podStatusPhaseQuery(pod); err != nil {
		podStatus.Status = false
		podStatus.Reason = err.Error()
		return &podStatus, err
	}

	if podStatus.Status != true {
		return &podStatus, errors.New(fmt.Sprintf("namespace/%s -- pod/%s not health", podStatus.Namespace, podStatus.Name))
	} else {
		return &podStatus, nil
	}
}

// pod status phase query
func podStatusPhaseQuery(pod *corev1.Pod) error {
	switch pod.Status.Phase {
	case corev1.PodPending:
		return errors.New(fmt.Sprintf("pod/%s is Pending", pod.Name))
	case corev1.PodFailed:
		return errors.New(fmt.Sprintf("pod/%s is Failed", pod.Name))
	case corev1.PodRunning:
		return nil
	case corev1.PodSucceeded:
		return nil
	case corev1.PodUnknown:
		return errors.New(fmt.Sprintf("pod/%s is Unknown", pod.Name))
	default:
		return nil
	}
}
