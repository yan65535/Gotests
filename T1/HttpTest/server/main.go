package main

// TCP编程 流程
//1.监听端口
//2.接收客户端请求建立链接
//3.创建goroutine处理链接。

// TCP server端
import (
	"bufio" //
	"fmt"
	"net"
)

func process(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn) //关闭连接
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:]) //读取数据
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("收到client端发来的数据", recvStr)
		conn.Write([]byte(recvStr)) //发送数据
	}
}
func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("listen failed,err", err)
		return
	}
	for {
		conn, err := listen.Accept() //建立链接
		if err != nil {
			fmt.Println("accept failed,err", err)
			continue
		}
		go process(conn)

	}
}
