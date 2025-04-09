package main

import (
	"context"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"io/ioutil"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "latexOcr/proto/latexOcr" // 替换为实际的包路径
)

func main() {
	// 创建一个带有超时的上下文
	zkConn, _, err := zk.Connect([]string{"localhost:2181"}, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to ZooKeeper: %v", err)
	}
	defer zkConn.Close()
	pythonServiceAddr := ""
	// 获取Python服务的节点信息
	pythonServicePath := "/services/latex_ocr"
	children, _, err := zkConn.Children(pythonServicePath)
	if err != nil {
		log.Fatalf("Failed to get Python service nodes: %v", err)
	}
	if len(children) > 0 {
		nodePath := pythonServicePath + "/" + children[0]
		data, _, err := zkConn.Get(nodePath)
		if err != nil {
			log.Fatalf("Failed to get node data: %v", err)
		}
		pythonServiceAddr = string(data)
		fmt.Println(pythonServiceAddr)
		log.Printf("Found Python service at: %s", pythonServiceAddr)

		// 这里可以使用pythonServiceAddr来连接到Python服务
		// 例如：使用grpc.Dial来连接到Python服务
		//pythonConn, err := grpc.Dial(pythonServiceAddr, grpc.WithInsecure())
		//if err != nil {
		//    log.Fatalf("Failed to connect to Python service: %v", err)
		//}
		//defer pythonConn.Close()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用 DialContext 连接到 gRPC 服务器
	conn, err := grpc.DialContext(ctx, pythonServiceAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 创建客户端
	client := pb.NewLatexServiceClient(conn)

	// 读取图片文件
	imageData, err := ioutil.ReadFile("img/img_1.png")
	if err != nil {
		log.Fatalf("Failed to read image file: %v", err)
	}

	// 准备请求
	request := &pb.LatexRequest{
		ImageData: imageData,
	}

	// 调用服务方法
	response, err := client.RecognizeLatex(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to call RecognizeLatex: %v", err)
	}

	// 输出响应结果
	log.Printf("识别结果: %s", response.Result)
}
