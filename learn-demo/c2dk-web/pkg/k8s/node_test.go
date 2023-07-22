package k8s

import (
	"fmt"
	"testing"
)

func TestDeleteNode(t *testing.T) {
	fmt.Println(DeleteNode("192.168.0.111"))
}
