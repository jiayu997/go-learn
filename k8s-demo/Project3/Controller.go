package Project3

import (
	"fmt"
	flag "github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"path/filepath"
	"time"
)

type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
}

// Run开始和watch同步
func (c Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// 停止控制器后关闭队列
	defer c.queue.ShutDown()

	// 启动
	go c.informer.Run(stopCh)

	// 等待所有相关的缓存同步，然后再开始处理队列中的项目
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Time out waiting for caches to sync"))
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWoker, time.Second, stopCh)
	}
	<-stopCh
}

func (c Controller) runWoker() {
	for c.processNextItem() {

	}
}

func (c Controller) processNextItem() bool {
	// 等到工作队列有一个新元素
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	// 告诉队列我们已经完成了处理此key的操作，将为其他worker解锁该key
	defer c.queue.Done(key)

	// 调用保护业务逻辑的方法
	err := c.syncToStdout(key.(string))

	// 如果在执行业务逻辑过程中出现错误，则处理错误
	c.handleErr(err, key)
	return true
}

func (c Controller) syncToStdout(key string) error {
	// 从本地存储中获取key对应的对象
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		klog.Error("Fetching object with key %s from store failed with %v", key, err)
		return err
	}
	if !exists {
		fmt.Printf("Pod %s does not exits anymore \n", key)
	} else {
		fmt.Printf("Sync/Add/Update for Pod %s\n", obj.(*v1.Pod).GetName())
	}
	return nil
}

func (c Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// 忘记每次成功同步时 key 的#AddRateLimited历史记录。
		// 这样可以确保不会因过时的错误历史记录而延迟此 key 更新的以后处理
		c.queue.Forget(key)
		return
	}
	// 如果出现问题，重试5次
	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing pod %v:%v", key, err)
		c.queue.AddRateLimited(key)
		return
	}
	c.queue.Forget(key)
	// 多次重试，无法处理key
	runtime.HandleError(err)
}

func NewController(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller) *Controller {
	return &Controller{
		informer: informer,
		queue:    queue,
		indexer:  indexer,
	}
}

func initClient() (*kubernetes.Clientset, error) {
	var err error
	var config *rest.Config
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(可选)kubeconfig文件绝对路径")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig 文件绝对路径")
	}

	if config, err = rest.InClusterConfig(); err != nil {
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			panic(err.Error())
		}
	}

	return kubernetes.NewForConfig(config)
}

func Project3() {
	clientset, err := initClient()
	if err != nil {
		klog.Fatal(err)
	}
	// create pod listwatcher
	podListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", metav1.NamespaceAll, fields.Everything())

	// 创建队列
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	indexer, informer := cache.NewIndexerInformer(podListWatcher, &v1.Pod{}, 0*time.Second, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			//fmt.Printf("Add operation\n")
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			//fmt.Printf("Update operation\n")
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			//fmt.Printf("Delete operation\n")
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}, cache.Indexers{})

	controller := NewController(queue, indexer, informer)

	// start controller
	stopCh := make(chan struct{})
	defer close(stopCh)
	go controller.Run(1, stopCh)
	select {}
}
