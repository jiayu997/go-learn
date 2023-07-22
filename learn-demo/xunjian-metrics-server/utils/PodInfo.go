package utils

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/xuri/excelize/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// 目前暂时使用RESTCLIENT 客户端，因其为所有客户端的父类
type PodInfo struct {
	Namespace     string          //命名空间
	PodName       string          // pod name
	PodStatus     string          // pod status
	PodIp         string          // pod ip
	CpuUse        float64         // pod cpu useage
	MemoryUse     int64           // pod memory useage
	ContainerList []ContainerInfo //container list
}

type ContainerInfo struct {
	ContainerName    string
	ContainerImage   string // container image name
	CpuRe            string // pod cpu request  container1+container2+...containern
	MemoryRe         string // pod memory request  container1+container2+...cotainern
	CpuLimit         string // pod cpu limits
	MemoryLimit      string // pod momory limits
	HealthyStatus    bool   // pod healthy check status
	ContainerRestart int32  // pod restart total
}

func getPodMetrics(namespace string, podName string) (cpuQuantity float64, memQuantity int64) {
	config := GenerateConfig()
	mc, err := metrics.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
	//	podMetrics, err := mc.MetricsV1beta1().PodMetricses(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	podMetrics, err := mc.MetricsV1beta1().PodMetricses(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range podMetrics.Containers {
		if len(podMetrics.Containers) == 0 {
			cpuQuantity = 0
			memQuantity = 0
			break
		} else {
			cpuQuantity += container.Usage.Cpu().ToDec().AsApproximateFloat64()
			memQuantity += container.Usage.Memory().AsDec().UnscaledBig().Int64()
		}
	}
	// fmt.Println(decimal(cpuQuantity), memQuantity/(1024*1024))
	return decimal(cpuQuantity), memQuantity / (1024 * 1024)
}

func handleContainer(result *corev1.PodList, PodList chan PodInfo, wg *sync.WaitGroup) {
	// 用于处理pod metrics信息
	//	fmt.Println(getPodMetrics("harbor", "harbor-registry-7db8ddcb76-jch9c"))

	for _, pod := range result.Items {
		var tmp PodInfo
		var cpuQuantity float64
		var memQuantity int64

		// Pod 异常处理
		//fmt.Println(pod.Name, pod.Status.Conditions[0].Status, pod.Status.Conditions[1].Status, pod.Status.Conditions[2].Status, pod.Status.Conditions[3].Status)
		// pod.Status.Conditions[1].Status == "False"
		// fmt.Println(pod.Name, pod.Status.Conditions, "--------------")
		// fmt.Println(pod.Status.Conditions[1].Status)
		if pod.Status.Phase == "Failed" || pod.Status.Phase == "Unknown" || pod.Status.Phase == "Pending" {
			fmt.Printf("POD名称：%-40s\t状态：%-s\n", pod.Name, pod.Status.Phase)
			cpuQuantity, memQuantity = 0, 0
			tmp.PodStatus = "Failed"
		} else if (pod.Status.Phase == "Running" || pod.Status.Phase == "Succeeded") && pod.Status.Conditions[1].Status == "False" {
			//fmt.Println(pod.Status.Conditions[1].Status)
			fmt.Printf("POD名称：%-40s\t状态：Failed\n", pod.Name)
			cpuQuantity, memQuantity = 0, 0
			tmp.PodStatus = "Failed"
		} else { // Succeeded|Running
			// pod cpu and memory usage
			cpuQuantity, memQuantity = getPodMetrics(pod.Namespace, pod.Name)
			tmp.PodStatus = "Running"
		}
		// 写入Pod信息
		tmp.Namespace = pod.Namespace
		tmp.PodName = pod.Name
		tmp.PodIp = pod.Status.PodIP
		tmp.CpuUse = cpuQuantity
		tmp.MemoryUse = memQuantity
		Clen := len(pod.Spec.Containers)
		if Clen == 0 {
			tmp.ContainerList = make([]ContainerInfo, 1)
		} else {
			tmp.ContainerList = make([]ContainerInfo, len(pod.Spec.Containers))
		}
		for containerIndex := range pod.Spec.Containers {
			// 处于pending状态的pod，这里的信息为空切片
			if len(pod.Status.ContainerStatuses) == 0 {
				tmp.ContainerList[containerIndex].ContainerRestart = 0
				tmp.ContainerList[containerIndex].HealthyStatus = false
			} else {
				tmp.ContainerList[containerIndex].ContainerRestart = pod.Status.ContainerStatuses[containerIndex].RestartCount
				tmp.ContainerList[containerIndex].HealthyStatus = pod.Status.ContainerStatuses[containerIndex].Ready
			}
			tmp.ContainerList[containerIndex].ContainerName = pod.Spec.Containers[containerIndex].Name
			tmp.ContainerList[containerIndex].ContainerImage = pod.Spec.Containers[containerIndex].Image
			tmp.ContainerList[containerIndex].CpuRe = pod.Spec.Containers[containerIndex].Resources.Requests.Cpu().String()
			tmp.ContainerList[containerIndex].MemoryRe = pod.Spec.Containers[containerIndex].Resources.Requests.Memory().String()
			tmp.ContainerList[containerIndex].CpuLimit = pod.Spec.Containers[containerIndex].Resources.Limits.Cpu().String()
			tmp.ContainerList[containerIndex].MemoryLimit = pod.Spec.Containers[containerIndex].Resources.Limits.Memory().String()
		}
		// 当超过30个时会被阻塞
		PodList <- tmp
	}
	close(PodList)
	wg.Done()
}

