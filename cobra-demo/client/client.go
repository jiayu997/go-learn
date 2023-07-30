package client

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
)

// init client set from kubeconfig
func InitKubeConfigFromKubeConfig() (clientset *kubernetes.Clientset) {
	var kubeConfig *string

	// init kubeconfig path
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// load kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		log.Fatal("no kubeconfig")
	}

	// create clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return clientset
}

// init client set from
