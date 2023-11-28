package disk

import (
	liclog "chinacreator.com/c2/license/log"
	"errors"
	"github.com/shirou/gopsutil/disk"
	"strings"
)

func CheckDiskSerialNumber(diskSerialNumber string) error {
	diskInfos,_ := disk.IOCounters("SerialNumber")

	if len(diskInfos)==0{
		liclog.Info.Println("未检测到硬盘相关数据...")
		return nil
	}
	var serialNumExist = false
	for _,value := range diskInfos{
		if value.SerialNumber!="" && strings.EqualFold(diskSerialNumber,value.SerialNumber) {
			serialNumExist = true
		}
	}
	if !serialNumExist {
		return errors.New("硬件指标验证不通过.")
	}
	return nil
}
