package datanode

import (
	"fmt"
	"io"

	pb "github.com/Suraj-singh04/vaultfs/proto/datanode"
)

type GrpcHandler struct {
	node *DataNode
	pb.UnimplementedDataNodeServer
}

func NewGrpcHandler(node *DataNode) *GrpcHandler {
	return &GrpcHandler{node: node}
}

func (h *GrpcHandler) StoreChunk(stream pb.DataNode_StoreChunkServer) error {
	var chunkID string
	var data []byte

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		chunkID = req.ChunkId
		data = append(data, req.Data...)
	}

	success := h.node.StoreChunk(chunkID, data)

	if !success {
		return fmt.Errorf("failed to store chunk %s", chunkID)
	}

	return stream.SendAndClose(&pb.StoreChunkResponse{Success: true})
}

func (h *GrpcHandler) FetchChunk(req *pb.RetrieveChunkRequest, stream pb.DataNode_FetchChunkServer) error {
	data, err := h.node.FetchChunk(req.ChunkId)
	if err != nil {
		return err
	}

	chunkSize := 1024 * 1024

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize

		if end > len(data) {
			end = len(data)
		}

		err := stream.Send(&pb.RetrieveChunkResponse{
			Data: data[i:end],
		})

		if err != nil {
			return err
		}
	}

	return nil

}
