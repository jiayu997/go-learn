package metrics

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"
	"xunjian-prometheus-server/kubernetes"
	"xunjian-prometheus-server/tool"

	v1 "k8s.io/api/core/v1"
)

type PodInfo struct {
	PodName      string
	PodIP        string
	Status       string
	CpuUse       float64
	MemoryUse    int64
	RestartCount int64
	HealthCheck  bool
}

type BusinnessMetric struct {
	Name            string
	Namespace       string
	TypeName        string // 业务名称
	ControllerName  string // deployment|daemonset|statefulset
	Replicas        int
	HostNetwork     bool
	HealthCheck     bool
	CpuRequest      []string
	MemoryRequest   []string
	CpuLimit        []string
	MemoryLimit     []string
	VolumeMounts    []string
	UserHealth_Path string
	UserHealth_port string
	Pod             []PodInfo
}

func initBusinnessConfig(BusinnessList *[]BusinnessMetric) {
	// 只创建属于deployment的
	for _, businness := range tool.Conf.CheckList.Business {
		if businness["name"] != "" && businness["controller"] != "" && businness["namespace"] != "" {
			var tmp BusinnessMetric
			tmp.Name = businness["name"]
			tmp.TypeName = businness["type_name"]
			tmp.Namespace = businness["namespace"]
			tmp.ControllerName = businness["controller"]
			tmp.UserHealth_Path = businness["health_path"]
			tmp.UserHealth_port = businness["health_port"]
			*BusinnessList = append(*BusinnessList, tmp)
		}
	}
}

func handleHealthPod(path, port, ip string) bool {
	if path == "" || port == "" {
		return false
	}
	// 设置超时
	client := http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest("GET", path+":"+port, nil)
	if err != nil {
		return false
	}
	response, err := client.Do(req)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return false
	} else {
		return true
	}
	//	body, _ := ioutil.ReadAll(response.Body)
	//	if tool.CompareString(string(body), "Ok") {
	//		return false
	//	} else {
	//		return true
	//	}
}

//func getDeploymentInfos() {
//	//all_deployment := kubernetes.GetDeployments()
//	//deployment := kubernetes.GetDeployment("coredns", "kube-system")
//	//fmt.Println(deployment)
//
//	// deployment处理
//	podlist := kubernetes.GetPods()
//	for index := 0; index < len(DeploymentList); index++ {
//		tmp_deployment := kubernetes.GetDeployment(DeploymentList[index].Name, DeploymentList[index].Namespace)
//		DeploymentList[index].Replicas = int(*tmp_deployment.Spec.Replicas)
//		DeploymentList[index].HostNetwork = tmp_deployment.Spec.Template.Spec.HostNetwork
//		if len(tmp_deployment.Spec.Template.Spec.Containers) != 0 {
//			for _, j := range tmp_deployment.Spec.Template.Spec.Containers {
//				// 健康检查判断
//				if j.LivenessProbe != nil || j.ReadinessProbe != nil {
//					DeploymentList[index].HealthCheck = true
//				}
//				// request获取
//				DeploymentList[index].CpuRequest = append(DeploymentList[index].CpuRequest, j.Resources.Requests.Cpu().String())
//				DeploymentList[index].MemoryRequest = append(DeploymentList[index].MemoryRequest, j.Resources.Requests.Memory().String())
//				DeploymentList[index].CpuLimit = append(DeploymentList[index].CpuLimit, j.Resources.Limits.Cpu().String())
//				DeploymentList[index].MemoryLimit = append(DeploymentList[index].MemoryLimit, j.Resources.Limits.Memory().String())
//				for _, k := range j.VolumeMounts {
//					DeploymentList[index].VolumeMounts = append(DeploymentList[index].VolumeMounts, k.MountPath)
//				}
//			}
//		}
//		//处理deployment下的pods
//		for _, pod := range podlist.Items {
//			//fmt.Println(DeploymentList[index].Name, pod.Name)
//			//fmt.Println(tool.CompareString(pod.Name, "^"+DeploymentList[index].Name+"(.*)"))
//			if tool.CompareString(pod.Name, "^"+DeploymentList[index].Name+`-[0-9a-zA-Z]{2,}-[0-9a-zA-Z]{2,}$`) {
//				DeploymentList[index].Pod = append(DeploymentList[index].Pod, DeployPod{
//					PodName:      DeploymentList[index].Name,
//					PodIP:        pod.Status.PodIP,
//					CpuUse:       getPodCpu(pod.Namespace, pod.Name),
//					MemoryUse:    getPodMemory(pod.Namespace, pod.Name),
//					RestartCount: getPodRestartCount(pod.Namespace, pod.Name),
//					Status:       string(pod.Status.Phase),
//					HealthCheck:  handleHealthPod(DeploymentList[index].UserHealth_Path, DeploymentList[index].UserHealth_Path, pod.Status.PodIP),
//				})
//			}
//		}
//	}
//	//fmt.Println(DeploymentList)
//}

