package restClient

//
//import (
//	"context"
//	"fmt"
//	v1 "k8s.io/api/core/v1"
//	"k8s.io/client-go/kubernetes/scheme"
//	"k8s.io/client-go/rest"
//	"k8s.io/client-go/tools/clientcmd"
//)
//
//func Demo1() {
//	// config
//	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
//	if err != nil {
//		panic(err)
//	}
//	config.GroupVersion = &v1.SchemeGroupVersion
//	config.NegotiatedSerializer = scheme.Codecs
//	config.APIPath = "/api"
//
//	// client
//	restClient, err := rest.RESTClientFor(config)
//	if err != nil {
//		panic(err)
//	}
//
//	// get data
//	pod := v1.PodList{}
//	err = restClient.Get().Namespace("default").Resource("pods").Do(context.TODO()).Into(&pod)
//	if err != nil {
//		panic(err)
//	} else {
//		fmt.Println(len(pod.Items))
//	}
//}
//
