package main

import "fmt"

// defer 先进后出特性
func main() {
	list_ := make([]int64, 5)
	for i, i2 := range list_ {
		defer func() {
			fmt.Println(i, i2)
		}()
	}
}
