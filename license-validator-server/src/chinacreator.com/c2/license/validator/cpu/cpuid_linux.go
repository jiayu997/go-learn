package cpu

import (
	"errors"
	"strings"
	"fmt"
	"os/exec"
)



func CheckCPU(cpuid string) error {
	//TODO gsub去除换行
	out,err := exec.Command("sh","-c","dmidecode -t 4 | grep ID |sort -u |awk -F': ' '{gsub(/ /,\"\",$2);print $2;}'").Output() ;
	if err != nil {
		fmt.Println("获取CPUID失败...",err)
	}
	//去除换行
	cpuids := strings.Split(string(out),"\n")

	for _, cpu := range cpuids {
		cpu := strings.TrimSpace(cpu) ;
		if strings.EqualFold(cpu,cpuid) {
			return nil
		}
	}
	return errors.New("硬件指标验证不通过.")
}