func getDeploymentInfo(podlist *v1.PodList, business BusinnessMetric) BusinnessMetric {
	tmp_deployment := kubernetes.GetDeployment(business.Name, business.Namespace)
	business.Replicas = int(*tmp_deployment.Spec.Replicas)
	business.HostNetwork = tmp_deployment.Spec.Template.Spec.HostNetwork
	if len(tmp_deployment.Spec.Template.Spec.Containers) != 0 {
		for _, j := range tmp_deployment.Spec.Template.Spec.Containers {
			if j.LivenessProbe != nil || j.ReadinessProbe != nil {
				business.HealthCheck = true
			}
			business.CpuRequest = append(business.CpuRequest, j.Resources.Requests.Cpu().String())
			business.MemoryRequest = append(business.MemoryRequest, j.Resources.Requests.Memory().String())
			business.CpuLimit = append(business.CpuLimit, j.Resources.Limits.Cpu().String())
			business.MemoryLimit = append(business.MemoryLimit, j.Resources.Limits.Memory().String())
			for _, k := range j.VolumeMounts {
				business.VolumeMounts = append(business.VolumeMounts, k.MountPath)
			}
		}
	}
	for _, pod := range podlist.Items {
		if tool.CompareString(pod.Name, "^"+business.Name+`-[0-9a-zA-Z]{2,}-[0-9a-zA-Z]{2,}$`) {
			business.Pod = append(business.Pod, PodInfo{
				PodName:      business.Name,
				PodIP:        pod.Status.PodIP,
				CpuUse:       getPodCpu(pod.Namespace, pod.Name),
				MemoryUse:    getPodMemory(pod.Namespace, pod.Name),
				RestartCount: getPodRestartCount(pod.Namespace, pod.Name),
				Status:       string(pod.Status.Phase),
				HealthCheck:  handleHealthPod(business.UserHealth_Path, business.UserHealth_port, pod.Status.PodIP),
			})
		}
	}
	return business
}

func getDaemonsetInfo(podlist *v1.PodList, business BusinnessMetric) BusinnessMetric {
	tmp_daemonset := kubernetes.GetDaemonset(business.Name, business.Namespace)
	business.HostNetwork = tmp_daemonset.Spec.Template.Spec.HostNetwork
	if len(tmp_daemonset.Spec.Template.Spec.Containers) != 0 {
		for _, j := range tmp_daemonset.Spec.Template.Spec.Containers {
			if j.LivenessProbe != nil || j.ReadinessProbe != nil {
				business.HealthCheck = true
			}
			business.CpuRequest = append(business.CpuRequest, j.Resources.Requests.Cpu().String())
			business.MemoryRequest = append(business.MemoryRequest, j.Resources.Requests.Memory().String())
			business.CpuLimit = append(business.CpuLimit, j.Resources.Limits.Cpu().String())
			business.MemoryLimit = append(business.MemoryLimit, j.Resources.Limits.Memory().String())
			for _, k := range j.VolumeMounts {
				business.VolumeMounts = append(business.VolumeMounts, k.MountPath+"\n")
			}
		}
	}
	for _, pod := range podlist.Items {
		if tool.CompareString(pod.Name, "^"+business.Name+`-[0-9a-zA-Z]{1,}`) {
			business.Pod = append(business.Pod, PodInfo{
				PodName:      business.Name,
				PodIP:        pod.Status.PodIP,
				CpuUse:       getPodCpu(pod.Namespace, pod.Name),
				MemoryUse:    getPodMemory(pod.Namespace, pod.Name),
				RestartCount: getPodRestartCount(pod.Namespace, pod.Name),
				Status:       string(pod.Status.Phase),
				HealthCheck:  handleHealthPod(business.UserHealth_Path, business.UserHealth_port, pod.Status.PodIP),
			})
		}
	}
	business.Replicas = len(business.Pod)
	return business
}

