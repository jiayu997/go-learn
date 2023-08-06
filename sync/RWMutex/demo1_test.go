package RWMutex

import (
	"sync"
	"testing"
	"time"
)

func TestLockAndUnlock1(t *testing.T) {
	mutex := new(sync.RWMutex)
	t.Log("Main Lock the RWMutex")
	mutex.Lock()
	t.Log("Main The Lock is locked")

	ch := make([]chan int, 5)

	for i := 0; i < 5; i++ {
		// 无缓存
		ch[i] = make(chan int)
		go func(i int, c chan int) {
			t.Logf("Not Lock: %d", i)
			// 如果已经有锁了，则会被阻塞,直到被解锁
			mutex.Lock()
			t.Logf("Locked: %d", i)
			t.Logf("Unlock the lock: %d", i)
			mutex.Unlock()
			c <- i
		}(i, ch[i])
	}
	time.Sleep(time.Second * 2)
	t.Log("Main Unlock the lock")
	mutex.Unlock()
	time.Sleep(time.Second * 2)

	for _, c := range ch {
		<-c
	}
}
