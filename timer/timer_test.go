package timer

import (
	"fmt"
	"testing"
	"time"
)

func TestChannel(test *testing.T) {
	type Person struct {
		Name string
	}
	t1 := make(chan Person, 0)
	test.Logf("%T\n", t1) // chan timer.Person
	t2 := new(chan Person)
	test.Logf("%T\n", t2) // *chan timer.Person
	t2 = &t1
	go func() {
		t1 <- Person{Name: "jiayu"}
	}()
	go func() {
		fmt.Println(<-(*t2))
	}()
	time.Sleep(time.Second * 5)
}

func TestWaitChannel(test *testing.T) {
	timer := time.NewTimer(3 * time.Second)
	t := make(chan struct{}, 0)
	defer close(t)
	select {
	case <-t:
		fmt.Println("wait timeout")
	case <-timer.C:
		fmt.Println("timeout")
	}
}

// 其返回值代表定时器有没有超时：
// true：定时器超时前停止，后续不会再有事件发送
// false：定时器超时后停止
func TestStopTimer(t *testing.T) {
	timer1 := time.NewTimer(time.Second * 3)
	if timer1.Stop() {
		t.Log("timer1 stop success")
	} else {
		t.Log("timer1 stop failed")
	}
	timer2 := time.NewTimer(time.Second * 3)
	time.Sleep(time.Second * 5)

	if timer2.Stop() {
		t.Log("timer2 stop success")
	} else {
		t.Log("timer2 stop failed")
	}
}

func TestResetTimer(t *testing.T) {
	timer := time.NewTimer(time.Second * 3)
	<-timer.C
	t.Log("time out!")
	if timer.Stop() {
		t.Log("timer stop success")
	} else {
		t.Log("timer stop failed")
	}

	// 已经过期的定时器或已经停止的定时器，可以通过重置来重新激活
	// 重置的动作实质上是先停掉定时器，再启动。其返回值也即停掉计时器的返回值。
	// 将时间过期重置为5s
	if timer.Reset(time.Second * 5) {
		t.Log("not timeout")
	} else {
		t.Log("already timeout")
	}
}

// 有时我们就是想等指定的时间，没有需求提前停止定时器，也没有需求复用该定时器，那么可以使用匿名的定时器：
func TestAfter(t *testing.T) {
	t.Logf("%v", time.Now())
	t.Log("xx")
	<-time.After(time.Second * 1)
	t.Logf("%v", time.Now())
}

// time.After 这里会不断生成 timer，虽然最终会回收，但是会造成无意义的cpu资源消耗
// 因为：其底层会调用NewTimer(d).C
//
//	for {
//		select {
//			case <-time.After(time.Second * 1):
//			t.Log("one second")
//			}
//		}
//	}
func TestAfterLoop(t *testing.T) {
	timer := time.NewTimer(1 * time.Second)
	for {
		timer.Reset(time.Second * 1)
		select {
		case <-timer.C:
			t.Log("每秒执行一次")
		}
	}
}

// 我们可以使用 AfterFunc 更加简洁的实现延迟一个方法的调用：
// time.AfterFunc()是异步执行的，所以需要在函数最后sleep等待指定的协程退出，否则可能函数结束时协程还未执行
func TestAfterFunc(t *testing.T) {
	ch := make(chan struct{}, 0)
	t.Log("start to ch <- after 3 seconds")

	time.AfterFunc(time.Second*3, func() {
		t.Log("start to ch <- success")
		ch <- struct{}{}
	})
	<-ch
	t.Log("start to <- ch success")
}
