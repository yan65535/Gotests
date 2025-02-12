package main

import "fmt"

type iPeople interface {
	say(string)
}
type people struct {
	name string
	age  int
}

func (p people) say(a string) {
	fmt.Println(a)
}

func main() {

	a := &people{
		name: "ywh",
		age:  20,
	}
	a.say("hello")

}
