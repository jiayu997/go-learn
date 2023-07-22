package Project2

import (
	"context"
	v13 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	informer "k8s.io/client-go/informers/core/v1"
	ingressInformer "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	service "k8s.io/client-go/listers/core/v1"
	ingress "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"reflect"
	"time"
)

type controller struct {
	clientset     kubernetes.Interface
	ingressLister ingress.IngressLister
	serviceLister service.ServiceLister
	queue         workqueue.RateLimitingInterface
}

func (c *controller) addService(obj interface{}) {
	c.enqueue(obj)
}

func (c *controller) updateService(oldObj interface{}, newObj interface{}) {
	// 比较anotation
	if reflect.DeepEqual(oldObj, newObj) {
		return
	}

	c.enqueue(newObj)
}

func (c *controller) deleteIngress(obj interface{}) {
	ingress := obj.(*v1.Ingress)
	ownerReference := v12.GetControllerOf(ingress)

	if ownerReference == nil {
		return
	}

	if ownerReference.Kind != "Service" {
		return
	}

	c.queue.Add(ingress.Namespace + "/" + ingress.Name)

}

func (c *controller) Run(stopCh chan struct{}) {
	for i := 0; i < 5; i++ {
		go wait.Until(c.worker, time.Minute, stopCh)
	}

	<-stopCh
}

func (c *controller) worker() {
	for c.processNextItem() {

	}
}

// 从queue 获取key
func (c *controller) processNextItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	// 当key处理完后，需要在queue中移除
	defer c.queue.Done(item)

	key := item.(string)

	err := c.sysncService(key)
	if err != nil {
		c.handlerError(key, err)
	}
	return true
}

func (c *controller) sysncService(key string) error {
	namespaceKey, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// 删除
	service, err := c.serviceLister.Services(namespaceKey).Get(name)
	if errors.IsNotFound(err) {
		return nil
	}

	if err != nil {
		return err
	}

	// 新增和删除
	_, ok := service.GetAnnotations()["ingress/http"]
	ingress, err := c.ingressLister.Ingresses(namespaceKey).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if ok && errors.IsNotFound(err) {
		// 如果ingress不存在
		ig := c.constructIngress(service)
		_, err := c.clientset.NetworkingV1().Ingresses(namespaceKey).Create(context.TODO(), ig, v12.CreateOptions{})
		if err != nil {
			return err
		}

	} else if !ok && ingress != nil {
		// 删除ingress
		err := c.clientset.NetworkingV1().Ingresses(namespaceKey).Delete(context.TODO(), name, v12.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil

}

func (c *controller) constructIngress(service *v13.Service) *v1.Ingress {
	ingress := v1.Ingress{}

	ingress.ObjectMeta.OwnerReferences = []v12.OwnerReference{
		*v12.NewControllerRef(service, v13.SchemeGroupVersion.WithKind("Service")),
	}

	ingress.Namespace = service.Namespace
	ingress.Name = service.Name
	pathType := v1.PathTypePrefix
	icn := "nginx"
	ingress.Spec = v1.IngressSpec{
		IngressClassName: &icn,
		Rules: []v1.IngressRule{
			{
				Host: "example.com",
				IngressRuleValue: v1.IngressRuleValue{
					HTTP: &v1.HTTPIngressRuleValue{
						Paths: []v1.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: v1.IngressBackend{
									Service: &v1.IngressServiceBackend{
										Name: service.Name,
										Port: v1.ServiceBackendPort{
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
	}
	return &ingress
}

func (c *controller) handlerError(key string, err error) {
	if c.queue.NumRequeues(key) <= 10 {
		c.queue.AddRateLimited(key)
		return
	}
	runtime.HandleError(err)
	c.queue.Forget(key)
}

func (c *controller) enqueue(obj interface{}) {
	// 生成objkey
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}

	// 队列放对象的key
	c.queue.Add(key)
}

func NewController(client kubernetes.Interface, serviceInformer informer.ServiceInformer, ingressInformer ingressInformer.IngressInformer) controller {
	c := controller{
		clientset: client,
		// Lister函数执行后(serviceInformer.Informer也会注册informer)，会生成一个Informer，并加到sharedInformerFactory struct中的informers中，这也就是为什么factory.start要后执行(start 会去判断这个factory中是否有informer在)
		// lister() 返回一个type serviceLister struct,通过这个，我们可以去index索引器中获取数据
		// Lister() 可以让我们根据给定的标签选择器，去给定的index索引器中，查找相应的service资源集合
		serviceLister: serviceInformer.Lister(),
		ingressLister: ingressInformer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
	}
	// 在service和ingressInformer lister执行完成后，创建serviceInformer的Factory 中的informers有二个了, === 将这个informer注册到Factory中

	// AddEventHandler中有一个函数 unc (p *processorListener) run() 启动协程一直在监听Add/Update/Delete事件
	// serviceInformer.Informer会返回一个service资源类型的shareindexformer结构体(他实现了ShareIndexInformer),如果这个shareindexinformer已被注册，则返回这个
	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
	})

	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})
	return c
}
