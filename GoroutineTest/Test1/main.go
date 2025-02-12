package main

import (
	"fmt"
	"sync"
)

//sync.WaitGroup：适用于等待多个 goroutine 处理完任务。
//channel：适用于 goroutine 之间的通信，数据流动和协作。

var wg sync.WaitGroup

func pA() {

	fmt.Println("A")

}
func pB() {

	fmt.Println("B")

}

func main() {
	Done := make(chan bool, 2)
	go func() {
		fmt.Println("A")
		pA()
		Done <- true
	}()
	go func() {
		fmt.Println("B")
		pB()
		Done <- true
	}()
	<-Done //阻塞 等待接信号
	<-Done //同上
	fmt.Println("Done")
}
