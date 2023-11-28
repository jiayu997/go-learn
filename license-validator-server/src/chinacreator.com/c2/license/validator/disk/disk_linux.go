package disk

import (
	"strings"
	"errors"
	"os/exec"
	liclog "chinacreator.com/c2/license/log"
)

func CheckDiskSerialNumber(diskSerialNumber string) error {
	out,err := exec.Command("sh","-c","fdisk -l |grep \"Disk identifier\" |awk {'print $3'} ").Output() ;
	if err != nil {
		liclog.Info.Println("未检测到硬盘相关数据...")
	}
	diskNums := strings.Split(string(out),"\n")

	for _, diskNum := range diskNums {
		diskNum := strings.TrimSpace(diskNum) ;
		if strings.EqualFold(diskNum,diskSerialNumber) {
			return nil
		}
	}
	return errors.New("硬件指标验证不通过.")
}