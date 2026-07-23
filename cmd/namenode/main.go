package main

import (
	"log"
	"net"

	"github.com/Suraj-singh04/vaultfs/internal/namenode"
	pb "github.com/Suraj-singh04/vaultfs/proto/namenode"
	"google.golang.org/grpc"
)

func main() {
	node := namenode.NewNameNode()

	handler := namenode.NewGrpcHandler(node)

	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNameNodeServer(grpcServer, handler)

	log.Println("NameNode is running on port 50051...")

	grpcServer.Serve(listener)
}
