//package main
//
//import (
//	"fmt"
//)
//
//func main() {
//	done := make(chan int, 5)
//	for i := 0; i < 5; i++ {
//		go func() {
//			fmt.Println(i)
//			done <- i
//		}()
//	}
//	for i := 0; i < 5; i++ {
//		<-done
//	}
//
//}

package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	//done := make(chan int, 5)

	for i := 0; i < 5; i++ {

		go func(j int) {

			defer wg.Done()

			fmt.Println(j)
			//done <- j
		}(i)
	}
	wg.Add(5)
	wg.Wait()

}
