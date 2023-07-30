package metrics

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
	"xunjian-prometheus-server/tool"

	"github.com/olekukonko/tablewriter"
)

func podTableView(PodList []Pod) {
	var errorPodList []Pod
	// 筛选异常pod
	for _, pod := range PodList {
		if !tool.CompareString(pod.Status, "Running|Succeeded") {
			errorPodList = append(errorPodList, pod)
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	if len(errorPodList) == 0 {
		// "2006-01-02 15:04:05 MST Mon"
		table.SetHeader([]string{"Pod巡检结果", "时间"})
		tmp := []string{"正常", time.Now().Format("2006-01-02 15:04")}
		table.Append(tmp)
		table.Render()
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Pod命名空间", "Pod名称", "Pod状态", "时间"})
		for _, pod := range errorPodList {
			table.Rich([]string{pod.NameSpace, pod.Name, pod.Status, time.Now().Format("2006-01-02 15:04")}, []tablewriter.Colors{
				tablewriter.Colors{},
				tablewriter.Colors{},
				colorSet(pod.Status, "Running|Succeeded"),
				tablewriter.Colors{},
			})
		}
		table.Render()
	}
	fmt.Printf("\n")
}

// 从状态和重启次数判断
func nodeTableView(NodeList []Node) {
	//	var errorNodeList []Node
	//	// 筛选异常Node
	//	//fmt.Println(NodeList)
	//	for _, node := range NodeList {
	//		if !node.Schedulable {
	//			//fmt.Println("unschduler")
	//			errorNodeList = append(errorNodeList, node)
	//			continue
	//		}
	//		if node.Ready == "NotReady" {
	//			//fmt.Println("notready")
	//			errorNodeList = append(errorNodeList, node)
	//			continue
	//		}
	//		if node.CpuUsePercent >= 85.0 {
	//			//fmt.Println("cpu")
	//			errorNodeList = append(errorNodeList, node)
	//			continue
	//		}
	//		if node.DiskUsePercent >= 85.0 {
	//			//fmt.Println("disk")
	//			errorNodeList = append(errorNodeList, node)
	//			continue
	//		}
	//		if node.MemoryUsePercent >= 85.0 {
	//			//fmt.Println("memory")
	//			errorNodeList = append(errorNodeList, node)
	//			continue
	//		}
	//		if node.OverFlow {
	//			//fmt.Println("overflow")
	//			errorNodeList = append(errorNodeList, node)
	//			continue
	//		}
	//	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	if len(NodeList) == 0 {
		table.SetHeader([]string{"Node巡检结果", "时间"})
		tmp := []string{"正常", time.Now().Format("2006-01-02 15:04")}
		table.Append(tmp)
		table.Render()
	} else {
		table.SetHeader([]string{"节点主机名", "节点IP", "节点是否可调度", "Kubelet状态", "节点CPU使用率", "节点CPU负荷过载", "节点内存使用率", "节点磁盘使用率", "时间"})
		for _, node := range NodeList {
			table.Rich([]string{node.NodeName, node.IP, strconv.FormatBool(node.Schedulable), node.Ready, strconv.FormatFloat(node.CpuUsePercent, 'f', 2, 64) + "%", strconv.FormatBool(node.OverFlow), strconv.FormatFloat(node.MemoryUsePercent, 'f', 2, 64) + "%", strconv.FormatFloat(node.DiskUsePercent, 'f', 2, 64) + "%", time.Now().Format("2006-01-02 15:04")}, []tablewriter.Colors{
				tablewriter.Colors{},
				tablewriter.Colors{},
				colorSet(strconv.FormatBool(node.Schedulable), "true"),
				colorSet(node.Ready, "^Ready$"),
				colorSet(node.CpuUsePercent, 85.0),
				colorSet(node.OverFlow, false),
				colorSet(node.MemoryUsePercent, 85.0),
				colorSet(node.DiskUsePercent, 85.0),
				tablewriter.Colors{},
			})
		}
		table.Render()
	}
	fmt.Printf("\n")
}

func harborTableView(HarborList []HarborMetric) {
	//	var errorHarborList []HarborMetric
	//
	//	for _, harbor := range HarborList {
	//		if harbor.Status != "Ok" {
	//			errorHarborList = append(errorHarborList, harbor)
	//			continue
	//		}
	//		if harbor.TestImageStatus != "Ok" {
	//			errorHarborList = append(errorHarborList, harbor)
	//			continue
	//		}
	//	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	if len(HarborList) == 0 {
		table.SetHeader([]string{"Harbor巡检结果", "时间"})
		table.Append(
			[]string{"正常", time.Now().Format("2006-01-02 15:04")},
		)
		table.Render()
	} else {
		table.SetHeader([]string{"业务名称", "Harbor IP/端口", "Harbor 账号/密码", "Harbor组件健康检查", "Harbor镜像Pull测试", "CPU使用率", "内存使用率", "磁盘使用率", "时间"})
		for _, harbor := range HarborList {
			table.Rich([]string{
				harbor.TypeName,
				harbor.IP + ":" + harbor.Port,
				harbor.Username + "/" + harbor.Password,
				harbor.Status,
				harbor.ImageStatus,
				strconv.FormatFloat(harbor.CpuUsePercent, 'f', 2, 64) + "%",
				strconv.FormatFloat(harbor.MemoryUsePercent, 'f', 2, 64) + "%",
				strconv.FormatFloat(harbor.DiskUsePercent, 'f', 2, 64) + "%",
				time.Now().Format("2006-01-02 15:04")},
				[]tablewriter.Colors{
					tablewriter.Colors{},
					tablewriter.Colors{},
					tablewriter.Colors{},
					colorSet(harbor.Status, "Ok"),
					colorSet(harbor.ImageStatus, "Ok"),
					colorSet(harbor.CpuUsePercent, 85.0),
					colorSet(harbor.MemoryUsePercent, 85.0),
					colorSet(harbor.DiskUsePercent, 85.0),
					tablewriter.Colors{},
				})
		}
		table.Render()
	}
	fmt.Printf("\n")
}

