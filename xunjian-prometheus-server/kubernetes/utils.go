package kubernetes

import (
	"context"
	"log"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func generateConfig() *restclient.Config {
	// 加载配置文件
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func GetPods() *corev1.PodList {
	config := generateConfig()

	// 配置API路径
	config.APIPath = "api"

	// 配置GV版本
	config.GroupVersion = &corev1.SchemeGroupVersion

	// 配置数据序列化
	config.NegotiatedSerializer = scheme.Codecs

	// 定义数据接收变量
	result := &corev1.PodList{}

	// 实例化rest client
	restClient, err := rest.RESTClientFor(config)
	if err != err {
		log.Fatal(err)
	}

	err = restClient.Get().
		Resource("pods").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func GetNodes() *corev1.NodeList {
	config := generateConfig()

	// 配置API路径
	config.APIPath = "api"

	// 配置GV版本
	config.GroupVersion = &corev1.SchemeGroupVersion

	// 配置数据序列化
	config.NegotiatedSerializer = scheme.Codecs

	// 定义数据接收变量
	result := &corev1.NodeList{}

	// 实例化
	restClient, err := rest.RESTClientFor(config)
	if err != err {
		log.Fatal(err)
	}

	err = restClient.Get().
		Resource("nodes").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func GetDNSServer() *corev1.Service {
	config := generateConfig()

	// 配置API路径
	config.APIPath = "api"

	config.GroupVersion = &corev1.SchemeGroupVersion

	config.NegotiatedSerializer = scheme.Codecs

	result := &corev1.Service{}

	restclient, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatal(err)
	}
	err = restclient.Get().
		Resource("services").
		Namespace("kube-system").
		Name("kube-dns").
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func GetPVCs() *corev1.PersistentVolumeClaimList {
	config := generateConfig()
	config.APIPath = "api"

	config.GroupVersion = &corev1.SchemeGroupVersion

	config.NegotiatedSerializer = scheme.Codecs

	result := &corev1.PersistentVolumeClaimList{}

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = restClient.Get().
		Resource("PersistentVolumeClaims").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(result)
	if err != nil {
		log.Fatal(err.Error())
	}
	return result
}

func GetPVs() *corev1.PersistentVolumeList {
	config := generateConfig()

	config.APIPath = "api"

	config.GroupVersion = &corev1.SchemeGroupVersion

	config.NegotiatedSerializer = scheme.Codecs

	result := &corev1.PersistentVolumeList{}

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatal(err)
	}

	err = restClient.Get().
		Resource("PersistentVolumes").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(result)
	if err != nil {
		log.Fatal(err.Error())
	}
	return result
}

func GetDeployment(name, namespace string) *v1.Deployment {
	config := generateConfig()

	// apis
	config.APIPath = "apis"

	// apps/v1
	config.GroupVersion = &v1.SchemeGroupVersion

	config.NegotiatedSerializer = scheme.Codecs

	//result := &v1.DeploymentList{}
	result := &v1.Deployment{}

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatal(err)
	}
	err = restClient.Get().
		Resource("deployments").
		Namespace(namespace).
		Name(name).
		Do(context.TODO()).
		Into(result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func GetStatefulset(name, namespace string) *v1.StatefulSet {
	config := generateConfig()

	config.APIPath = "apis"

	config.GroupVersion = &v1.SchemeGroupVersion

	config.NegotiatedSerializer = scheme.Codecs

	result := &v1.StatefulSet{}

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatal(err)
	}
	err = restClient.Get().
		Resource("statefulsets").
		Namespace(namespace).
		Name(name).
		Do(context.TODO()).
		Into(result)
	if err != nil {
		log.Fatal(err)
	}
	return result

}

func GetDaemonset(name, namespace string) *v1.DaemonSet {
	config := generateConfig()

	config.APIPath = "apis"

	config.GroupVersion = &v1.SchemeGroupVersion

	config.NegotiatedSerializer = scheme.Codecs

	result := &v1.DaemonSet{}

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatal(err)
	}
	err = restClient.Get().
		Resource("daemonsets").
		Namespace(namespace).
		Name(name).
		Do(context.TODO()).
		Into(result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
