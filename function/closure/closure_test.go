package closure

import (
	"fmt"
	"testing"
)

// 结果
//
//	defer i:  0
//	defer i ptr: 0x14000016340
//	closure i:  100
//	closure i ptr: 0x14000016340
//
// 闭包是引用传递，所以 Golang 中使用匿名函数的时候要特别注意区分清楚引用传递和值传递。根据实际需要，我们在不需要引用传递的地方通过匿名函数参数赋值的方式实现值传递 [7]
func TestClosureType(t *testing.T) {
	i := 0

	// 闭包：i是引用传递(不是值传递,有地址的变量)
	defer func() {
		fmt.Println("closure i: ", i)
		fmt.Printf("closure i ptr: %v\n", &i)
	}()

	// 非闭包：i是值传递
	defer fmt.Printf("defer i ptr: %p\n", &i)
	defer fmt.Println("defer i: ", i)

	// 修改值
	i = 100
}

func TestClosureForRange(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}

	for _, v := range s {
		go func() {
			t.Log(v) //for_range会导致v的地址为最后一个而闭包是引用类型，把地址传递进去了，所以 v -> 5
		}()
	}

	for _, v := range s {
		// 这个v是值传递，相当于新开了一个变量空间，这个变量值等于每次传递过来的值
		go func(v int) {
			t.Log(v)
		}(v)
	}
}

func TestClosureDefer1(t *testing.T) {
	x, y := 1, 2

	defer func(a int) {
		fmt.Printf("x: %d,y: %d\n", a, y) // y 为闭包引用，最终结果为x：1，y：2 + 100
	}(x)

	x += 100
	y += 100
	fmt.Println(x, y)
}

func TestClosureDefer2(t *testing.T) {
	x, y := 1, 2

	func(a int) {
		fmt.Printf("x: %d,y: %d\n", a, y) // y 为闭包引用，最终结果为x：1，y：2
	}(x)

	x += 100
	y += 100
	fmt.Println(x, y)
}

// 结果
// 0
// 1
// 2
// 0x14000016320 3
// 0x14000016320 3
// 0x14000016320 3
func test1() []func() {
	var s []func()

	for i := 0; i < 3; i++ {
		// j := i  每次重新定义修复
		fmt.Println(i)
		s = append(s, func() {
			fmt.Println(&i, i) // i永远会等于3(最后一次循环是2，但是i++后就是3了)
		})
	}
	return s
}

func TestFuncList1(t *testing.T) {
	for _, f := range test1() {
		f()
	}
}
