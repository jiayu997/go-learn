package metrics

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"sync"
	"time"
	"xunjian-prometheus-server/kubernetes"
	"xunjian-prometheus-server/tool"
)

type ResourceMetric struct {
	BackupFile bool
	Crontab    bool
	Component  bool
	DNS_SERVER string
	DNS_Status string
}

func checkCrontab() int {
	cmd := exec.Command("/bin/bash", "-c", "crontab -l | grep k8s_backup.sh")

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
	}
	return 0
}

func checkBackupFile() int {
	cmd := exec.Command("/bin/bash", "-c", "ls -al /opt/backup | grep -E 'bak|backup' &>/dev/null")

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
	}
	return 0
}

func checkK8sCs() int {
	cmd := exec.Command("/bin/bash", "-c", "kubectl get cs | grep Unhealthy &>/dev/null")

	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return 0
			//return exitError.ExitCode()
		}
	}
	return 1

}
func getDnsIP(dnsName, ip string) []string {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 5 * time.Second,
			}
			return d.DialContext(ctx, "udp", ip+":53")
		},
	}
	ipList, err := r.LookupHost(context.Background(), dnsName)
	if err != nil {
		return nil
	}
	return ipList
}

func dnsCheck(dns *ResourceMetric) {
	result := kubernetes.GetDNSServer()
	ipList := getDnsIP("kubernetes.default.svc.cluster.local", result.Spec.ClusterIP)
	if len(ipList) == 0 {
		dns.DNS_SERVER = "kubernetes.default.svc.cluster.local"
		dns.DNS_Status = "Failed"
	} else {
		dns.DNS_SERVER = "kubernetes.default.svc.cluster.local"
		dns.DNS_Status = "Ok"
	}
}

func TestResource(wg *sync.WaitGroup, resultChan chan interface{}) {
	if tool.Conf.Controller["resource"] != "true" {
		wg.Done()
		return
	}
	fmt.Println("--------------- 开始备份/DNS检查 ----------")
	var resource ResourceMetric
	// 备份检查
	if checkCrontab() == 0 {
		resource.Crontab = true
	} else {
		resource.Crontab = false
	}
	if checkBackupFile() == 0 {
		resource.BackupFile = true
	} else {
		resource.BackupFile = false
	}
	if checkK8sCs() == 0 {
		resource.Component = true
	} else {
		resource.Component = false
	}
	// dns检查
	dnsCheck(&resource)
	fmt.Println("--------------- 结束备份/DNS检查 ----------")
	resultChan <- resource
	wg.Done()
}