func storageTableView(storage Storage) {
	var errorPvList []PvMetric
	var errorPvcList []PvcMetric
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	for _, pv := range storage.Pv {
		if pv.Status != "Bound" {
			errorPvList = append(errorPvList, pv)
		}
	}
	for _, pvc := range storage.Pvc {
		if pvc.Status != "Bound" {
			errorPvcList = append(errorPvcList, pvc)
		}
	}
	var max_length int
	var min_length int
	if len(errorPvList) < len(errorPvcList) {
		max_length = len(errorPvcList)
		min_length = len(errorPvList)
	} else {
		max_length = len(errorPvList)
		min_length = len(errorPvcList)
	}
	if max_length == 0 {
		table.SetHeader([]string{"PV巡检结果", "PVC巡检结果", "时间"})
		table.Append(
			[]string{"正常", "正常", time.Now().Format("2006-01-02 15:04")},
		)
		table.Render()
		fmt.Println("\n")
	} else {
		table.SetHeader([]string{"PV名称", "PV状态", "PVC名称", "PVC名称空间", "PVC状态", "时间"})
		for j := 0; j < max_length; j++ {
			if j < min_length {
				table.Rich([]string{
					errorPvList[j].Name,
					errorPvList[j].Status,
					errorPvcList[j].Name,
					errorPvcList[j].NameSpace,
					errorPvcList[j].Status,
					time.Now().Format("2006-01-02 15:04"),
				}, []tablewriter.Colors{
					tablewriter.Colors{},
					colorSet(errorPvList[j].Status, "Bound|Available"),
					tablewriter.Colors{},
					tablewriter.Colors{},
					colorSet(errorPvcList[j].Status, "Bound|Available"),
					tablewriter.Colors{},
				})
			} else {
				if len(errorPvList) == len(errorPvcList) {
					table.Rich([]string{
						errorPvList[j].Name,
						errorPvList[j].Status,
						errorPvcList[j].Name,
						errorPvcList[j].NameSpace,
						errorPvcList[j].Status,
						time.Now().Format("2006-01-02 15:04"),
					}, []tablewriter.Colors{
						tablewriter.Colors{},
						colorSet(errorPvList[j].Status, "Bound|Available"),
						tablewriter.Colors{},
						tablewriter.Colors{},
						colorSet(errorPvcList[j].Status, "Bound|Available"),
						tablewriter.Colors{},
					})
				} else if len(errorPvList) > len(errorPvcList) {
					table.Rich([]string{
						errorPvList[j].Name,
						errorPvList[j].Status,
						"---",
						"---",
						"---",
						time.Now().Format("2006-01-02 15:04"),
					}, []tablewriter.Colors{
						tablewriter.Colors{},
						colorSet(errorPvList[j].Status, "Bound|Available"),
						tablewriter.Colors{},
						tablewriter.Colors{},
						tablewriter.Colors{},
						tablewriter.Colors{},
					})
				} else {
					table.Rich([]string{
						"---",
						"---",
						errorPvcList[j].Name,
						errorPvcList[j].NameSpace,
						errorPvcList[j].Status,
						time.Now().Format("2006-01-02 15:04"),
					}, []tablewriter.Colors{
						tablewriter.Colors{},
						tablewriter.Colors{},
						tablewriter.Colors{},
						tablewriter.Colors{},
						colorSet(errorPvcList[j].Status, "Bound|Available"),
						tablewriter.Colors{},
					})
				}
			}
		}
		table.Render()
		fmt.Println("\n")
	}
}

