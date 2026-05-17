package clientapi

import (
	"context"
	"fmt"

	pbDataNode "github.com/Suraj-singh04/vaultfs/proto/datanode"
	pbNameNode "github.com/Suraj-singh04/vaultfs/proto/namenode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Coordinator struct {
	nameNodeClient  pbNameNode.NameNodeClient
	dataNodeClients map[string]pbDataNode.DataNodeClient
}

func NewCoordinator(nameNodeAddr string) (*Coordinator, error) {
	conn, err := grpc.NewClient(nameNodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	NameNodeClient := pbNameNode.NewNameNodeClient(conn)

	return &Coordinator{
		nameNodeClient:  NameNodeClient,
		dataNodeClients: make(map[string]pbDataNode.DataNodeClient),
	}, nil
}

func (c *Coordinator) sendChunkToNode(ctx context.Context, nodeID, nodeAddr, chunkID string, data []byte) error {
	nodeClient, exists := c.dataNodeClients[nodeID]

	if !exists {
		conn, err := grpc.NewClient(nodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		nodeClient = pbDataNode.NewDataNodeClient(conn)
		c.dataNodeClients[nodeID] = nodeClient
	}

	stream, err := nodeClient.StoreChunk(ctx)
	if err != nil {
		return err
	}

	chunkSize := 1024 * 1024
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		err := stream.Send(&pbDataNode.StoreChunkRequest{
			ChunkId: chunkID,
			Data:    data[i:end],
		})
		if err != nil {
			return err
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("datanode rejected chunk %s", chunkID)
	}

	return nil
}
