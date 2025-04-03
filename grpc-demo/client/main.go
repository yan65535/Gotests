package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "grpc-demo/proto/protoFiles/hello" // 替换为实际的包路径
)

func main() {
	// 连接到服务端
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloServiceClient(conn)

	// 准备请求
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	name := "yhw"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
