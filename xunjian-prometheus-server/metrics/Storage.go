package metrics

import (
	"fmt"
	"sync"
	"xunjian-prometheus-server/kubernetes"
	"xunjian-prometheus-server/tool"

	corev1 "k8s.io/api/core/v1"
)

type PvMetric struct {
	Name   string
	Status string
}

type PvcMetric struct {
	Name      string
	NameSpace string
	Status    string
}

type Storage struct {
	Pv  []PvMetric
	Pvc []PvcMetric
}

func handlerPvResult(storageList *Storage, pv_re *corev1.PersistentVolumeList) {
	for _, pv := range pv_re.Items {
		var tmp PvMetric
		tmp.Name = pv.Name
		tmp.Status = string(pv.Status.Phase)
		storageList.Pv = append(storageList.Pv, tmp)
	}
}

func handlerPvcResult(storageList *Storage, pvc_re *corev1.PersistentVolumeClaimList) {
	for _, pvc := range pvc_re.Items {
		var tmp PvcMetric
		tmp.Name = pvc.Name
		tmp.NameSpace = pvc.Namespace
		tmp.Status = string(pvc.Status.Phase)
		storageList.Pvc = append(storageList.Pvc, tmp)
	}
}

func TestStorage(wg *sync.WaitGroup, resultChan chan interface{}) {
	if tool.Conf.Controller["storage"] != "true" {
		wg.Done()
		return
	}
	fmt.Println("--------------- 开始PV&PVC检查 ----------")
	var storageList Storage
	handlerPvResult(&storageList, kubernetes.GetPVs())
	handlerPvcResult(&storageList, kubernetes.GetPVCs())
	fmt.Println("--------------- 结束PV&PVC检查 ----------")
	resultChan <- storageList
	wg.Done()
}
