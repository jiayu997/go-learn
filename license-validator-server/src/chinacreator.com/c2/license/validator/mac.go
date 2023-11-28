package validator

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func CheckMac(macs []string) error {
	netInfos, _ := net.Interfaces()
	var macExist = false
	for _, netInfo := range netInfos {
		for _, mac := range macs {
			var formatMac = strings.Replace(mac, "-", ":", -1)
			formatMac = strings.ToLower(formatMac)

			if netInfo.HardwareAddr.String() == formatMac {
				fmt.Println("%s:%s",netInfo.HardwareAddr.String(),formatMac)
				return nil
			}
		}
	}
	if !macExist {
		return errors.New("硬件指标验证不通过..")
	}
	return nil
}
