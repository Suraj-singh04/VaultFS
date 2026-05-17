package clientapi

import (
	"context"
	"io"

	pbDataNode "github.com/Suraj-singh04/vaultfs/proto/datanode"
	pbNameNode "github.com/Suraj-singh04/vaultfs/proto/namenode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (c *Coordinator) FetchFile(fileName string) ([]byte, error) {
	ctx := context.Background()

	resp, err := c.nameNodeClient.GetFileLocation(ctx, &pbNameNode.GetFileLocationRequest{
		FileName: fileName,
	})

	if err != nil {
		return nil, err
	}

	chunkData := make([]byte, 0)

	for chunkId, nodeList := range resp.ChunkLocations {
		nodeId := nodeList.NodeIds[0]

		nodeClient, exists := c.dataNodeClients[nodeId]

		if !exists {
			conn, err := grpc.NewClient(nodeId, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, err
			}
			nodeClient = pbDataNode.NewDataNodeClient(conn)
			c.dataNodeClients[nodeId] = nodeClient
		}

		stream, err := nodeClient.FetchChunk(ctx, &pbDataNode.RetrieveChunkRequest{
			ChunkId: chunkId,
		})

		if err != nil {
			return nil, err
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			chunkData = append(chunkData, resp.Data...)
		}
	}
	return chunkData, nil
}
