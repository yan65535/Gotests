package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan bool)

	go func() {
		time.Sleep(10 * time.Second) // 模拟长任务
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("Task completed")
	case <-time.After(5 * time.Second): // 阻塞等待，等待五秒 如果17行的<-done还没有解除阻塞 将会执行Timeout
		fmt.Println("Timeout")
	}
}
