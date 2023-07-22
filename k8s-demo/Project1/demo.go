package Project1

//
//import (
//	"k8s.io/client-go/informers"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/client-go/rest"
//	"k8s.io/client-go/tools/clientcmd"
//	"log"
//)
//
//func DemoMain() {
//	// 1. config
//	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
//	if err != nil {
//		inClusterConfig, err := rest.InClusterConfig()
//		if err != nil {
//			log.Fatalln("can't get config")
//		}
//		config = inClusterConfig
//	}
//
//	// client set
//	clientset, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		log.Fatalln("can't create client")
//	}
//
//	// factory
//	factory := informers.NewSharedInformerFactory(clientset, 0)
//	serviceInformer := factory.Core().V1().Services()
//	ingressInformer := factory.Networking().V1().Ingresses()
//
//	controller := NewController(clientset, serviceInformer, ingressInformer)
//
//	stopCh := make(chan struct{})
//	factory.Start(stopCh)
//
//	// 二者通同时关闭
//	factory.WaitForCacheSync(stopCh)
//	controller.Run(stopCh)
//}
//
