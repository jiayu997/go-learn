package utils

import (
	"context"
	"fmt"
	"log"
	"strconv"

	excelize "github.com/xuri/excelize/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	//	"k8s.io/client-go/rest"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type NodeList struct {
	NodeName             string
	NodeIP               string
	Schedulable          bool
	NodeStatus           string // Ready|DiskPressure|MemoryPressure|PIDPressure|NetworkUnavailable
	NodeArchitecture     string
	NodeKernelVersion    string
	NodeOperation        string
	NodeCpuUse           float64
	NodeMemoryUse        int64
	NodeCpuUsePercent    float64
	NodeMemoryUsePercent float64
}

func getNodeMetrics(nodeName string) (float64, int64) {
	config := GenerateConfig()
	mc, err := metrics.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
	nodeMetric, err := mc.MetricsV1beta1().NodeMetricses().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}
	return decimal(nodeMetric.Usage.Cpu().AsApproximateFloat64()), nodeMetric.Usage.Memory().AsDec().UnscaledBig().Int64() / (1024 * 1024)
	// fmt.Printf("%s\t%f核\t%dMi\n", nodeMetric.Name, nodeMetric.Usage.Cpu().AsApproximateFloat64(), nodeMetric.Usage.Memory().AsDec().UnscaledBig().Int64()/(1024*1024))
	//	nodeMetrics, err := mc.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	//
	//	for _, nodeMetric := range nodeMetrics.Items {
	//		fmt.Println(nodeMetric)
	//	}
}

