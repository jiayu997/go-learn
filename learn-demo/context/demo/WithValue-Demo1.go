package context

import (
	"context"
	"fmt"
	"time"
)

func son(ctx context.Context, msg chan int) {
	t := time.Tick(time.Second)
	for _ = range t {
		select {
		case m := <-msg:
			fmt.Printf("接收到值：%d\n", m)
		case <-ctx.Done():
			fmt.Println("子协程结束了", ctx.Value("name"))
		}
	}
}

func main() {
	ctx := context.WithValue(context.Background(), "name", "jiayu")
	ctx, clear := context.WithCancel(ctx)
	message := make(chan int)
	go son(ctx, message)
	for i := 0; i < 10; i++ {
		message <- i
	}
	clear()
	time.Sleep(time.Second)
	fmt.Println("主进程结束了", ctx.Value("name"))
}