func mySQLTableView(MySQLList []MySQLMetric) {
	//	var errorMySQLList []MySQLMetric
	//	for _, mysql := range MySQLList {
	//		if mysql.Status != "Ok" {
	//			errorMySQLList = append(errorMySQLList, mysql)
	//		}
	//	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	if len(MySQLList) == 0 {
		table.SetHeader([]string{"MySQL巡检结果", "时间"})
		table.Append(
			[]string{"正常", time.Now().Format("2006-01-02 15:04")},
		)
		table.Render()
	} else {

		table.SetHeader([]string{"业务名称", "MySQL IP/端口", "MySQL 账号/密码", "MySQL检查", "时间"})
		for _, mysql := range MySQLList {
			table.Rich([]string{mysql.TypeName, mysql.IP + ":" + mysql.Port, mysql.Username + "/" + mysql.Password, mysql.Status, time.Now().Format("2006-01-02 15:04")}, []tablewriter.Colors{
				tablewriter.Colors{},
				tablewriter.Colors{},
				tablewriter.Colors{},
				colorSet(mysql.Status, "Ok"),
				tablewriter.Colors{},
			})
		}
		table.Render()
	}
	fmt.Printf("\n")
}

func pgsTableView(PgsList []PgsMetric) {
	//	var errorPgsList []PgsMetric
	//	for _, pgs := range PgsList {
	//		if pgs.Status != "Ok" {
	//			errorPgsList = append(errorPgsList, pgs)
	//		}
	//	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	if len(PgsList) == 0 {
		table.SetHeader([]string{"Postgres巡检结果", "时间"})
		table.Append(
			[]string{"正常", time.Now().Format("2006-01-02 15:04")},
		)
		table.Render()
	} else {
		table.SetHeader([]string{"业务名称", "Postgress IP/端口", "Postgress 账号/密码", "Postgress检查", "时间"})
		for _, pgs := range PgsList {
			table.Rich([]string{pgs.TypeName, pgs.IP + ":" + pgs.Port, pgs.Username + "/" + pgs.Password, pgs.Status, time.Now().Format("2006-01-02 15:04")}, []tablewriter.Colors{
				tablewriter.Colors{},
				tablewriter.Colors{},
				tablewriter.Colors{},
				colorSet(pgs.Status, "Ok"),
				tablewriter.Colors{},
			})
		}
		table.Render()
	}
	fmt.Printf("\n")
}

func redisTableView(RedisList []RedisMetric) {
	//	var errorRedisList []RedisMetric
	//	for _, redis := range RedisList {
	//		if redis.Status != "Ok" {
	//			errorRedisList = append(errorRedisList, redis)
	//		}
	//	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	if len(RedisList) == 0 {
		table.SetHeader([]string{"Redis巡检结果", "时间"})
		table.Append(
			[]string{"正常", time.Now().Format("2006-01-02 15:04")},
		)
		table.Render()
	} else {
		table.SetHeader([]string{"业务名称", "Redis IP/端口", "Redis 密码", "Redis检查", "时间"})
		for _, redis := range RedisList {
			table.Rich([]string{redis.TypeName, redis.IP + ":" + redis.Port, redis.Password, redis.Status, time.Now().Format("2006-01-02 15:04")}, []tablewriter.Colors{
				tablewriter.Colors{},
				tablewriter.Colors{},
				tablewriter.Colors{},
				colorSet(redis.Status, "Ok"),
				tablewriter.Colors{},
			})
		}
		table.Render()
	}
	fmt.Printf("\n")
}

