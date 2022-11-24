package main

import (
	"catLog"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"proto/msg"

	"google.golang.org/grpc"
)

type server struct {
	msg.HelloServer
}

func (s *server) SayHello(context.Context, *msg.HelloRequest) (*msg.HelloReply, error) {
	catLog.Log("grpc hello")
	return &msg.HelloReply{}, nil
}

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	msg.RegisterHelloServer(grpcServer, &server{})
	grpcServer.Serve(lis)

}
