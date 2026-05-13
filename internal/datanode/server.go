package datanode

import (
	"os"
	"path/filepath"
)

func (d *DataNode) StoreChunk(chunkID string, data []byte) bool {
	path := filepath.Join(d.BaseDir, chunkID)

	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return false
	}

	return true
}

func (d *DataNode) FetchChunk(chunkID string) ([]byte, error) {
	path := filepath.Join(d.BaseDir, chunkID)

	return os.ReadFile(path)
}

func (d *DataNode) DeleteChunk(chunkID string) bool {
	path := filepath.Join(d.BaseDir, chunkID)

	err := os.Remove(path)
	if err != nil {
		return false
	}

	return true
}