func nfsTableView(NfsList []NfsMetric) {
	//	var errorNfsList []NfsMetric
	//	for _, nfs := range NfsList {
	//		if nfs.Status != "Ok" {
	//			errorNfsList = append(errorNfsList, nfs)
	//		}
	//	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	if len(NfsList) == 0 {
		table.SetHeader([]string{"NFS巡检结果", "时间"})
		table.Append(
			[]string{"正常", time.Now().Format("2006-01-02 15:04")},
		)
		table.Render()
	} else {
		table.SetHeader([]string{"业务名称", "NFS IP", "NFS 路径", "NFS检查", "CPU使用率", "内存使用率", "磁盘使用率", "时间"})
		//table.SetColumnColor(tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{
		//	tablewriter.Bold,
		//	tablewriter.FgRedColor,
		//})
		//table.SetColumnColor(tablewriter.Colors{}, tablewriter.Colors{}, colorSet())
		for _, nfs := range NfsList {
			table.Rich([]string{
				nfs.TypeName,
				nfs.IP,
				nfs.DataDir,
				nfs.Status,
				strconv.FormatFloat(nfs.CpuUsePercent, 'f', 2, 64) + "%",
				strconv.FormatFloat(nfs.MemoryUsePercent, 'f', 2, 64) + "%",
				strconv.FormatFloat(nfs.DiskUsePercent, 'f', 2, 64) + "%",
				time.Now().Format("2006-01-02 15:04")},
				[]tablewriter.Colors{
					tablewriter.Colors{},
					tablewriter.Colors{},
					tablewriter.Colors{},
					colorSet(nfs.Status, "Ok"),
					colorSet(nfs.CpuUsePercent, 85.0),
					colorSet(nfs.MemoryUsePercent, 85.0),
					colorSet(nfs.DiskUsePercent, 85.0),
					tablewriter.Colors{},
				})
		}
		table.Render()
	}
	fmt.Printf("\n")
}

func resourceTableView(resource ResourceMetric) {
	table := tablewriter.NewWriter(os.Stdout)
	// "2006-01-02 15:04:05 MST Mon"
	table.SetHeader([]string{"备份文件", "定时备份", "K8S组件状态", "DNS名称", "DNS解析", "时间"})
	table.SetRowLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.Rich([]string{
		strconv.FormatBool(resource.BackupFile),
		strconv.FormatBool(resource.Crontab),
		strconv.FormatBool(resource.Component),
		resource.DNS_SERVER,
		resource.DNS_Status,
		time.Now().Format("2006-01-02 15:04"),
	}, []tablewriter.Colors{
		colorSet(resource.BackupFile, true),
		colorSet(resource.Crontab, true),
		colorSet(resource.Component, true),
		tablewriter.Colors{
			tablewriter.ALIGN_CENTER,
		},
		colorSet(resource.DNS_Status, "Ok"),
		tablewriter.Colors{
			tablewriter.ALIGN_CENTER,
		},
	})
	table.Render()
	fmt.Println("\n")
}

func businessTableView(BusinessList []BusinnessMetric) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"名称", "副本", "宿主机网络", "配置健康检查", "CPU/内存需求", "CPU/内存限制", "CPU/内存使用", "挂载点", "PodStatus", "重启次数", "健康检查结果"})
	table.SetRowLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	for _, business := range BusinessList {
		var tmp_cpu_memory_request, tmp_cpu_memory_limit, tmp_mounts, tmp_pod_name, tmp_pod_status, tmp_pod_health_check, tmp_pod_resource_use, tmp_pod_restart_total string
		// 及时没节点没有配置，client-go默认会返回0
		//fmt.Println(len(business.CpuRequest), len(business.CpuLimit), len(business.MemoryRequest), len(business.MemoryLimit))
		//fmt.Println(business.CpuRequest, business.CpuLimit, business.MemoryRequest, business.MemoryLimit)
		var container_len int = len(business.CpuLimit)
		//fmt.Println(container_len, "len-len")
		for k := 0; k < container_len; k++ {
			if container_len == 1 {
				tmp_cpu_memory_request = business.CpuRequest[0] + "/" + business.MemoryRequest[0]
				tmp_cpu_memory_limit = business.CpuLimit[0] + "/" + business.MemoryLimit[0]
			} else {
				tmp_cpu_memory_request = tmp_cpu_memory_request + business.CpuRequest[k] + "/" + business.MemoryRequest[k] + "\n"
				tmp_cpu_memory_limit = tmp_cpu_memory_limit + business.CpuLimit[k] + "/" + business.MemoryLimit[k] + "\n"
			}
		}
		for _, h := range business.VolumeMounts {
			tmp_mounts = tmp_mounts + h + "\n"
		}
		for _, f := range business.Pod {
			tmp_pod_name += f.PodName + "\n" //pod name
			tmp_pod_status += f.Status + "\n"
			tmp_pod_health_check += strconv.FormatBool(f.HealthCheck) + "\n"
			tmp_pod_resource_use += strconv.FormatFloat(f.CpuUse, 'f', 2, 64) + "/" + strconv.Itoa(int(f.MemoryUse)/1048576) + "Mi" + "\n"
			tmp_pod_restart_total += strconv.Itoa(int(f.RestartCount)) + "\n"
		}
		table.Rich([]string{
			business.TypeName,
			strconv.Itoa(business.Replicas),
			strconv.FormatBool(business.HostNetwork),
			strconv.FormatBool(business.HealthCheck),
			tmp_cpu_memory_request,
			tmp_cpu_memory_limit,
			tmp_pod_resource_use,
			tmp_mounts,
			tmp_pod_status,
			tmp_pod_restart_total,
			tmp_pod_health_check,
		}, []tablewriter.Colors{
			tablewriter.Colors{},
			colorSet(1, business.Replicas),
			colorSet(strconv.FormatBool(business.HostNetwork), "true"),
			colorSet(strconv.FormatBool(business.HealthCheck), "true"),
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
		})
	}
	table.Render()
	fmt.Println("\n")
}

