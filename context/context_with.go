package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go Watch(ctx, "goroutine1")
	go Watch(ctx, "goroutine2")

	time.Sleep(6 * time.Second) // 让goroutine1和goroutine2执行6s
	fmt.Println("end watching!!!")
	cancel() // 通知goroutine1和goroutine2关闭
	time.Sleep(1 * time.Second)
}

func Watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s exit!\n", name) // 主goroutine调用cancel后，会发送一个信号到ctx.Done()这个channel，这里就会收到信息
			return
		default:
			fmt.Printf("%s watching...\n", name)
			time.Sleep(time.Second)
		}
	}
}
