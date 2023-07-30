package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	// 后台删除，不夯前台，提升删除速度
	deletePolicy metav1.DeletionPropagation = "Background"
	// 立即删除
	gracetime     int64                = 0
	deleteOptions metav1.DeleteOptions = metav1.DeleteOptions{
		PropagationPolicy:  &deletePolicy,
		GracePeriodSeconds: &gracetime,
	}
)

// 初始化client
func initClientSet() (*kubernetes.Clientset, error) {
	// 获取reset config
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		return nil, err
	}

	// 实例clienetset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func GetClient() (*kubernetes.Clientset, error) {
	return initClientSet()
}

//func main() {
//	// 初始化flag
//	var single = pflag.BoolP("single", "s", false, "删除异常POD(default false)")
//	var multi = pflag.BoolP("multi", "m", false, "批量删除POD(default false)")
//	var restart = pflag.BoolP("restart", "r", false, "删除重启次数大于xx的POD(default false)")
//	var restartCount = pflag.Int32P("count", "c", 10, "过滤重启次数大于xx的POD(default 10)")
//	var threads = pflag.Int32P("thread", "t", 1, "删除并发数")
//
//	pflag.Parse()
//
//	// 初始化clientset
//	clientset := initClientSet()
//	if *multi {
//		multiDeleteErrorPods(clientset)
//	} else if *single {
//		singleDeletePods(clientset, *threads)
//	} else if *restart {
//		singleDeleteRestartPods(clientset, *restartCount, *threads)
//	} else {
//		pflag.PrintDefaults()
//	}
//}
