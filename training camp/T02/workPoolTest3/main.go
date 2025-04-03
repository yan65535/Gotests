package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type CsvProcessor struct {
	goodCsv *csv.Writer
	badCsv  *csv.Writer
	mutex   sync.Mutex
	wg      sync.WaitGroup
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
		log.Fatalf("读取CSV文件失败: %v", err)
	}

	goodFile := mustFile(os.Create("good.csv"))
	defer goodFile.Close()

	badFile := mustFile(os.Create("bad.csv"))
	defer badFile.Close()

	goodCsv := csv.NewWriter(goodFile)
	defer goodCsv.Flush()

	badCsv := csv.NewWriter(badFile)
	defer badCsv.Flush()

	processor := &CsvProcessor{
		goodCsv: goodCsv,
		badCsv:  badCsv,
	}

	goodCsv.Write(records[0])
	badCsv.Write(records[0])

	fmt.Println("总记录数:", len(records))

	// 创建工作池
	numWorkers := 10
	jobs := make(chan []string, len(records)-1)
	results := make(chan struct{}, len(records)-1)

	// 启动工作池中的工人
	for w := 0; w < numWorkers; w++ {
		go worker(processor, jobs, results)
	}

	// 分配任务
	go func() {
		for i, record := range records {
			if i == 0 {
				continue // 跳过表头
			}
			jobs <- record
		}
		close(jobs)
	}()

	// 等待所有任务完成
	for a := 0; a < len(records)-1; a++ {
		<-results
	}

	fmt.Println("所有任务完成")
}

func mustFile(file *os.File, err error) *os.File {
	if err != nil {
		log.Fatalf("打开文件失败: %v", err)
	}
	return file
}

func checkUrl(url string) bool {
	url = strings.TrimSpace(url)
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("创建请求错误: %s, 错误信息: %v\n", url, err)
		return false
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("请求错误: %s, 错误信息: %v\n", url, err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

func (p *CsvProcessor) addGood(record []string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.goodCsv.Write(record)
	fmt.Println("有效链接:", record)
}

func (p *CsvProcessor) addBad(record []string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.badCsv.Write(record)
	fmt.Println("无效链接:", record)
}

func worker(processor *CsvProcessor, jobs <-chan []string, results chan<- struct{}) {
	for record := range jobs {
		if checkUrl(record[4]) {
			processor.addGood(record)
		} else {
			processor.addBad(record)
		}
		results <- struct{}{}
	}
}
