package demo

//package main
//
//import (
//	"fmt"
//	"strings"
//
//	flag "github.com/spf13/pflag"
//)
//
//// 定义命令行参数对应的变量
//// 长短兼容 --name 或者 -n
//var cliName = flag.StringP("name", "n", "nick", "Input Your Name")
//
//// 长短兼容 --age -a
//var cliAge = flag.IntP("age", "a", 22, "Input Your Age")
//
//// 长短兼容 --gender -g
//var cliGender = flag.StringP("gender", "g", "male", "Input Your Gender")
//
//// 长命令  --corp
//var cliCorp = flag.String("corp", "hnkc", "Input Your Corp")
//
//func wordSepNormalizeFunc(f *flag.FlagSet, name string) flag.NormalizedName {
//	from := []string{"-", "_"}
//	to := "."
//	for _, sep := range from {
//		name = strings.Replace(name, sep, to, -1)
//	}
//	return flag.NormalizedName(name)
//}
//
//func main() {
//	flag.Lookup("phone").NoOptDefVal = "1234"
//
//	// 设置标准化参数名称的函数
//	flag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
//
//	// 为 age 参数设置 NoOptDefVal
//	flag.Lookup("age").NoOptDefVal = "25"
//
//	// 把 badflag 参数标记为即将废弃的，请用户使用 des-detail 参数
//	flag.CommandLine.MarkDeprecated("badflag", "please use --des-detail instead")
//	// 把 badflag 参数的 shorthand 标记为即将废弃的，请用户使用 des-detail 的 shorthand 参数
//	flag.CommandLine.MarkShorthandDeprecated("badflag", "please use -d instead")
//
//	// 在帮助文档中隐藏参数 gender
//	flag.CommandLine.MarkHidden("badflag")
//
//	// 把用户传递的命令行参数解析为对应变量的值
//	flag.Parse()
//
//	fmt.Println("name=", *cliName)
//	fmt.Println("age=", *cliAge)
//	fmt.Println("gender=", *cliGender)
//	fmt.Println("phone=", *phone)
//}
