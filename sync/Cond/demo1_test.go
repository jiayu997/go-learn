package Cond

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var done = false

func read(name string, c *sync.Cond) {
	// 上锁
	c.L.Lock()

	for !done {
		// wait 方法会先调用c.L.UnLock 然后 调用c.L.Lock
		// 在wait 期间，其他的go routing就可以去Lock
		c.Wait()
	}
	fmt.Printf("name: %s start reading\n", name)
	c.L.Unlock()
}

// write() 中的暂停了 1s，一方面是模拟耗时，另一方面是确保前面的 3 个 read 协程都执行到 c.Wait()中去
func write(name string, c *sync.Cond) {
	fmt.Printf("name: %s start writing\n", name)
	time.Sleep(time.Second * 2)
	c.L.Lock()
	done = true
	c.L.Unlock()
	fmt.Printf("name: %s wake all\n", name)
	c.Broadcast()
}

func TestNotifyAll(t *testing.T) {
	cond := sync.NewCond(&sync.Mutex{})
	go read("reader1", cond)
	go read("reader2", cond)
	go read("reader3", cond)

	write("writer", cond)
	time.Sleep(time.Second * 5)
}
