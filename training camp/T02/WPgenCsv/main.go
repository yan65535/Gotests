package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var client = &http.Client{
	Timeout: 8 * time.Second, // Client为线程安全 直接全局定义，避免资源消耗
}

type Task struct {
	Record []string
	Index  int
}

type Result struct {
	Record []string
	IsGood bool
}

func main() {
	now := time.Now()
	defer func() {
		fmt.Printf("程序耗时： %v\n", time.Since(now))
	}()

	data := mustFile(os.Open("data.csv"))
	defer data.Close()
	dataT := csv.NewReader(data)
	records, err := dataT.ReadAll()
	if err != nil {
		panic(err)
	}

	goodFile := mustFile(os.Create("good.csv"))
	defer goodFile.Close()

	badFile := mustFile(os.Create("bad.csv"))
	defer badFile.Close()

	goodCsv := csv.NewWriter(goodFile)
	defer goodCsv.Flush()

	badCsv := csv.NewWriter(badFile)
	defer badCsv.Flush()

	goodCsv.Write(records[0])
	badCsv.Write(records[0])

	fmt.Println("总记录数:", len(records))

	// 创建任务和结果通道
	taskCh := make(chan Task, 1000)
	resultCh := make(chan Result, 1000)

	// 创建工作池
	workerCount := 100
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(taskCh, resultCh, &wg)
	}

	// 启动结果处理协程
	var resultWg sync.WaitGroup
	resultWg.Add(1)
	go processResults(resultCh, goodCsv, badCsv, &resultWg)

	// 分发任务
	for i, record := range records {
		if i == 0 {
			continue // 跳过表头
		}
		taskCh <- Task{Record: record, Index: i}
	}
	close(taskCh) // 关闭任务通道，通知工人退出

	// 等待所有工人完成
	wg.Wait()
	close(resultCh) // 关闭结果通道，通知结果处理退出

	// 等待结果处理完成
	resultWg.Wait()
}

func mustFile(file *os.File, err error) *os.File {
	if err != nil {
		panic(err)
	}
	return file
}

func checkUrl(url string) bool {
	url = strings.TrimSpace(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("创建请求错误: %s, 错误信息: %v\n", url, err)
		return false
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求错误: %s, 错误信息: %v\n", url, err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

func worker(taskCh <-chan Task, resultCh chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskCh {
		isGood := checkUrl(task.Record[4])
		resultCh <- Result{Record: task.Record, IsGood: isGood}
	}
}

func processResults(resultCh <-chan Result, goodCsv *csv.Writer, badCsv *csv.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	for result := range resultCh {
		if result.IsGood {
			goodCsv.Write(result.Record)
			fmt.Println("有效链接:", result.Record)
		} else {
			badCsv.Write(result.Record)
			fmt.Println("无效链接:", result.Record)
		}
	}
}
