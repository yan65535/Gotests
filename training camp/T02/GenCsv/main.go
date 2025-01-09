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

var wg sync.WaitGroup
var mutex sync.Mutex

var client = &http.Client{
	Timeout: 20 * time.Second, // Client为线程安全 直接全局定义，避免资源消耗
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

	// os.Create 默认会覆盖原文件；如果不存在就会新建
	goodFile := mustFile(os.Create("good.csv"))
	defer goodFile.Close()

	badFile := mustFile(os.Create("bad.csv"))
	defer badFile.Close()

	badCsv := csv.NewWriter(badFile)
	defer badCsv.Flush() // 从缓存中刷新到badFile。

	goodCsv := csv.NewWriter(goodFile)
	defer goodCsv.Flush()

	goodCsv.Write(records[0])
	badCsv.Write(records[0])

	fmt.Println("总记录数:", len(records))

	for i, record := range records {
		if i == 0 {
			continue // 跳过表头
		}

		wg.Add(1)
		go func(record []string) {
			defer wg.Done()
			if checkUrl(record[4]) {
				addGood(record, goodCsv)
			} else {
				addBad(record, badCsv)
			}
		}(record)
	}

	wg.Wait() // 等待所有协程完成
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

func addGood(record []string, goodW *csv.Writer) {
	mutex.Lock()
	goodW.Write(record)
	defer mutex.Unlock()
	fmt.Println("有效链接:", record)

}

func addBad(record []string, badW *csv.Writer) {
	mutex.Lock()
	badW.Write(record)
	defer mutex.Unlock()
	fmt.Println("无效链接:", record)
}
