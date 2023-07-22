package main

import (
	"sync"
	"xunjian-prometheus-server/metrics"
	"xunjian-prometheus-server/tool"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(10)
	var resultChan = make(chan interface{}, 10)
	// 初始化配置文件
	tool.InitConfig()

	// 获取pod信息
	go metrics.GetPodInfo(&wg, resultChan)

	// 获取node信息
	go metrics.GetNodeInfo(&wg, resultChan)

	// 获取MySQL信息
	go metrics.TestMySQL(&wg, resultChan)

	// 获取Redis相关信息
	go metrics.TestRedis(&wg, resultChan)

	// 获取postgress信息
	go metrics.TestPgs(&wg, resultChan)

	// 获取NFS信息
	go metrics.TestNfs(&wg, resultChan)

	// 获取Harbor信息
	go metrics.TestHarbor(&wg, resultChan)

	// 获取备份信息与DNS信息
	go metrics.TestResource(&wg, resultChan)

	// 获取pv&pvc信息
	go metrics.TestStorage(&wg, resultChan)

	// 设置自定义组件检查
	go metrics.TestDeployment(&wg, resultChan)

	//等待协程完成，并关闭管道
	wg.Wait()
	close(resultChan)

	// 用于控制输出显示结果
	metrics.ShowTableView(resultChan)
}
