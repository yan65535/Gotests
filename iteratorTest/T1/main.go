package main

import "fmt"

func main() {
	n := 8
	for f := range Fibonacci(n) {
		fmt.Println(f)
	}
}
func Fibonacci(n int) func(yield func(int) bool) {
	a, b, c := 0, 1, 1
	return func(yield func(int) bool) {
		for range n {
			if !yield(a) {
				return
			}
			a, b = b, c
			c = a + b
		}
	}
}