func getStatefulsetInfo(podlist *v1.PodList, business BusinnessMetric) BusinnessMetric {
	tmp_statefulset := kubernetes.GetStatefulset(business.Name, business.Namespace)
	business.Replicas = int(*tmp_statefulset.Spec.Replicas)
	business.HostNetwork = tmp_statefulset.Spec.Template.Spec.HostNetwork
	if len(tmp_statefulset.Spec.Template.Spec.Containers) != 0 {
		for _, j := range tmp_statefulset.Spec.Template.Spec.Containers {
			if j.LivenessProbe != nil || j.ReadinessProbe != nil {
				business.HealthCheck = true
			}
			business.CpuRequest = append(business.CpuRequest, j.Resources.Requests.Cpu().String())
			business.MemoryRequest = append(business.MemoryRequest, j.Resources.Requests.Memory().String())
			business.CpuLimit = append(business.CpuLimit, j.Resources.Limits.Cpu().String())
			business.MemoryLimit = append(business.MemoryLimit, j.Resources.Limits.Memory().String())
			for _, k := range j.VolumeMounts {
				business.VolumeMounts = append(business.VolumeMounts, k.MountPath+"\n")
			}
		}
	}
	for _, pod := range podlist.Items {
		if tool.CompareString(pod.Name, "^"+business.Name+`-[0-9]{1,}`) {
			business.Pod = append(business.Pod, PodInfo{
				PodName:      business.Name,
				PodIP:        pod.Status.PodIP,
				CpuUse:       getPodCpu(pod.Namespace, pod.Name),
				MemoryUse:    getPodMemory(pod.Namespace, pod.Name),
				RestartCount: getPodRestartCount(pod.Namespace, pod.Name),
				Status:       string(pod.Status.Phase),
				HealthCheck:  handleHealthPod(business.UserHealth_Path, business.UserHealth_port, pod.Status.PodIP),
			})
		}
	}
	return business
}

func getBusinessInfo(BusinnessList []BusinnessMetric) {
	// 获取集群所有pod
	podlist := kubernetes.GetPods()
	for index, business := range BusinnessList {
		if business.ControllerName == "deployment" {
			BusinnessList[index] = getDeploymentInfo(podlist, BusinnessList[index])
		} else if business.ControllerName == "statefulset" {
			BusinnessList[index] = getStatefulsetInfo(podlist, BusinnessList[index])
		} else if business.ControllerName == "daemonset" {
			BusinnessList[index] = getDaemonsetInfo(podlist, BusinnessList[index])
		}
	}
}

func TestDeployment(wg *sync.WaitGroup, resultChan chan interface{}) {
	if tool.Conf.Controller["businness"] != "true" {
		wg.Done()
		return
	}
	fmt.Println("--------------- 开始自定义检查 ----------")
	var BusinnessList []BusinnessMetric
	initBusinnessConfig(&BusinnessList)
	getBusinessInfo(BusinnessList)
	fmt.Println("--------------- 结束自定义检查 ----------")
	wg.Done()
	resultChan <- BusinnessList
}
