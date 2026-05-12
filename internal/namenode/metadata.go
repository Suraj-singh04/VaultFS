package namenode

import "time"

type DataNode struct {
	ID   string
	IP   string
	Port int
	AvailableStorage int64
	LastHeartbeat time.Time 
}

type FileMetaData struct {
	FileName string
	ChunkIDs []string
}

type ChunkMetaData struct {
	ChunkID string
	DataNodes []string
}

type NameNode struct {
	DataNodes map[string]DataNode
	Files map[string]FileMetaData
	Chunks map[string]ChunkMetaData
}