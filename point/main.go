package main

import (
	"fmt"
	"time"
)

func modify(slice []int) {
	slice[0] = 100
	slice = append(slice, 200)
	fmt.Println(slice)
}

func checkType(T interface{}) {
	switch T.(type) {
	case int:
		fmt.Println("int")
	case float64:
		fmt.Println("float64")
	case string:
		fmt.Println("string")
	default:
		fmt.Println("not a type")
	}

}

func main() {
	s := []int{1, 2, 3}
	modify(s)
	fmt.Println(s)
	checkType(s)
	ch := make(chan int, 2)
	//go func() {
	ch <- 1
	ch <- 2
	close(ch)
	//}()

	time.Sleep(time.Second)
	for v := range ch {
		fmt.Print(v, " ")
	}

}
