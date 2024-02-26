package main

import (
	"fmt"
	pb "github.com/piklen/pb/user"

	"google.golang.org/grpc"
	"net"
	"user/conf"
	"user/service"
)

const (
	//设置最大grpc的传输大小为50MB
	maxMsgSize = 50 * 1024 * 1024
)

func main() {
	conf.Init()
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer(
		//设置grpc传输接受和发送的最大传输量
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	) // 创建gRPC服务器
	pb.RegisterUserServiceServer(s, &service.Server{}) // 在gRPC服务端注册user服务
	fmt.Printf("创建gRPC服务器成功！！！！")
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
