package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	id      int
	randnum int
}
type Result struct {
	task   Task
	result int
}

var tasks = make(chan Task, 10)
var results = make(chan Result, 10)

func process(num int) int {
	sum := 0
	for num != 0 {
		digit := num % 10
		sum += digit
		num /= 10
	}
	time.Sleep(2 * time.Second)
	return sum
}
func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		result := Result{task, process(task.randnum)}
		results <- result
	}
}
func createWorkerPool(numOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()
	close(results)
}
func allocate(numOfTasks int) {
	for i := 0; i < numOfTasks; i++ {
		randnum := rand.Intn(999)
		task := Task{i, randnum}
		tasks <- task
	}
	close(tasks)
}
func getResult(done chan bool) {
	for result := range results {
		fmt.Printf("Task id %d, randnum %d , sum %d\n", result.task.id, result.task.randnum, result.result)
	}
	done <- true
}
func main() {
	startTime := time.Now()
	numOfWorkers := 20
	numOfTasks := 100

	var done = make(chan bool)
	go getResult(done)
	go allocate(numOfTasks)
	go createWorkerPool(numOfWorkers)
	// 必须在allocate()和getResult()之后创建工作池
	<-done
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time taken ", diff.Seconds(), "seconds")
}
