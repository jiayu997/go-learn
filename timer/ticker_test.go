package timer

import (
	"testing"
	"time"
)

// 在创建Ticker时会指定一个时间，作为事件触发的周期。这也是Ticker与Timer的最主要的区别

func TestNewTicker(t *testing.T) {
	ticker := time.NewTicker(time.Second * 1)

	go func() {
		for range ticker.C {
			t.Log("ticker.C/s")
		}
	}()

	go func() {
		for range (*ticker).C {
			t.Log("(*ticker).C/s")
		}
	}()

	time.Sleep(time.Second * 5)
	ticker.Stop()
}

func TestTick(t *testing.T) {
	for {
		select {
		// 这里会不断生成 ticker，而且 ticker 会进行重新调度，造成泄漏
		case <-time.Tick(time.Second * 1):
			t.Log("run/s")
		}
	}
}
