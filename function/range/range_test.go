package _range

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

// 结果
//
//		0x140000b80a8
//		0x140000b80a8
//		0x140000b80a8
//		tmp1 1 &{Name:3 Age:3}
//		tmp1 2 &{Name:3 Age:3}
//		tmp1 3 &{Name:3 Age:3}
//		tmp2 1 &{Name:1 Age:1}
//		tmp2 2 &{Name:2 Age:2}
//		tmp2 3 &{Name:3 Age:3}
//	 遍历 map 为随机序输出，slice 为索引序输出
//	 range v 是值拷贝，且只会声明初始化一次
func TestRangeStruct(t *testing.T) {
	tmp1 := make(map[string]*Person)
	tmp2 := make(map[string]*Person)
	stu := []Person{
		{Name: "1", Age: 1},
		{Name: "2", Age: 2},
		{Name: "3", Age: 3},
	}

	// stu的地址是一个不变值，如果我们引用其地址，那么在遍历完成后，这个地址会指向stu的最后一个值,从而导致结果一致
	for _, stu := range stu {
		// 都指向了同一个stu的内存指针,因为 for range 中的 v 只会声明初始化一次,不会每次循环都初始化，最后赋值会覆盖前面的
		fmt.Printf("%p\n", &stu)

		// 直接取地址会导致结果最终指向最后一个结果, stu.Name 这里是值，不会每个结果一致
		// &stu 由于在for range期间地址不变，在遍历完成后，最终都会指向最后一个结果
		tmp1[stu.Name] = &stu

		// 这里的newStu会有一个新地址，其地址会指向当前遍历的stu
		newStu := stu
		tmp2[stu.Name] = &newStu
	}

	for i, v := range tmp1 {
		fmt.Printf("tmp1 %v %+v\n", i, v)
	}

	for i, v := range tmp2 {
		fmt.Printf("tmp2 %v %+v\n", i, v)
	}
}

// 结果
//
//	index: 0 value: &{Name:1 Age:1}
//	index: 1 value: &{Name:2 Age:2}
//	index: 2 value: &{Name:3 Age:3}
func TestRangeStruct2(t *testing.T) {
	tmp1 := make([]*Person, 0)
	stu := []Person{
		{Name: "1", Age: 1},
		{Name: "2", Age: 2},
		{Name: "3", Age: 3},
	}

	for index := range stu {
		tmp1 = append(tmp1, &stu[index])
	}
	for i, v := range tmp1 {
		fmt.Printf("index: %d value: %+v\n", i, v)
	}
}

// for range 等效于下面的for
//
//	range_test.go:80: value1: 0x1400000c0c0
//	range_test.go:84: value2: 0x1400000c0d8
//	range_test.go:80: value1: 0x1400000c0c0
//	range_test.go:84: value2: 0x1400000c0f0
//	range_test.go:80: value1: 0x1400000c0c0
//	range_test.go:84: value2: 0x1400000c108
func TestForStruct(t *testing.T) {
	tmp1 := make([]*Person, 0)
	tmp2 := make([]*Person, 0)

	stu := []Person{
		{Name: "1", Age: 1},
		{Name: "2", Age: 2},
		{Name: "3", Age: 3},
	}

	// 用来模拟for range 用的v，等效的
	{
		var value1 Person
		for_temp := stu
		len_temp := len(for_temp)
		for index_temp := 0; index_temp < len_temp; index_temp++ {
			value_temp := for_temp[index_temp]

			value1 = value_temp

			tmp1 = append(tmp1, &value1)
			t.Logf("value1: %p\n", &value1)

			value2 := value_temp
			tmp2 = append(tmp2, &value2)
			t.Logf("value2: %p\n", &value2)
		}
	}
}

// 结果
//
//	值: 1 地址: 0x14000096670
//	值: 2 地址: 0x14000096690
//	值: 3 地址: 0x140000966b0
func TestForValue(t *testing.T) {
	t1 := []string{"1", "2", "3"}
	for _, v := range t1 {
		// value 每次都重新申请一次，地址都是不同的
		value := v
		fmt.Printf("值: %v 地址: %p\n", value, &value)
	}
}
