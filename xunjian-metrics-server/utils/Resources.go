package utils

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"sync"
	"time"

	excelize "github.com/xuri/excelize/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

func getPvcInfo(excelF *excelize.File, wg *sync.WaitGroup) {
	config := GenerateConfig()

	config.APIPath = "api"

	config.GroupVersion = &corev1.SchemeGroupVersion

	config.NegotiatedSerializer = scheme.Codecs

	result := &corev1.PersistentVolumeClaimList{}

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = restClient.Get().
		Resource("PersistentVolumeClaims").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(result)
	if err != nil {
		log.Fatal(err.Error())
	}

	var rowNum int = 1
	excelF.SetCellValue("K8S资源巡检", "C1", "PVC名称")
	excelF.SetCellValue("K8S资源巡检", "D1", "PVC异常")

	for _, pvc := range result.Items {
		if pvc.Status.Phase != "Bound" {
			rowNum++
			fmt.Printf("PVC名称：%-40s\t状态：%v\n", pvc.Name, pvc.Status.Phase)
			excelF.SetCellValue("K8S资源巡检", "C"+fmt.Sprintf("%d", rowNum), pvc.Name)
			excelF.SetCellValue("K8S资源巡检", "D"+fmt.Sprintf("%d", rowNum), pvc.Status.Phase)
			excelF.SetCellStyle("K8S资源巡检", "D"+fmt.Sprintf("%d", rowNum), "D"+fmt.Sprintf("%d", rowNum), contentRedStyle(excelF))
		}
	}
	if rowNum == 1 {
		fmt.Printf("%-45s\t%s\n", "PVC状态：", "正常")
	}
	wg.Done()
}

func getPvInfo(excelF *excelize.File, wg *sync.WaitGroup) {
	config := GenerateConfig()

	config.APIPath = "api"

	config.GroupVersion = &corev1.SchemeGroupVersion

	config.NegotiatedSerializer = scheme.Codecs

	result := &corev1.PersistentVolumeList{}

	var rowNum int = 1

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = restClient.Get().
		Resource("PersistentVolumes").
		VersionedParams(&metav1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(result)
	if err != nil {
		log.Fatal(err.Error())
	}
	excelF.SetCellValue("K8S资源巡检", "A1", "PV名称")
	excelF.SetCellValue("K8S资源巡检", "B1", "PV异常")

	for _, pv := range result.Items {
		if pv.Status.Phase != "Bound" {
			rowNum++
			fmt.Printf("PV名称：%-40s\t状态：%v\n", pv.Name, pv.Status.Phase)
			excelF.SetCellValue("K8S资源巡检", "A"+fmt.Sprintf("%d", rowNum), pv.Name)
			excelF.SetCellValue("K8S资源巡检", "B"+fmt.Sprintf("%d", rowNum), pv.Status.Phase)
			excelF.SetCellStyle("K8S资源巡检", "B"+fmt.Sprintf("%d", rowNum), "B"+fmt.Sprintf("%d", rowNum), contentRedStyle(excelF))
		}
	}
	if rowNum == 1 {
		fmt.Printf("%-45s\t%s\n", "PV状态：", "正常")
	}
	wg.Done()
}

// 获取k8s apiserver dns信息
func getApiServerDNS(excelF *excelize.File, wg *sync.WaitGroup) {
	config := GenerateConfig()

	config.APIPath = "api"

	config.GroupVersion = &corev1.SchemeGroupVersion

	config.NegotiatedSerializer = scheme.Codecs

	result := &corev1.Service{}

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = restClient.Get().
		Resource("services").
		Namespace("kube-system").
		Name("kube-dns").
		VersionedParams(&metav1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(result)
	if err != nil {
		log.Fatal(err.Error())
	}

	excelF.SetCellValue("K8S资源巡检", "E1", "DNS检查")
	ipList := dnsCheck("kubernetes.default.svc.cluster.local", result.Spec.ClusterIP)
	if len(ipList) == 0 {
		fmt.Printf("%-45s\t%s\n", "K8S集群DNS状态：", "异常")
		excelF.SetCellValue("K8S资源巡检", "E2", "DNS异常")
		excelF.SetCellStyle("K8S资源巡检", "E2", "E2", contentRedStyle(excelF))
	} else {
		fmt.Printf("%-45s\t%s\n", "K8S集群DNS状态：", "正常")
		excelF.SetCellValue("K8S资源巡检", "E2", "DNS正常")
		excelF.SetCellStyle("K8S资源巡检", "E2", "E2", contentNormalStyle(excelF))
	}
	wg.Done()
}

func dnsCheck(dnsName string, ip string) []string {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 10 * time.Second,
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

func getBackupInfo(excelF *excelize.File, wg *sync.WaitGroup) {
	excelF.SetCellValue("K8S资源巡检", "F1", "元数据备份检查")
	excelF.SetCellValue("K8S资源巡检", "G1", "备份文件检查")
	if checkCrontab() != 0 {
		fmt.Printf("%-45s\t%s\n", "K8S原数据备份状态：", "异常")
		excelF.SetCellValue("K8S资源巡检", "F2", "异常")
		excelF.SetCellStyle("K8S资源巡检", "F2", "F2", contentRedStyle(excelF))
	} else {
		fmt.Printf("%-45s\t%s\n", "K8S原数据备份状态：", "正常")
		excelF.SetCellValue("K8S资源巡检", "F2", "正常")
		excelF.SetCellStyle("K8S资源巡检", "F2", "F2", contentNormalStyle(excelF))
	}
	if checkBackupFile() != 0 {
		fmt.Printf("%-45s\t%s\n", "K8S备份文件状态：", "异常")
		excelF.SetCellValue("K8S资源巡检", "G2", "异常")
		excelF.SetCellStyle("K8S资源巡检", "G2", "G2", contentRedStyle(excelF))
	} else {
		fmt.Printf("%-45s\t%s\n", "K8S备份文件状态：", "正常")
		excelF.SetCellValue("K8S资源巡检", "G2", "正常")
		excelF.SetCellStyle("K8S资源巡检", "G2", "G2", contentNormalStyle(excelF))
	}
	wg.Done()
}

func GetK8SResource(excelName string) {
	fmt.Println("----------------------------- K8S资源检查 --------------------------------")
	fmt.Println("-----------------------  检查列表：PVC/PV/DNS/备份状态  -------------------")
	// open file
	excelF := openExcel(excelName)
	var wg sync.WaitGroup
	wg.Add(4)

	createSheet(excelF, "K8S资源巡检")

	// 生成dns巡检
	go getApiServerDNS(excelF, &wg)

	// 生成pv巡检
	go getPvInfo(excelF, &wg)

	// 生成pvc巡检
	go getPvcInfo(excelF, &wg)

	// 生成集群备份巡检
	go getBackupInfo(excelF, &wg)

	wg.Wait()
	excelF.Save()
	excelF.Close()
	fmt.Println("--------------------------------------------------------------------------")
}
