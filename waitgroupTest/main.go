package main

import (
	"sync"
	//"time"
)

const N = 10

var wg = &sync.WaitGroup{}

func main() {

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {

			println(i)
			defer wg.Done()
		}(i)
	}
	wg.Wait()

}
