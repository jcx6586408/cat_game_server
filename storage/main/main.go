package main

import (
	"catLog"
	"flag"
	"fmt"
	"log"
	"net"
	"proto/msg"
	"storage"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50054, "The server port")
)

func main() {
	flag.Parse()
	storage.StorageInstance.Run()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	msg.RegisterStorageServer(grpcServer, &storage.Storage{})
	grpcServer.Serve(lis)
	catLog.Log("存储监听")
}
