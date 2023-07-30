package k8s

import (
	"context"
	"errors"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 删除ip这个节点
func DeleteNode(ip string) error {
	var exist bool
	var info map[string]string = make(map[string]string)

	client, err := initClientSet()
	if err != nil {
		return err
	}
	// 获取所有的NODE
	nodelist, err := client.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, node := range nodelist.Items {
		for _, nip := range node.Status.Addresses {
			if nip.Address == ip {
				exist = true
				info["nodename"] = node.Name
				info["address"] = nip.Address
			}
		}
	}
	if exist {
		err := client.CoreV1().Nodes().Delete(context.TODO(), info["nodename"], v1.DeleteOptions{})
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return errors.New("node not found")
	}
}

// 判断这个ip是否在集群
func SearchNode(ip string) (bool, error) {
	client, err := initClientSet()
	if err != nil {
		return false, err
	}

	// 获取当前K8S集群所有node节点信息
	nodelist, err := client.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, node := range nodelist.Items {
		for _, nip := range node.Status.Addresses {
			if nip.Address == ip {
				return true, nil
			}
		}
	}
	return false, errors.New("node not found")
}
