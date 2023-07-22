package Index

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	// 索引名分类
	NamespaceIndexName = "namespace"
	NodeNameIndexName  = "nodeName"
)

// 索引器函数
func NamespaceIndexFunc(obj interface{}) ([]string, error) {
	pod, ok := obj.(v1.Pod)
	if !ok {
		return nil, fmt.Errorf("类型错误")
	}
	return []string{pod.Namespace}, nil
}

// 基于nodename 的索引器函数
func NodeNameIndexFuc(obj interface{}) ([]string, error) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		return nil, fmt.Errorf("类型错误")
	}
	return []string{pod.Spec.NodeName}, nil
}

func Demo1() {
	// 索引的方式，需要我们自己实现
	index := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{
		NamespaceIndexName: NamespaceIndexFunc,
		NodeNameIndexName:  NodeNameIndexFuc,
	})

	pod1 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-1",
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			NodeName: "node1",
		},
	}

	pod2 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-2",
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			NodeName: "node2",
		},
	}
	pod3 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-3",
			Namespace: "kube-system",
		},
		Spec: v1.PodSpec{
			NodeName: "node2",
		},
	}
	// pod加入到索引器中
	index.Add(pod1)
	index.Add(pod2)
	index.Add(pod3)

	// 查询, indexname=索引器名称，indexvalue=索引键
	pods, err := index.ByIndex(NamespaceIndexName, "default")
	if err != nil {
		panic(err)
	}
	for _, pod := range pods {
		fmt.Println(pod.(*v1.Pod).Name)
	}
}
