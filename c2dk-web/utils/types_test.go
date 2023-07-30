package utils

import (
	"fmt"
	"testing"
)

func TestInitNetworkInterface(t *testing.T) {
	Amp.Master = "192.168.0.10"
	err := Amp.InitNetworkInterface()
	if err != nil {
		t.Log(err)
	}
	fmt.Println(Amp.NetworkInterface)
}
