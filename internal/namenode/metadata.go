package namenode

import "time"
import "sync"

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
	ConfirmedNodes []string
}

type NameNode struct {
	mu sync.RWMutex
	DataNodes map[string]DataNode
	Files map[string]FileMetaData
	Chunks map[string]ChunkMetaData
}