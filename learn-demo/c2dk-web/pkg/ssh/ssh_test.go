package ssh

import (
	"fmt"
	"testing"
)

func TestCmd(t *testing.T) {
	client, err := GetSSHClient("localhost", "boysandgirls", "22")
	if err != nil {
		t.Log(err.Error(), "client get error")
	}

	std, err := RunCmd(client, "ls /root")
	if err != nil {
		t.Log(err, "ssh command failed")
	}
	fmt.Println("command", string(std))
}
