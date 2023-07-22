package Project1

//
//import (
//	"context"
//	v1 "k8s.io/api/networking/v1"
//	"k8s.io/apimachinery/pkg/api/errors"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/apimachinery/pkg/util/wait"
//	informer "k8s.io/client-go/informers/core/v1"
//	ingressInformer "k8s.io/client-go/informers/networking/v1"
//	"k8s.io/client-go/kubernetes"
//	coreList "k8s.io/client-go/listers/core/v1"
//	ingressList "k8s.io/client-go/listers/networking/v1"
//	"k8s.io/client-go/tools/cache"
//	"k8s.io/client-go/util/workqueue"
//	"reflect"
//	"time"
//)
//
//const workNum = 5
//
//type controller struct {
//	client        kubernetes.Interface
//	ingressLister ingressList.IngressLister
//	serviceList   coreList.ServiceLister
//	queue         workqueue.RateLimitingInterface
//}
//
//func (c *controller) enqueue(obj interface{}) {
//	key, err := cache.MetaNamespaceKeyFunc(obj)
//	if err != nil {
//		return
//	}
//	c.queue.Add(key)
//}
//
//func (c *controller) deleteIngress(obj interface{}) {
//	ingress := obj.(*v1.Ingress)
//	ownerReference := v12.GetControllerOf(ingress)
//	if ownerReference != nil {
//		return
//	}
//	if ownerReference.Kind != "Service" {
//		return
//	}
//	c.queue.Add(ingress.Namespace + "/" + ingress.Name)
//}
//
//func (c *controller) addService(obj interface{}) {
//	c.enqueue(obj)
//}
//
//func (c *controller) updateService(oldObj interface{}, newObj interface{}) {
//	if reflect.DeepEqual(oldObj, newObj) {
//		return
//	}
//	c.enqueue(newObj)
//}
//
//func (c *controller) worker() {
//	for c.processNextItem() {
//
//	}
//}
//
//func (c *controller) constructIngress(namespaceKey string, name string) v1.Ingress {
//	ingress := v1.Ingress{}
//	ingress.Name = name
//	ingress.Namespace = namespaceKey
//	pathType := v1.PathTypePrefix
//	ingress.Spec = v1.IngressSpec{
//		Rules: []v1.IngressRule{
//			{
//				Host: "example.com",
//				IngressRuleValue: v1.IngressRuleValue{
//					HTTP: &v1.HTTPIngressRuleValue{
//						Paths: []v1.HTTPIngressPath{
//							{
//								Path:     "/",
//								PathType: &pathType,
//								Backend: v1.IngressBackend{
//									Service: &v1.IngressServiceBackend{
//										Name: name,
//										Port: v1.ServiceBackendPort{
//											Number: 80,
//										},
//									},
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//	}
//	return ingress
//}
//
//func (c *controller) syncService(key string) error {
//	namespaceKey, name, err := cache.SplitMetaNamespaceKey(key)
//	if err != nil {
//		return err
//	}
//	// 删除
//	service, err := c.serviceList.Services(namespaceKey).Get(name)
//	if errors.IsNotFound(err) {
//		return nil
//	}
//	if err != nil {
//		return err
//	}
//	// 新增和删除
//	_, ok := service.GetAnnotations()["ingress/http"]
//	ingress, err := c.ingressLister.Ingresses(namespaceKey).Get(name)
//	if err != nil {
//		return err
//	}
//
//	if ok && errors.IsNotFound(err) {
//		// create ingress
//		ig := c.constructIngress(namespaceKey, name)
//		c.client.NetworkingV1().Ingresses(namespaceKey).Create(context.TODO(), ig, metav1.CreateOptions{})
//	} else if !ok && ingress != nil {
//		// delete ingress
//		c.client.NetworkingV1().Ingresses(namespaceKey).Delete(context.TODO(), name, metav1.DeleteOptions{})
//	}
//}
//
//func (c *controller) handlerError(err interface{}) {
//
//}
//
//func (c *controller) processNextItem() bool {
//	item, shutdown := c.queue.Get()
//	if shutdown {
//		return false
//	}
//	defer c.queue.Done(item)
//	key := item.(string)
//
//	err := c.syncService(key)
//	if err != nil {
//		c.handlerError(err)
//		return false
//	}
//}
//
//func (c *controller) Run(stopCh chan struct{}) {
//	for i := 0; i < workNum; i++ {
//		go wait.Until(c.worker, time.Minute, stopCh)
//	}
//	<-stopCh
//}
//
//func NewController(client kubernetes.Interface, serviceInformer informer.ServiceInformer, ingressInformer ingressInformer.IngressInformer) controller {
//	c := controller{
//		client:        client,
//		ingressLister: ingressInformer.Lister(),
//		serviceList:   serviceInformer.Lister(),
//		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
//	}
//	// service 的新增与更新
//	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
//		AddFunc:    c.addService,
//		UpdateFunc: c.updateService,
//	})
//
//	// ingress 删除逻辑
//	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
//		DeleteFunc: c.deleteIngress,
//	})
//	return c
//}
//
