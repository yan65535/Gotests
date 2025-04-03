package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "grpc-demo/proto/protoFiles/hello" // 替换为实际的包路径
)

// 实现 HelloService 接口
type server struct {
	pb.UnimplementedHelloServiceServer
}

// 实现 SayHello 方法
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Helloss, " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
