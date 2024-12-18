package main

import "fmt"

func main() {
	Done := make(chan int64)
	go func() {
		var sum int64
		for i := 0; i < 10; i++ {
			sum += int64(i)
		}
		Done <- sum
	}()
	t := <-Done
	fmt.Println(t)
}
