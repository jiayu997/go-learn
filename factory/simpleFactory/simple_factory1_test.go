package simpleFactory

import (
	"fmt"
	"testing"
)

type Printer interface {
	Print(name string) string
}

type CnPrinter struct{}

func (*CnPrinter) Print(name string) string {
	return fmt.Sprintf("cn: %s", name)
}

type EnPrinter struct{}

func (*EnPrinter) Print(name string) string {
	return fmt.Sprintf("en: %s", name)
}

// 简单工厂的优点是，简单，缺点嘛，如果具体产品扩产，就必须修改工厂内部，增加Case，一旦产品过多就会导致简单工厂过于臃肿
func TestNew1(t *testing.T) {
	var print Printer
	pName := "cn"
	switch pName {
	case "cn":
		print = new(CnPrinter)
	case "en":
		print = new(EnPrinter)
	default:
		print = new(CnPrinter)
	}
	fmt.Println(print.Print("chinese"))
}
