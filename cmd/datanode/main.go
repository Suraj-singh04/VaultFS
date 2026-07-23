package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/Suraj-singh04/vaultfs/internal/datanode"
	pb "github.com/Suraj-singh04/vaultfs/proto/datanode"
	pbNameNode "github.com/Suraj-singh04/vaultfs/proto/namenode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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



	go grpcServer.Serve(listener)

	nameNodeAddr := os.Getenv("NAMENODE_ADDR")

	conn, err := grpc.NewClient(nameNodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect to NameNode: %v", err)
	}

	nameNodeClient := pbNameNode.NewNameNodeClient(conn)

	ctx := context.Background()
	resp, err := nameNodeClient.RegisterNode(ctx, &pbNameNode.RegisterNodeRequest{
		Ip:             os.Getenv("SERVICE_ADDR"),
		Port:           50052,
		AvailableSpace: 10 * 1024 * 1024 * 1024,
	})
	if err != nil {
		log.Fatalf("Failed to register with NameNode: %v", err)
	}

	log.Printf("Registered with NameNode, assigned ID: %s", resp.NodeId)

	
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			nameNodeClient.Heartbeat(context.Background(), &pbNameNode.HeartbeatRequest{
				NodeId:         resp.NodeId,
				AvailableSpace: 10 * 1024 * 1024 * 1024,
			})
		}
	}()

	select {}
}
