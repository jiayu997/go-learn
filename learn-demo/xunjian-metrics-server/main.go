package main

import (
	"time"
	"xunjian/utils"
)

func main() {
	fileName := time.Now().Format("2006-01-02") + "巡检报告.xlsx"
	utils.CreateExcel(fileName)

	utils.GetPodInfo(fileName)
	utils.GetK8SResource(fileName)
	utils.GetNodeInfo(fileName)
	//utils.SendEmail(time.Now().Format("2006-01-02"))
}
