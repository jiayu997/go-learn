package ClientSet

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Demo2() {
	// config
	config, _ := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)

	// client set
	clientset, _ := kubernetes.NewForConfig(config)
	watch, _ := clientset.CoreV1().Pods("monitor").Watch(context.TODO(), metav1.ListOptions{})

	for {
		result := <-watch.ResultChan()
		fmt.Println(result)
	}
}
