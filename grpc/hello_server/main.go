package main

import (
	"context"
	"fmt"
	"hello_server/pb"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	reply := "hello" + in.GetName()
	return &pb.HelloResponse{Reply: reply}, nil
}

func main() {
	// 启动服务
	l, err := net.Listen("tcp", ":8999")
	if err != nil {
		fmt.Println(err)
		return
	}

	//创建 rpc 服务
	s := grpc.NewServer()

	//注册服务
	pb.RegisterGreeterServer(s, &server{})

	//启动服务
	err = s.Serve(l)
	if err != nil {
		fmt.Println(err)
		return
	}
}
