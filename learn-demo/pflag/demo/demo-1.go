package demo

//package demo
//
//import (
//	"fmt"
//	"net"
//	"time"
//
//	"github.com/spf13/pflag"
//)
//
//func pflagDefine() {
//	//64位整数，不带单标志位的
//	var pflagint64 *int64 = pflag.Int64("number1", 1234, "this is int 64, without single flag")
//
//	//64位整数，带单标志位的
//	var pflagint64p *int64 = pflag.Int64P("number2", "n", 2345, "this is int 64, without single flag")
//
//	//这种可以把变量的定义和变量取值分开，适合于struct，全局变量等地方
//	var pflagint64var int64
//	pflag.Int64Var(&pflagint64var, "number3", 1234, "this is int64var")
//
//	//上面那一种的增加短标志位版
//	var pflagint64varp int64
//	pflag.Int64VarP(&pflagint64varp, "number4", "m", 1234, "this is int64varp")
//
//	//slice版本,其实是上面的增强版，但是支持多个参数，也就是导成一个slice
//	var pflagint64slice *[]int64 = pflag.Int64Slice("number5", []int64{1234, 3456}, "this is int64 slice")
//
//	//bool版本
//	var pflagbool *bool = pflag.Bool("bool", true, "this is bool")
//
//	//bytes版本
//	var pflagbyte *[]byte = pflag.BytesBase64("byte64", []byte("ea"), "this is byte base64")
//
//	//count版本
//	var pflagcount *int = pflag.Count("count", "this is count")
//
//	//duration版本
//	var pflagduration *time.Duration = pflag.Duration("duration", 10*time.Second, "this is duration")
//
//	//float版本
//	var pflagfloat *float64 = pflag.Float64("float64", 123.345, "this is florat64")
//
//	//IP版本
//	var pflagip *net.IP = pflag.IP("ip1", net.IPv4(192, 168, 1, 1), "this is ip, without single flag")
//
//	//mask版本
//	var pflagmask *net.IPMask = pflag.IPMask("mask", net.IPv4Mask(255, 255, 255, 128), "this is net mask")
//
//	//string版本
//	var pflagstring *string = pflag.String("string", "teststring", "this is string")
//
//	//uint版本
//	var pflaguint *uint64 = pflag.Uint64("uint64", 12345, "this is uint64")
//
//	pflag.Parse()
//	fmt.Println("number1 int64 is ", *pflagint64)
//	fmt.Println("number2 int64 is ", *pflagint64p)
//	fmt.Println("number3 int64var is ", pflagint64var)
//	fmt.Println("number4 int64varp is", pflagint64varp)
//	fmt.Println("number5 int64slice is", *pflagint64slice)
//	fmt.Println("bool is ", *pflagbool)
//	fmt.Println("byte64 is ", *pflagbyte)
//	fmt.Println("count is ", *pflagcount)
//	fmt.Println("duration is ", *pflagduration)
//	fmt.Println("float is ", *pflagfloat)
//	fmt.Println("ip1 net.ip is ", *pflagip)
//	fmt.Println("mask is %s", *pflagmask)
//	fmt.Println("string is ", *pflagstring)
//	fmt.Println("uint64 is ", *pflaguint)
//
//}
//
//func main() {
//	pflagDefine()
//}
//