// 设置单元格样式
func colorSet(content, flag interface{}) []int {
	// 判断二个是否类型一致,如果没找到返回默认样式
	//fmt.Println(content, flag)
	if reflect.TypeOf(&content).Kind() != reflect.TypeOf(&flag).Kind() {
		return tablewriter.Colors{}
	}
	switch content.(type) {
	case float64:
		//fmt.Println(reflect.TypeOf(content), reflect.TypeOf(flag), "----------------")
		value, _ := content.(float64)
		pattern, _ := flag.(float64)
		//fmt.Println("float64", value, pattern)
		if value >= pattern {
			return tablewriter.Colors{
				tablewriter.Bold,
				tablewriter.FgRedColor,
				tablewriter.ALIGN_CENTER,
			}
		} else {
			return tablewriter.Colors{}
		}
	case int:
		value, _ := content.(int)
		pattern, _ := flag.(int)
		//fmt.Println("int", value, pattern)
		if value >= pattern {
			return tablewriter.Colors{
				tablewriter.Bold,
				tablewriter.FgRedColor,
			}
		} else {
			return tablewriter.Colors{}
		}
	case string: // 如果内容与给出的不匹配，则红颜色
		value, _ := content.(string)
		pattern, _ := flag.(string)
		//fmt.Println("string", value, pattern)
		if tool.CompareString(value, pattern) {
			return tablewriter.Colors{
				tablewriter.ALIGN_CENTER,
			}
		} else {
			return tablewriter.Colors{
				tablewriter.Bold,
				tablewriter.FgRedColor,
			}
		}
	case bool: // 如果内容与给出的不匹配，则红颜色
		value, _ := content.(bool)
		pattern, _ := flag.(bool)
		//fmt.Println("bool")
		if value == pattern {
			//fmt.Println("1111")
			return tablewriter.Colors{
				tablewriter.ALIGN_CENTER,
			}
		} else {
			//fmt.Println("2222")
			return tablewriter.Colors{
				tablewriter.Bold,
				tablewriter.FgRedColor,
			}
		}
	}
	return tablewriter.Colors{}
}

func ShowTableView(resultChan chan interface{}) {
	for {
		result, ok := <-resultChan
		if !ok {
			break
		}
		switch result.(type) {
		case []Pod:
			podTableView(result.([]Pod))
		case []Node:
			nodeTableView(result.([]Node))
		case []HarborMetric:
			harborTableView(result.([]HarborMetric))
		case Storage:
			storageTableView(result.(Storage))
		case []MySQLMetric:
			mySQLTableView(result.([]MySQLMetric))
		case []PgsMetric:
			pgsTableView(result.([]PgsMetric))
		case []RedisMetric:
			redisTableView(result.([]RedisMetric))
		case []NfsMetric:
			nfsTableView(result.([]NfsMetric))
		case ResourceMetric:
			resourceTableView(result.(ResourceMetric))
		case []BusinnessMetric:
			businessTableView(result.([]BusinnessMetric))
		}
	}
}
