package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	url := "https://baidu.com"

	// 自定义 HTTP 客户端
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置超时时间为 5 秒
	}

	// 发出 GET 请求
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// 打印响应状态和内容
	fmt.Println("Status:", resp.Status)
	fmt.Println("Body:", string(body))
}
