package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 定义一个结构体 Counter
type Counter struct {
	count int32
}

// 使用原子操作增加计数
func (c *Counter) Increment() {
	atomic.AddInt32(&c.count, 1)
}

// 使用原子操作获取当前计数
func (c *Counter) Get() int32 {
	time.Sleep(3 * time.Second)
	return atomic.LoadInt32(&c.count)
}

func main() {
	var counter Counter
	var wg sync.WaitGroup
	numGoroutines := 100

	wg.Add(numGoroutines)

	// 启动多个协程来并发增加计数
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Println("=-=-")
	// 输出最终计数结果
	fmt.Println("Final Counter:", counter.Get())
}
