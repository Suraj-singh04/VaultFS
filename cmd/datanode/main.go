package main

import (
	"log"
	"net"
	"os"

	"github.com/Suraj-singh04/vaultfs/internal/datanode"
	pb "github.com/Suraj-singh04/vaultfs/proto/datanode"
	"google.golang.org/grpc"
)

func main() {
	id := os.Getenv("NODE_ID")
	baseDir := os.Getenv("BASE_DIR")

	node := datanode.NewDataNode(id, baseDir, 10*1024*1024*1024)

	handler := datanode.NewGrpcHandler(node)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("DataNode listening on: :50052")

	grpcServer := grpc.NewServer()
	pb.RegisterDataNodeServer(grpcServer, handler)

	grpcServer.Serve(listener)
}
