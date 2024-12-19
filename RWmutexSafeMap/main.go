package main

import (
	"fmt"
	"strconv"
	"sync"
)

var wg sync.WaitGroup
var rwMutex sync.RWMutex

func main() {
	m := make(map[string]int)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			rwMutex.Lock()
			m[strconv.Itoa(i)] += i
			rwMutex.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			rwMutex.RLock()
			fmt.Println(m[strconv.Itoa(i)])
			rwMutex.RUnlock()
		}
	}()
	wg.Wait()
}
