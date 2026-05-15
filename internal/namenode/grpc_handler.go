package namenode

import (
	"context"

	pb "github.com/Suraj-singh04/vaultfs/proto/namenode"
)

type GrpcHandler struct {
	node *NameNode
	pb.UnimplementedNameNodeServer
}

func NewGrpcHandler(node *NameNode) *GrpcHandler {
	return &GrpcHandler{node: node}
}

func (h *GrpcHandler) RegisterNode(ctx context.Context, req *pb.RegisterNodeRequest) (*pb.RegisterNodeResponse, error) {

	id := h.node.RegisterNode(req.Ip, int(req.Port), req.AvailableSpace)

	return &pb.RegisterNodeResponse{
		NodeId: id,
	}, nil
}

func (h *GrpcHandler) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {

	alive := h.node.Heartbeat(req.NodeId, req.AvailableSpace)

	return &pb.HeartbeatResponse{Alive: alive}, nil
}

func (h *GrpcHandler) GetFileLocations(ctx context.Context, req *pb.GetFileLocationRequest) (*pb.GetFileLocationResponse, error) {

	locations := h.node.GetFileLocations(req.FileName)

	protoLocations := make(map[string]*pb.NodeList)

	for chunkId, nodeIds := range locations {
		protoLocations[chunkId] = &pb.NodeList{NodeIds: nodeIds}
	}

	return &pb.GetFileLocationResponse{
		ChunkLocations: protoLocations,
	}, nil
}

func (h *GrpcHandler) AllocateChunks(ctx context.Context, req *pb.AllocateChunksRequest) (*pb.AllocateChunksResponse, error) {
	chunkLocations := h.node.AllocateChunk(req.FileName, req.ChunkIds)

	protoLocations := make(map[string]*pb.NodeList)

	for chunkId, nodeIds := range chunkLocations {
		protoLocations[chunkId] = &pb.NodeList{NodeIds: nodeIds}
	}

	return &pb.AllocateChunksResponse{
		ChunkLocations: protoLocations,
	}, nil
}

func (h *GrpcHandler) ConfirmChunk(ctx context.Context, req *pb.ConfirmChunkRequest) (*pb.ConfirmChunkResponse, error) {
	confirmed := h.node.ConfirmChunk(req.ChunkId, req.NodeId)
	return &pb.ConfirmChunkResponse{Success: confirmed}, nil
}
