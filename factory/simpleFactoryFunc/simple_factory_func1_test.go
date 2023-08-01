package simpleFactoryFunc

import (
	"fmt"
	"testing"
)

type Feature interface {
	Photo() string
	Phone() string
}

type Apple struct {
	Name string
	Type string
}

func (apple *Apple) Photo() string {
	return apple.Name + "/" + apple.Type + "take photo"
}

func (apple *Apple) Phone() string {
	return apple.Name + "/" + apple.Type + "make a phone"
}

// 由AppleFactory来创建对象
type AppleFactory interface {
	Create(name, ty string) *Apple
}

// 由它来实现Apple创建
type AppleFactoryStruct struct{}

func (ap *AppleFactoryStruct) Create(name, ty string) *Apple {
	return &Apple{
		Name: name,
		Type: ty,
	}
}

// 由它来实现Apple创建
func NewApple(name, ty string) *Apple {
	return &Apple{
		Name: name,
		Type: ty,
	}
}

type Huawei struct {
	Name string
	Type string
}

func (huawei *Huawei) Photo() string {
	return huawei.Name + "/" + huawei.Type + "take photo"
}

func (huawei *Huawei) Phone() string {
	return huawei.Name + "/" + huawei.Type + "make a phone"
}

// 由HuaweiFactory来创建对象
type HuaweiFactory interface {
	Create(name, ty string) *Huawei
}

// 由它来实现Huawei创建
type HuaweiFactoryStruct struct{}

func (hua *HuaweiFactoryStruct) Create(name, ty string) *Huawei {
	return &Huawei{
		Name: name,
		Type: ty,
	}
}

// 由它来实现Huawei创建
func NewHuawei(name, ty string) *Huawei {
	return &Huawei{
		Name: name,
		Type: ty,
	}
}

func TestNewMethod1(t *testing.T) {
	// 这种方式，每增加一个手机，这个手机都要实现Feature功能才行,以及一个Newxx()方法
	huawei := NewHuawei("jiayu", "Mate8 ")
	apple := NewApple("jiayu", "14Pro Max ")
	fmt.Println(huawei.Photo())
	fmt.Println(huawei.Phone())
	fmt.Println(apple.Photo())
	fmt.Println(apple.Phone())
}

func TestNewMethod2(t *testing.T) {
	// 这种方式，每增加一个手机,这个手机要实现Feature功能，同时要增加一个工厂结构体和方法(由它来创建手机)
	huaweiFactory := &HuaweiFactoryStruct{}
	huawei := huaweiFactory.Create("jiayu", "Mate8 ")
	fmt.Println(huawei.Photo())
	fmt.Println(huawei.Phone())

	appleFactory := &AppleFactoryStruct{}
	apple := appleFactory.Create("jiayu", "14Pro Max")
	fmt.Println(apple.Photo())
	fmt.Println(apple.Phone())
}
