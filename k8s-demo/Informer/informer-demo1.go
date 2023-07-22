package Informer

import (
	"fmt"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func Demo1() {
	// config
	config, _ := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)

	// clientset
	clientset, _ := kubernetes.NewForConfig(config)

	// informer
	//factory := informers.NewSharedInformerFactory(clientset, 0)
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace("monitor"))
	informer := factory.Core().V1().Pods().Informer()

	// add event handler function
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		// 只有添加事件才会触发
		AddFunc: func(obj interface{}) {
			fmt.Println("Add Event")
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("Update Event")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Delete Event")
		},
	})

	// start informer
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)
	<-stopCh
}
