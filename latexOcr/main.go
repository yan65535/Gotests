package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "latexOcr/proto/latexOcr" // 替换为实际的包路径
)

func main() {
	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用 DialContext 连接到 gRPC 服务器
	conn, err := grpc.DialContext(ctx, "localhost:50052", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 创建客户端
	client := pb.NewLatexServiceClient(conn)

	// 读取图片文件
	imageData, err := ioutil.ReadFile("img/img.png")
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
