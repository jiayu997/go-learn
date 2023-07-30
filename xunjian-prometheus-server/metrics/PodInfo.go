package metrics

import (
	"fmt"
	"log"
	"sync"
	"xunjian-prometheus-server/kubernetes"
	"xunjian-prometheus-server/tool"
)

type Pod struct {
	NameSpace     string   `yaml:"namespace"`
	Name          string   `yaml:"name"`
	Status        string   `yaml:"status"`
	CpuUse        float64  `yaml:"cpuuse"`
	MemoryUse     int64    `yaml:"memoryuser"`
	Restart       int32    `yaml:"restart"`
	HealthCheck   bool     `yaml:"healthcheck"`
	CpuRequest    []string `yaml:"cpurequest"`
	MemoryRequest []string `yaml:"memoryrequest"`
	CpuLimit      []string `yaml:"cpulimit"`
	MemoryLimit   []string `yaml:"memorylimit"`
	Resource      bool
}

func podlist() []Pod {
	result := kubernetes.GetPods()

	// 判断prometheus配置文件是否配置
	if tool.Conf.Prometheus["ip"] == "" || tool.Conf.Prometheus["port"] == "" {
		log.Fatal("prometheus config error")
	}

	podlist := make([]Pod, len(result.Items))

	for index, pod := range result.Items {
		// pod 名称
		podlist[index].Name = pod.Name

		// pod 命名空间
		podlist[index].NameSpace = pod.Namespace

		// pod 状态
		podlist[index].Status = string(pod.Status.Phase)

		// pod重启次数,多个容器重启相加
		if len(pod.Status.ContainerStatuses) == 0 {
			podlist[index].Restart = 0
		} else {
			for _, j := range pod.Status.ContainerStatuses {
				podlist[index].Restart += j.RestartCount
			}
		}

		// pod是否配置检查检查
		flag := false
		for _, k := range pod.Spec.Containers {
			//fmt.Println(k.LivenessProbe)
			//fmt.Println(k.ReadinessProbe)
			// 如果没有，默认返回0
			podlist[index].CpuRequest = append(podlist[index].CpuRequest, k.Resources.Requests.Cpu().String())
			podlist[index].MemoryRequest = append(podlist[index].MemoryRequest, k.Resources.Requests.Memory().String())
			podlist[index].CpuLimit = append(podlist[index].CpuLimit, k.Resources.Limits.Cpu().String())
			podlist[index].MemoryLimit = append(podlist[index].MemoryLimit, k.Resources.Limits.Memory().String())
			if k.LivenessProbe != nil || k.ReadinessProbe != nil {
				flag = true
			}
		}
		if flag {
			podlist[index].HealthCheck = true
		} else {
			podlist[index].HealthCheck = false
		}
		podlist[index].MemoryUse = getPodMemory(podlist[index].NameSpace, podlist[index].Name)
		podlist[index].CpuUse = getPodCpu(podlist[index].NameSpace, podlist[index].Name)
	}
	return podlist
}

func GetPodInfo(wg *sync.WaitGroup, resultChan chan interface{}) {
	if tool.Conf.Controller["pod"] != "true" {
		wg.Done()
		return
	}
	fmt.Println("--------------- 开始Pod检查 ----------")
	PodList := podlist()
	fmt.Println("--------------- 结束Pod检查 ----------")
	resultChan <- PodList
	wg.Done()
}
