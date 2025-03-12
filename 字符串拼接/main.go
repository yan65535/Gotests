package main

import (
	"fmt"
	"strings"
)

func main() {
	var builder strings.Builder
	builder.WriteString("hello world")
	builder.WriteString("a555")
	a := builder.String()

	fmt.Println(a)
}
