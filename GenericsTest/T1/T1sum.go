package main

import (
	"fmt"
	"reflect"
)

func sum[T int | float64](a, b T) T {
	fmt.Println(reflect.TypeOf(a), reflect.TypeOf(b))
	fmt.Println("hello")
	return a + b

}
func main() {
	fmt.Println(sum(4, 210))

}
