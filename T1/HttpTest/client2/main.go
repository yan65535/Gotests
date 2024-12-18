package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
)

// 客户端
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080") //建立连接
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer conn.Close() // 关闭连接
	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := inputReader.ReadString('\n') // 读取用户输入
		str_ := "客户端2:" + input
		fmt.Println("客户端输入数据", str_, "，类型：", reflect.TypeOf(str_))
		inputInfo := strings.Trim(str_, "\r\n")
		if strings.ToUpper(inputInfo) == "Q" { // 如果输入q就退出
			return
		}
		_, err = conn.Write([]byte(inputInfo)) // 发送数据
		if err != nil {
			return
		}
		buf := [512]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("recv failed, err:", err)
			return
		}
		fmt.Println(string(buf[:n]))
	}
}
