package main

import (
	"context"
	"fmt"
)

func main() {
	Done := make(chan int64)
	go func() {

		var sum int64
		for i := 0; i < 10; i++ {
			sum += int64(i)
			fmt.Println(i, "==", sum)
		}
		Done <- sum
	}()
	go func() {

		var sum int64
		for i := 0; i < 10; i++ {
			sum += int64(i)
			fmt.Println(i, "--", sum)
		}
		Done <- sum
	}()
	context.Background()
	t := <-Done
	fmt.Println(t)
}
