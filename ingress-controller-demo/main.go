package main

import (
	"ingress-controller-demo/pkg"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"time"
)

func main() {
	// 1. config
	// 2. client
	// 3. informer
	// 4. add event handler
	// 5. informer start

	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalln(err)
		}
		config = inClusterConfig
	}

	// client set
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	factory := informers.NewSharedInformerFactory(clientset, time.Second*15)
	serviceInformer := factory.Core().V1().Services()
	ingressInformer := factory.Networking().V1().Ingresses()

	// 这里才生产并加informer加到factory中去了
	controller := pkg.NewController(clientset, serviceInformer, ingressInformer)

	// 启动informer
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	// 启动controller,来处理事件
	controller.Run(stopCh)
}
