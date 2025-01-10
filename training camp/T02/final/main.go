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

type Task struct {
	Record []string
	Index  int
}

type Result struct {
	Record []string
	IsGood bool
}

var client = &http.Client{
	Timeout: time.Second * 10,
}

var wg sync.WaitGroup

func main() {
	now := time.Now()
	defer func() {
		fmt.Printf("程序耗时:%v\n", time.Since(now))
	}()
	data := mustFile(os.Open("data.csv"))
	defer data.Close()
	dataT := csv.NewReader(data)
	records, err := dataT.ReadAll() //读取csv中所有数据
	if err != nil {
		panic(err)
	}
	goodFile := mustFile(os.Create("good.csv")) //创建有效连接文件
	defer goodFile.Close()

	badFile := mustFile(os.Create("bad.csv")) //创建无效连接csv文件
	defer badFile.Close()

	goodCsv := csv.NewWriter(goodFile) //创建有效链接写入文件对象
	defer goodCsv.Flush()
	badCsv := csv.NewWriter(badFile) //创建无效链接写入文件对象
	defer badCsv.Flush()
	fmt.Println("总记录数", len(records))

	taskChan := make(chan Task)
	resultChan := make(chan Result)
	workerCount := 1000

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(taskChan, resultChan)
	}
	var resultWg sync.WaitGroup
	resultWg.Add(1)

}
func mustFile(file *os.File, err error) *os.File {
	if err != nil {
		panic(err)
	}
	return file
}

func worker(ch <-chan Task, result chan<- Result) {

}
func process() {

}
func checkURL(url string) bool {
	url = strings.TrimSpace(url)                 //去除url中多余空格
	req, err := http.NewRequest("GET", url, nil) //创建新req对象
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求错误:%s,", url)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 400

}
