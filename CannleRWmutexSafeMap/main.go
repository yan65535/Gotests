package main

import (
	"fmt"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

type ISafeMap interface {
	Set(key string, value int)
	Get(key string) int
}

type safeMap struct {
	done chan struct{}
	m    map[string]int
}

func (m *safeMap) Set(key string, value int) {
	m.done <- struct{}{}
	m.m[key] = value
	<-m.done
}

func (m *safeMap) Get(key string) int {
	m.done <- struct{}{}
	v := m.m[key]
	<-m.done
	return v
}

func main() {
	var m = &safeMap{
		done: make(chan struct{}, 1),
		m:    make(map[string]int),
	}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			m.Set(strconv.Itoa(i), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			fmt.Println(m.Get(strconv.Itoa(i)))
		}
	}()
	wg.Wait()
}
