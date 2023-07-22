package Project2

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func DemoMain() {
	// config
	config, _ := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)

	// clientset
	clientset, _ := kubernetes.NewForConfig(config)

	// 实际返回实现改接口的结构体为：sharedInformerFactory
	factory := informers.NewSharedInformerFactory(clientset, 0)

	// 返回一个serviceinformer struct,kubernetes/staging/src/k8s.io/client-go/informers/core/v1/service.go
	servicesInformer := factory.Core().V1().Services()
	ingressesInformer := factory.Networking().V1().Ingresses()

	controller := NewController(clientset, servicesInformer, ingressesInformer)

	stopCh := make(chan struct{})
	factory.Start(stopCh)
	// 等待cache同步完成
	factory.WaitForCacheSync(stopCh)

	controller.Run(stopCh)
	<-stopCh

}