// 用于生成POD metrics excel表格
func generateExcel(PodList <-chan PodInfo, wg *sync.WaitGroup, excelName string) {
	fmt.Println("----------------------------- POD状态检查 --------------------------------")
	// 打开工作簿
	excelF := openExcel(excelName)

	// 记录行号
	rowNum := 1

	// 获取流式写读器
	streamWriter, err := excelF.NewStreamWriter("POD巡检")
	if err != nil {
		log.Fatal(err.Error())
	}

	// 设置流式样式
	// styleLongWidth, err := excelF.NewStyle(&excelize.Style{Font: &excelize.Font{Color: "#777777"}})

	// 设置titile列宽度
	streamWriter.SetColWidth(1, 1, 20)
	streamWriter.SetColWidth(2, 2, 45)
	streamWriter.SetColWidth(3, 4, 15)
	streamWriter.SetColWidth(5, 5, 18)
	streamWriter.SetColWidth(6, 7, 15)
	streamWriter.SetColWidth(8, 8, 20)
	streamWriter.SetColWidth(9, 9, 70)
	streamWriter.SetColWidth(10, 11, 18)
	streamWriter.SetColWidth(12, 14, 20)

	// 设置title样式
	headerStyle := titileStyle(excelF)

	// 设置POD titile
	err = streamWriter.SetRow("A1", []interface{}{
		excelize.Cell{Value: "命名空间", StyleID: headerStyle},
		excelize.Cell{Value: "POD 名称", StyleID: headerStyle},
		excelize.Cell{Value: "POD 状态", StyleID: headerStyle},
		excelize.Cell{Value: "POD IP地址", StyleID: headerStyle},
		excelize.Cell{Value: "POD CPU使用核数", StyleID: headerStyle},
		excelize.Cell{Value: "POD内存使用", StyleID: headerStyle},
		excelize.Cell{Value: "POD 重启数", StyleID: headerStyle},
		excelize.Cell{Value: "Container名称", StyleID: headerStyle},
		excelize.Cell{Value: "Container镜像", StyleID: headerStyle},
		excelize.Cell{Value: "Container健康检查", StyleID: headerStyle},
		excelize.Cell{Value: "CPU最小请求", StyleID: headerStyle},
		excelize.Cell{Value: "内存最小请求", StyleID: headerStyle},
		excelize.Cell{Value: "CPU请求限制", StyleID: headerStyle},
		excelize.Cell{Value: "内存请求限制", StyleID: headerStyle},
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		pod, ok := <-PodList
		if ok {
			rowNum++
			//PodMetricsPrint(pod)
			var restartCount int
			var containerName, containerCPURequest, containerImage, containerMemoryRequst, containerCPULimit, containerMemoryLimit, containerHealthyStatus string
			// contaienr处理
			for i := 0; i < len(pod.ContainerList); i++ {
				if len(pod.ContainerList) == 1 {
					restartCount += int(pod.ContainerList[i].ContainerRestart)
					containerName += pod.ContainerList[i].ContainerName
					containerImage += pod.ContainerList[i].ContainerImage
					containerCPURequest += pod.ContainerList[i].CpuRe
					containerMemoryRequst += pod.ContainerList[i].MemoryRe
					containerCPULimit += pod.ContainerList[i].CpuLimit
					containerMemoryLimit += pod.ContainerList[i].MemoryLimit
					if pod.ContainerList[i].HealthyStatus {
						containerHealthyStatus += "Ready"
					} else {
						containerHealthyStatus += "NotReady"
					}
				} else {
					if i == len(pod.ContainerList)-1 {
						restartCount += int(pod.ContainerList[i].ContainerRestart)
						containerName += pod.ContainerList[i].ContainerName
						containerImage += pod.ContainerList[i].ContainerImage
						containerCPURequest += pod.ContainerList[i].CpuRe
						containerMemoryRequst += pod.ContainerList[i].MemoryRe
						containerCPULimit += pod.ContainerList[i].CpuLimit
						containerMemoryLimit += pod.ContainerList[i].MemoryLimit
						if pod.ContainerList[i].HealthyStatus {
							containerHealthyStatus += "Ready"
						} else {
							containerHealthyStatus += "NotReady"
						}
					} else {
						restartCount += int(pod.ContainerList[i].ContainerRestart)
						containerName += pod.ContainerList[i].ContainerName + "\r\n"
						containerImage += pod.ContainerList[i].ContainerImage + "\r\n"
						containerCPURequest += pod.ContainerList[i].CpuRe + "\r\n"
						containerMemoryRequst += pod.ContainerList[i].MemoryRe + "\r\n"
						containerCPULimit += pod.ContainerList[i].CpuLimit + "\r\n"
						containerMemoryLimit += pod.ContainerList[i].MemoryLimit + "\r\n"
						if pod.ContainerList[i].HealthyStatus {
							containerHealthyStatus += "Ready" + "\r\n"
						} else {
							containerHealthyStatus += "NotReady" + "\r\n"
						}
					}
				}
			}
			// 检查容器重启次数
			if restartCount > 10 {
				fmt.Printf("POD名称：%-40s\t重启次数：%v\n", pod.PodName, restartCount)
			}
			row := make([]interface{}, 14)
			row[0] = pod.Namespace
			row[1] = pod.PodName
			row[2] = pod.PodStatus
			row[3] = pod.PodIp
			row[4] = fmt.Sprintf("%.6f核", pod.CpuUse)
			row[5] = fmt.Sprintf("%dMi", pod.MemoryUse)
			row[6] = restartCount
			row[7] = containerName
			row[8] = containerImage
			row[9] = containerHealthyStatus
			row[10] = containerCPURequest
			row[11] = containerMemoryRequst
			row[12] = containerCPULimit
			row[13] = containerMemoryLimit
			cell, _ := excelize.CoordinatesToCellName(1, rowNum)
			streamWriter.SetRow(cell, []interface{}{
				excelize.Cell{Value: row[0], StyleID: contentNormalStyle(excelF)},
				excelize.Cell{Value: row[1], StyleID: contentNormalStyle(excelF)},
				excelize.Cell{Value: row[2], StyleID: contentStyle(excelF, pod.PodStatus, "Failed")},
				excelize.Cell{Value: row[3], StyleID: contentNormalStyle(excelF)},
				excelize.Cell{Value: row[4], StyleID: contentNormalStyle(excelF)},
				excelize.Cell{Value: row[5], StyleID: contentNormalStyle(excelF)},
				excelize.Cell{Value: row[6], StyleID: contentStyle(excelF, restartCount, 10)},
				excelize.Cell{Value: row[7], StyleID: contentNormalStyle(excelF)},
				excelize.Cell{Value: row[8], StyleID: contentNormalStyle(excelF)},
				excelize.Cell{Value: row[9], StyleID: contentStyle(excelF, containerHealthyStatus, "NotReady.*")}, //需要修复红色填充
				excelize.Cell{Value: row[10], StyleID: contentNormalStyle(excelF)},
				excelize.Cell{Value: row[11], StyleID: contentNormalStyle(excelF)},
				excelize.Cell{Value: row[12], StyleID: contentNormalStyle(excelF)},
				excelize.Cell{Value: row[13], StyleID: contentNormalStyle(excelF)},
			}, excelize.RowOpts{Height: 30, Hidden: false})
		} else {
			break
		}
	}
	fmt.Println("--------------------------------------------------------------------------")
	// 回刷缓存
	if err := streamWriter.Flush(); err != nil {
		log.Fatal(err.Error())
	}
	excelF.Save()
	excelF.Close()
	wg.Done()
}

func GetPodInfo(excelName string) {
	// 等待组，用于让主进程等待二个协程处理完信息
	var wg sync.WaitGroup
	wg.Add(2)

	config := GenerateConfig()

	// 配置API 路径
	config.APIPath = "api"

	// 配置GV版本
	config.GroupVersion = &corev1.SchemeGroupVersion

	// 配置数据序列化
	config.NegotiatedSerializer = scheme.Codecs

	// 定义接收变量值
	result := &corev1.PodList{}

	// 实例化rest client
	restClient, err := rest.RESTClientFor(config)
	if err != err {
		log.Fatal(err.Error())
	}

	err = restClient.Get().
		Resource("pods").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(result)

	if err != nil {
		log.Fatal(err.Error())
	}

	// 同时处理30个任务
	PodList := make(chan PodInfo, 30)

	// 用来处理container一般信息
	go handleContainer(result, PodList, &wg)

	// 用来生成execl表格
	go generateExcel(PodList, &wg, excelName)

	// 等待子协程
	wg.Wait()
}
