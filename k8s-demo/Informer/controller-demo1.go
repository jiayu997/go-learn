package Informer

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

func initClient() (*kubernetes.Clientset, error) {
	// config
	config, _ := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)

	// clientset
	return kubernetes.NewForConfig(config)
}

// pod控制器
type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
}

func NewController(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller) *Controller {
	return &Controller{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
	}
}

func (c *Controller) Run(stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// 停止控制器后，需要关闭队列
	defer c.queue.ShuttingDown()

	// 启动控制器
	klog.Info("start pod controller")

	// 启动通用控制器
	go c.informer.Run(stopCh)

	// 等待缓存同步完成
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Time out for caches to sync"))
	}

	//
}

func ControllerDemo1() {

}
