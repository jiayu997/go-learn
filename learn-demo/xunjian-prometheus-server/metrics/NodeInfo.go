package metrics

import (
	"fmt"
	"log"
	"sync"
	"xunjian-prometheus-server/kubernetes"
	"xunjian-prometheus-server/tool"
)

type Node struct {
	NodeName         string
	IP               string
	Schedulable      bool
	Ready            string // 节点是否ready Ready|NotReady
	CpuUsePercent    float64
	MemoryUsePercent float64
	DiskUsePercent   float64
	CpuTotal         int     // 多少颗CPU
	Load15           float64 // 15分钟
	OverFlow         bool
}

func getnodelist() []Node {

	result := kubernetes.GetNodes()

	// 判断prometheus配置文件是否配置
	if tool.Conf.Prometheus["ip"] == "" || tool.Conf.Prometheus["port"] == "" {
		log.Fatal("prometheus config error")
	}

	nodelist := make([]Node, len(result.Items))

	for index, node := range result.Items {
		nodelist[index].NodeName = node.Name
		nodelist[index].IP = node.Status.Addresses[0].Address

		// 存在bug，如果节点宕机，但是并没有打上不可调度标签的话，会导致也是显示可调度的
		// 这里用来表示这个节点是否可调度
		nodelist[index].Schedulable = !node.Spec.Unschedulable

		// NetworkUnavailable
		// fmt.Println(node.Status.Conditions[0].Type, node.Status.Conditions[0].Status)
		// MemoryPressure
		// fmt.Println(node.Status.Conditions[1].Type, node.Status.Conditions[1].Status)
		// DiskPressure
		// fmt.Println(node.Status.Conditions[2].Type, node.Status.Conditions[2].Status)
		// PIDPressure
		// fmt.Println(node.Status.Conditions[3].Type, node.Status.Conditions[3].Status)
		// 节点是否Ready
		if node.Status.Conditions[4].Status == "True" {
			nodelist[index].Ready = "Ready"
		} else {
			nodelist[index].Ready = "NotReady"
		}
		nodelist[index].CpuUsePercent = getNodeCpuPercent(node.Name)
		nodelist[index].MemoryUsePercent = getNodeMemoryPercent(node.Name)
		nodelist[index].CpuTotal = getNodeCpuTotal(node.Name)
		nodelist[index].Load15 = getNodeLoad15(node.Name)
		nodelist[index].DiskUsePercent = getNodeDiskPercent(node.Name)
		if nodelist[index].Load15 > float64(nodelist[index].CpuTotal) {
			nodelist[index].OverFlow = true
		}
	}
	return nodelist
}

func GetNodeInfo(wg *sync.WaitGroup, resultChan chan interface{}) {
	if tool.Conf.Controller["node"] != "true" {
		wg.Done()
		return
	}
	fmt.Println("--------------- 开始节点检查 ----------")
	NodeList := getnodelist()
	fmt.Println("--------------- 结束节点检查 ----------")
	resultChan <- NodeList
	wg.Done()
}
