package Informer

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func Demo2() {
	// config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	// client set
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 初始化 informer factory, 30s全量找本地缓存同步一次
	factory := informers.NewSharedInformerFactory(clientset, time.Second*30)

	// 监听想要获取的对象informer
	deploymentInformer := factory.Apps().V1().Deployments()

	// 注册一下informer
	informer := deploymentInformer.Informer()

	// 创建Lister,只要第一次同步过来后,lister就可以从缓存中获取数据
	deploymentLister := deploymentInformer.Lister()

	// 注册事件处理程序
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deploy := obj.(*v1.Deployment)
			fmt.Println("add a deployment: ", deploy.Name, deploy.Namespace)

		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldDeploy := oldObj.(*v1.Deployment)
			newDeploy := oldObj.(*v1.Deployment)
			fmt.Println("update a deployment: ", oldDeploy.Name, newDeploy.Name)

		},
		DeleteFunc: func(obj interface{}) {
			deploy := obj.(*v1.Deployment)
			fmt.Println("delete a deployment", deploy.Name)
		},
	})

	// 启动informer (启动list & watch)
	stopCh := make(chan struct{})
	defer close(stopCh)
	factory.Start(stopCh)

	// 等待所有的informer同步完成
	factory.WaitForCacheSync(stopCh)

	// 通过Lister 获取缓存中的deployment,如果缓存中没有数据，则不会有数据
	deployments, err := deploymentLister.Deployments("monitor").List(labels.Everything())
	if err != nil {
		panic(err)
	}
	for index, deploy := range deployments {
		fmt.Printf("%d -> %s\n", index, deploy.Name)
	}

	// 阻塞主进程
	<-stopCh
}
