package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetAllSVC() {
	client, err := initClientSet()
	if err != nil {
		fmt.Println(err.Error())
	}
	serviceList, err := client.CoreV1().Services("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Printf(err.Error())
	}
	for _, service := range serviceList.Items {
		fmt.Println(service.Name, service.Namespace, service.Spec.Type)
	}

}
