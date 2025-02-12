package main

import "fmt"

func main() {
	m := make(map[string]int)
	l := make([]string, 0, len(m))
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3
	m["d"] = 4
	m["e"] = 5
	for k, _ := range m {
		l = append(l, k)
	}
	for _, v := range l {
		fmt.Println(m[v])
	}
	fmt.Println(m)
	fmt.Println(l)
}
