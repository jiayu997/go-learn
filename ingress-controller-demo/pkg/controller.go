package pkg

import (
	"context"
	ingressv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	informercorev1 "k8s.io/client-go/informers/core/v1"
	informernetworkv1 "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/listers/core/v1"
	networkv1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"reflect"
	"time"
)

const workNum = 5
const retry = 5

type controller struct {
	client        kubernetes.Interface
	ingressLister networkv1.IngressLister
	serviceLister corev1.ServiceLister
	queue         workqueue.RateLimitingInterface
}

func (c controller) addService(obj interface{}) {
	c.enqueue(obj)
}

func (c controller) updateService(oldObj interface{}, newObj interface{}) {
	// 应该要比较annotation
	if reflect.DeepEqual(oldObj, newObj) {
		return
	}
	c.enqueue(newObj)
}

func (c controller) deleteIngress(obj interface{}) {
	ingress := obj.(*ingressv1.Ingress)
	ownerReference := metav1.GetControllerOf(ingress)

	if ownerReference == nil {
		return
	}

	if ownerReference.Kind != "Service" {
		return
	}

	c.queue.Add(ingress.Namespace + "/" + ingress.Name)
}

// 队列里面只放key
func (c *controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceIndexFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}

	// 入队列
	c.queue.Add(key)
}

func NewController(client kubernetes.Interface, serviceInformer informercorev1.ServiceInformer, ingressInformer informernetworkv1.IngressInformer) *controller {
	c := controller{
		client:        client,
		ingressLister: ingressInformer.Lister(),
		serviceLister: serviceInformer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
	}

	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
	})

	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})

	return &c
}

func (c *controller) Run(stopCh chan struct{}) {
	for i := 0; i < workNum; i++ {
		// 实际处理逻辑
		go wait.Until(c.worker, time.Minute, stopCh)
	}
	<-stopCh
}

func (c *controller) worker() {
	for c.processNextItem() {

	}

}

func (c *controller) processNextItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	// 标记处理完成
	defer c.queue.Done(item)

	key := item.(string)

	// 调协service
	err := c.syncService(key)
	if err != nil {
		c.handlerError(item, err)
	}
	return true
}

func (c *controller) syncService(key string) error {
	namespaceKey, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// 删除
	service, err := c.serviceLister.Services(namespaceKey).Get(name)
	// 因为删除后，reflector会更新本地DeltaFIFO，然后更新indexer，所以service就不会存在
	if errors.IsNotFound(err) {
		return nil
	}

	if err != nil {
		return err
	}

	// 新增和更新
	_, ok := service.GetAnnotations()["ingress/http"]
	ingress, err := c.ingressLister.Ingresses(namespaceKey).Get(name)

	// ingress不是不存在
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	// service存在而ingress不存在
	if ok && errors.IsNotFound(err) {
		// 创建ingress
		ig := c.constructIngress(namespaceKey, name)
		_, err := c.client.NetworkingV1().Ingresses(namespaceKey).Create(context.TODO(), ig, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else if !ok && ingress != nil {
		// 删除ingress
		err := c.client.NetworkingV1().Ingresses(namespaceKey).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *controller) handlerError(item interface{}, err error) {
	if c.queue.NumRequeues(item) < retry {
		c.queue.AddRateLimited(item)
	}
	// 不让重试了
	runtime.HandleError(err)
	c.queue.Forget(item)
}

func (c *controller) constructIngress(namespaceKey string, name string) *ingressv1.Ingress {
	pathType := ingressv1.PathTypePrefix
	return &ingressv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceKey,
		},
		Spec: ingressv1.IngressSpec{
			Rules: []ingressv1.IngressRule{
				{
					Host: "example.com",
					IngressRuleValue: ingressv1.IngressRuleValue{
						HTTP: &ingressv1.HTTPIngressRuleValue{
							Paths: []ingressv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: ingressv1.IngressBackend{
										Service: &ingressv1.IngressServiceBackend{
											Name: name,
											Port: ingressv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
