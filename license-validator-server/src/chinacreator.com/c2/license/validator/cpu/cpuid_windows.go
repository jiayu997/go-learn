package cpu

import (
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"strings"
)


func CheckCPU(cpuid string) error {
	var CPUIDExist = false
	var infos, _ = cpu.Info()
	fmt.Println(infos)
	for _, cpuinfo := range infos {
		var PhysicalID, cpuid = strings.ToUpper(cpuinfo.PhysicalID), strings.ToUpper(cpuid)
		fmt.Println(PhysicalID)
		if strings.EqualFold(PhysicalID, cpuid) {
			CPUIDExist = true
			break
		}
	}
	if !CPUIDExist {
		return errors.New("硬件指标验证不通过.")
	}
	return nil
}
