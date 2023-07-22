package main

import (
	"context"
	"fmt"
	"time"
)

func son(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("son 超时了")
			return
		default:
			fmt.Println("son wait")
			time.Sleep(time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go son(ctx)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("father 超时了")
			goto END
		default:
			fmt.Println("father wait")
			time.Sleep(time.Second)
		}
	}
END:
}
