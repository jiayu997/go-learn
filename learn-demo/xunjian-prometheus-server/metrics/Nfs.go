package metrics

import (
	"fmt"
	"log"
	"sync"
	"syscall"
	"xunjian-prometheus-server/tool"
)

type NfsMetric struct {
	TypeName         string
	IP               string
	DataDir          string
	Status           string
	CpuUsePercent    float64
	MemoryUsePercent float64
	DiskUsePercent   float64
}

func initNfsList() []NfsMetric {
	NfsList := make([]NfsMetric, 0)
	for _, Client := range tool.Conf.CheckList.Nfs {
		var tmp NfsMetric
		if !tool.CheckIP(Client["ip"]) || Client["datadir"] == "" {
			continue
		}
		tmp.TypeName = Client["type_name"]
		tmp.IP = Client["ip"]
		tmp.DataDir = Client["datadir"]
		NfsList = append(NfsList, tmp)
	}
	return NfsList
}

func mountDatadir(server *NfsMetric) {
	err := syscall.Mount(":"+server.DataDir, "/tmp", "nfs", 0, "timeo=2,retry=1,soft,nolock,addr="+server.IP) // 20=2s
	if err != nil {
		server.Status = "Failed"
	} else {
		server.Status = "Ok"
		err = syscall.Unmount("/tmp", 0)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getNfsMetric(NfsList []NfsMetric) {
	nodeList := getnodelist()
	for index, nfs := range NfsList {
		for _, k := range nodeList {
			if k.IP == nfs.IP {
				NfsList[index].CpuUsePercent = k.CpuUsePercent
				NfsList[index].MemoryUsePercent = k.MemoryUsePercent
				NfsList[index].DiskUsePercent = k.DiskUsePercent
			}
		}
	}
}

func TestNfs(wg *sync.WaitGroup, resultChan chan interface{}) {
	if tool.Conf.Controller["nfs"] != "true" {
		wg.Done()
		return
	}
	fmt.Println("--------------- 开始NFS检查 ----------")
	NfsList := initNfsList()
	for i := 0; i < len(NfsList); i++ {
		mountDatadir(&NfsList[i])
	}
	getNfsMetric(NfsList)
	fmt.Println("--------------- 结束NFS检查 ----------")
	resultChan <- NfsList
	wg.Done()
}
