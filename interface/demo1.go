package main

import "fmt"

type HelloWorld interface {
	Hello()
}
type Person1 struct {
	Name       string
	Number     string
	HelloWorld //匿名字段
}
type Person2 struct {
	Name   string
	Number string
	tag    HelloWorld //非匿名字段
}
type hello struct {
}

func (h hello) Hello() {
	fmt.Println("hello")
}

func main() {
	h := hello{}
	p1 := &Person1{"DoveOne", "1", h}
	p1.Hello() //结构体内嵌接口时，匿名字段不可以直接引用该字段的方法

	p1.HelloWorld.Hello()
	p2 := &Person2{"DoveOne", "1", h}
	p2.tag.Hello() //非匿名字段必须指定字段名才能引用字段的方法
}
