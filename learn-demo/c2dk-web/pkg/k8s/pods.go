package k8s

import (
	"context"
	"errors"
	"fmt"
	"sync"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 获取所有异常POD
func singleDeletePods(clientset *kubernetes.Clientset, thread int32) {
	podlist, _ := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "status.phase!=Running,status.phase!=Succeeded,status.phase!=Unknown",
	})
	var wg sync.WaitGroup
	var ch = make(chan bool, thread)

	for _, pod := range podlist.Items {
		wg.Add(1)
		go func(pod v1.Pod) {
			ch <- true
			err := clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
			if err != nil {
				fmt.Printf("命名空间：%s--POD名称：%s 删除失败!\n", pod.Namespace, pod.Name)
			}
			<-ch
			wg.Done()
		}(pod)
	}
	wg.Wait()
}

// 批量删除异常pod
func multiDeleteErrorPods(clientset *kubernetes.Clientset) (string, error) {
	// 获取所有namespaces,批量删除必须要指定namespace
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", nil
	}

	// 删除所有pod
	for _, namespace := range namespaceList.Items {
		err = clientset.CoreV1().Pods(namespace.Name).DeleteCollection(context.TODO(), deleteOptions, metav1.ListOptions{
			FieldSelector: "status.phase!=Running,status.phase!=Succeeded,status.phase!=Unknown",
		})
		if err != nil {
			return namespace.Name + "下pods删除失败", err
		}
	}
	return "", err
}

// 删除重启次数过多pod,k8s目前不支持FieldSelector选择重启次数的，所以只能单个删除
func singleDeleteRestartPods(clientset *kubernetes.Clientset, count int32, thread int32) {
	podlist, _ := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	var wg sync.WaitGroup
	var ch = make(chan bool, thread)
	// 闭包传参，解决变量循环问题
	for _, pod := range podlist.Items {
		wg.Add(1)
		go func(pod v1.Pod) {
			ch <- true
			var restartCount int32
			for index, _ := range pod.Spec.Containers {
				restartCount += pod.Status.ContainerStatuses[index].RestartCount
			}
			if restartCount >= count {
				err := clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, deleteOptions)
				if err != nil {
					fmt.Println(pod.Namespace + "-" + pod.Name + "删除失败")
				}
			}
			<-ch
			wg.Done()
		}(pod)
	}
	wg.Wait()
}

func DeleteErrorPods() (string, error) {
	clientset, err := initClientSet()
	if err != nil {
		return "", errors.New("K8S集群连接失败")
	}
	podlist, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "status.phase!=Running,status.phase!=Succeeded,status.phase!=Unknown",
	})
	if err != nil {
		return "", err
	}
	if len(podlist.Items) == 0 {
		return "", errors.New("集群未发现异常POD")
	}
	return multiDeleteErrorPods(clientset)
}

// 查询所有POD
func GetAllPods() (*v1.PodList, error) {
	clientset, err := initClientSet()
	if err != nil {
		return nil, err
	}
	podlist, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podlist, nil
}
