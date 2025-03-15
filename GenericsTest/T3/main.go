package main

import "fmt"

type AA struct {
	name string
	age  int
}

func main() {
	a := AA{
		name: "1",
		age:  23,
	}
	fmt.Printf("%+v\n", a)
}