func handleNodeExcel(result *corev1.NodeList, excelName string) {
	fmt.Println("----------------------------- NODE状态检查 -------------------------------")
	var normalFlag bool = true

	// 打开工作簿
	excelF := openExcel(excelName)

	createSheet(excelF, "NODE节点巡检")

	// 记录行号
	rowNum := 1

	// 获取流式写读器
	streamWriter, err := excelF.NewStreamWriter("NODE节点巡检")
	if err != nil {
		log.Fatal(err.Error())
	}

	// 设置titile列宽度
	streamWriter.SetColWidth(1, 1, 20)
	streamWriter.SetColWidth(2, 2, 20)
	streamWriter.SetColWidth(3, 4, 20)
	streamWriter.SetColWidth(5, 5, 20)
	streamWriter.SetColWidth(6, 7, 20)
	streamWriter.SetColWidth(8, 8, 20)
	streamWriter.SetColWidth(9, 9, 20)
	streamWriter.SetColWidth(10, 11, 20)

	// 设置title样式
	headerStyle := titileStyle(excelF)

	// 设置POD titile
	err = streamWriter.SetRow("A1", []interface{}{
		excelize.Cell{Value: "NODE节点", StyleID: headerStyle},
		excelize.Cell{Value: "节点IP", StyleID: headerStyle},
		excelize.Cell{Value: "节点是否可调度", StyleID: headerStyle},
		excelize.Cell{Value: "节点状态", StyleID: headerStyle},
		excelize.Cell{Value: "节点系统架构", StyleID: headerStyle},
		excelize.Cell{Value: "节点内核版本", StyleID: headerStyle},
		excelize.Cell{Value: "节点系统类型", StyleID: headerStyle},
		excelize.Cell{Value: "节点CPU使用核数", StyleID: headerStyle},
		excelize.Cell{Value: "节点内存使用(Mi)", StyleID: headerStyle},
		excelize.Cell{Value: "节点CPU使用率", StyleID: headerStyle},
		excelize.Cell{Value: "节点内存使用率", StyleID: headerStyle},
	}, excelize.RowOpts{Height: 30, Hidden: false})
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, nodeInfo := range result.Items {
		var tmp NodeList
		var tmpStatus string

		// 行号记录
		rowNum++

		tmp.NodeName = nodeInfo.Name
		tmp.NodeIP = nodeInfo.Status.Addresses[0].Address
		tmp.NodeArchitecture = nodeInfo.Status.NodeInfo.Architecture
		tmp.NodeKernelVersion = nodeInfo.Status.NodeInfo.KernelVersion
		tmp.NodeOperation = nodeInfo.Status.NodeInfo.OperatingSystem

		// 节点异常处理
		//		fmt.Println(nodeInfo.Status.Conditions[0].Status, nodeInfo.Status.Conditions[0].Type)
		//		fmt.Println(nodeInfo.Status.Conditions[1].Status, nodeInfo.Status.Conditions[1].Type)
		//		fmt.Println(nodeInfo.Status.Conditions[2].Status, nodeInfo.Status.Conditions[2].Type)
		//		fmt.Println(nodeInfo.Status.Conditions[3].Status, nodeInfo.Status.Conditions[3].Type)
		//		fmt.Println(nodeInfo.Status.Conditions[4].Status, nodeInfo.Status.Conditions[4].Type)

		// 节点异常处理逻辑
		if nodeInfo.Status.Conditions[4].Status == "True" {
			tmp.Schedulable = !nodeInfo.Spec.Unschedulable
			tmp.NodeCpuUse, tmp.NodeMemoryUse = getNodeMetrics(tmp.NodeName)
			num, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", tmp.NodeCpuUse*100/nodeInfo.Status.Capacity.Cpu().AsApproximateFloat64()), 64)
			tmp.NodeCpuUsePercent = num
			tmp.NodeMemoryUsePercent = float64((tmp.NodeMemoryUse * 1024 * 1024 * 100) / nodeInfo.Status.Capacity.Memory().AsDec().UnscaledBig().Int64())
			// 由于每一种状态都属于切片中的一个，所以需要做特殊处理
			var tmpNum int
			for _, value := range nodeInfo.Status.Conditions {
				if value.Status == "True" {
					tmpNum++
					if tmpNum == 1 {
						tmpStatus += string(value.Type)
					} else {
						tmpStatus += "|" + string(value.Type)
					}
				}
			}
			if tmp.NodeCpuUsePercent > 90 {
				fmt.Printf("节点：%-40s CPU使用率过高，当前使用率为：%v\n", tmp.NodeName, tmp.NodeCpuUsePercent)
			}
			if tmp.NodeMemoryUsePercent > 90 {
				fmt.Printf("节点：%-40s 内存使用率过高，当前使用率为：%v\n", tmp.NodeName, tmp.NodeMemoryUsePercent)
			}
			tmp.NodeStatus = tmpStatus
		} else {
			fmt.Printf("节点：%-45s 状态异常\n", tmp.NodeName)
			tmp.Schedulable = false
			tmp.NodeCpuUse, tmp.NodeMemoryUse = 0, 0
			tmp.NodeStatus = "NotReady"
			tmp.NodeCpuUsePercent = 0
			tmp.NodeMemoryUsePercent = 0
			normalFlag = false
		}
		cell, _ := excelize.CoordinatesToCellName(1, rowNum)
		err = streamWriter.SetRow(cell, []interface{}{
			excelize.Cell{StyleID: contentNormalStyle(excelF), Value: tmp.NodeName},
			excelize.Cell{StyleID: contentNormalStyle(excelF), Value: tmp.NodeIP},
			excelize.Cell{StyleID: contentStyle(excelF, tmp.Schedulable, "false"), Value: tmp.Schedulable},
			excelize.Cell{StyleID: contentStyle(excelF, tmp.NodeStatus, "NotReady"), Value: tmp.NodeStatus},
			excelize.Cell{StyleID: contentNormalStyle(excelF), Value: tmp.NodeArchitecture},
			excelize.Cell{StyleID: contentNormalStyle(excelF), Value: tmp.NodeKernelVersion},
			excelize.Cell{StyleID: contentNormalStyle(excelF), Value: tmp.NodeOperation},
			excelize.Cell{StyleID: contentNormalStyle(excelF), Value: fmt.Sprintf("%.6f核", tmp.NodeCpuUse)},
			excelize.Cell{StyleID: contentNormalStyle(excelF), Value: fmt.Sprintf("%dMi", tmp.NodeMemoryUse)},
			excelize.Cell{StyleID: contentStyle(excelF, tmp.NodeCpuUsePercent, 90.0), Value: fmt.Sprintf("%.2f%%", tmp.NodeCpuUsePercent)},
			excelize.Cell{StyleID: contentStyle(excelF, tmp.NodeMemoryUsePercent, 90.0), Value: fmt.Sprintf("%.2f%%", tmp.NodeMemoryUsePercent)},
		}, excelize.RowOpts{Height: 30, Hidden: false})
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	if normalFlag {
		fmt.Printf("%-45s\t%s\n", "NODE节点检查：", "正常")
	}
	fmt.Println("--------------------------------------------------------------------------")
	// 回刷缓存
	if err := streamWriter.Flush(); err != nil {
		log.Fatal(err.Error())
	}
	excelF.Save()
	excelF.Close()
}

func GetNodeInfo(excelName string) {
	config := GenerateConfig()

	config.APIPath = "api"

	// 无组名资源组
	config.GroupVersion = &corev1.SchemeGroupVersion

	config.NegotiatedSerializer = scheme.Codecs

	result := &corev1.NodeList{}

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = restClient.Get().
		Resource("nodes").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(result)

	if err != nil {
		log.Fatal(err)
	}
	handleNodeExcel(result, excelName)
}
