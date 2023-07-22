package demo

//package main
//
//import (
//	"flag"
//	"fmt"
//)
//
//func main() {
//	//定义命令行参数方式1
//	var name string
//	var age int
//	var married bool
//	flag.StringVar(&name, "name", "张三", "姓名")
//	flag.IntVar(&age, "age", 18, "年龄")
//	flag.BoolVar(&married, "married", false, "婚否")
//	// 打印默认--选项以及值
//	flag.PrintDefaults()
//
//	// 设置已注册的flag的值。相当于改了name的值
//	flag.Set("name", "lbw")
//
//	//解析命令行参数
//	flag.Parse()
//	fmt.Println(name, age, married)
//
//	//返回flag后的 参数内容
//	fmt.Println(flag.Args())
//
//	//返回flag后的 args参数个数
//	fmt.Println(flag.NArg())
//
//	//返回flag个数
//	fmt.Println(flag.NFlag())
//
//	f1 := flag.Lookup("name")
//	if f1 == nil {
//		fmt.Println("--name not register")
//	} else {
//		fmt.Println(f1.DefValue, f1.Name, f1.Usage, f1.Value)
//	}
//
//	f2 := flag.Lookup("corp")
//	if f2 == nil {
//		fmt.Println("--corp not register")
//	} else {
//		fmt.Println(f2.DefValue, f2.Name, f2.Usage, f2.Value)
//	}
//}
