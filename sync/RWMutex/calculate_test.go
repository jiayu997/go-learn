package RWMutex

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type calculate struct {
	rw   sync.RWMutex
	step int64
	sum  int64
}

func (cal *calculate) Add() {
	// 上锁
	cal.rw.Lock()
	// 临界区
	cal.sum = cal.sum + cal.step
	fmt.Printf("Add: %d Sum: %d\n", cal.step, cal.sum)
	// end 临界区
	cal.rw.Unlock()
	time.Sleep(time.Second * 1)
}

// 有读锁的存在，写锁会被阻塞,为防止读锁一直存在，这里设置sleep随机时间
func (cal *calculate) Read() {
	// sleep 1 - 5
	t := rand.Intn(3) + 1
	// 上读锁
	cal.rw.RLock()
	// 临界区
	fmt.Printf("Read: %d Sleep: %d\n", cal.sum, t)
	// end 临界区
	cal.rw.RUnlock()
	time.Sleep(time.Duration(t) * time.Second)
	fmt.Println("Sleep Done")
}

func TestCalculate(t *testing.T) {
	var wg sync.WaitGroup
	cal := calculate{
		rw:   sync.RWMutex{},
		step: 1,
		sum:  0,
	}

	// 写数据
	go func() {
		for {
			cal.Add()
		}
	}()

	// 十个携程读数据
	wg.Add(10)
	for i := 1; i <= 10; i++ {
		go func() {
			cal.Read()
			wg.Done()
		}()
	}
	wg.Wait()
}
