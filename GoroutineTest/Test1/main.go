package main

import (
	"fmt"
	"sync"
)

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
		Done <- true
	}()
	go func() {
		fmt.Println("B")
		Done <- true
	}()
	<-Done //阻塞 等待接信号
	<-Done //同上
	fmt.Println("Done")
}
