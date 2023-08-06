package set

import (
	"kubernetes-src/staging/src/k8s.io/apimachinery/pkg/util/rand"
	"sync"
	"testing"
	"time"
)

type Computer interface {
	Add(item string)
}

type Compute struct {
	// 用于控制通知
	cond *sync.Cond

	// 数据集合
	Data Set
}

func NewCompute() *Compute {
	return &Compute{
		cond: sync.NewCond(&sync.Mutex{}),
		Data: NewSet(),
	}
}

func (c *Compute) Add(item string) {
	// 上锁
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	// 增加数据
	c.Data.Insert(item)
}

func (c *Compute) Read() string {
	// 上锁
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	// 当长度为0时，等待数据进来
	for c.Data.Len() == 0 {
		// wait first unlock and lock
		// 当我们调用signal和broadcast时这里会返回，如果长度为0,还是会来这里，所以这二个效果实际是一样的，只是看通知go routing的数量不一样而已
		c.cond.Wait()
	}

	// here must exist data
	data, _ := c.Data.PopAny()
	return data
}

func TestNotifyOne(t *testing.T) {
	com := NewCompute()
	AddFunc := func(i int, c *Compute) {
		r := rand.String(5)
		// 添加数据
		c.Add(r)
		t.Logf("Add thread: %d Add data: %s", i, r)
		// 通知某个go routing
		c.cond.Signal()
	}

	ReadFunc := func(i int, c *Compute) {
		// 读取数据
		t.Logf("Read thread: %d Get data: %s", i, c.Read())
	}

	// 启动5个携程专门用来写数据
	for i := 1; i <= 5; i++ {
		go func(i int, c *Compute) {
			// 一直添加数据
			for {
				AddFunc(i, c)
				time.Sleep(time.Second * 1)
			}
		}(i, com)
	}

	// 启动5个协程用来读数据
	for i := 1; i <= 5; i++ {
		go func(i int, c *Compute) {
			// 一直去获取数据
			for {
				ReadFunc(i, com)
			}
		}(i, com)
	}

	ch := make(chan struct{}, 0)
	<-ch
}

func TestNotifyAll(t *testing.T) {
	com := NewCompute()

	AddFunc := func(i int, c *Compute) {
		r := rand.String(5)
		// 添加数据
		c.Add(r)
		t.Logf("Add thread: %d Add data: %s", i, r)
		// 通知所有
		c.cond.Broadcast()
	}
	ReadFunc := func(i int, c *Compute) {
		// 读取数据
		t.Logf("Read thread: %d Get data: %s", i, c.Read())
	}

	// 启动2个携程专门用来写数据
	for i := 1; i <= 10; i++ {
		go func(i int, c *Compute) {
			// 一直添加数据
			for {
				AddFunc(i, c)
				time.Sleep(time.Second * 1)
			}
		}(i, com)
	}

	// 启动5个协程用来读数据
	for i := 1; i <= 5; i++ {
		go func(i int, c *Compute) {
			// 一直去获取数据
			for {
				ReadFunc(i, com)
			}
		}(i, com)
	}

	ch := make(chan struct{}, 0)
	<-ch
}
