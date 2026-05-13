package datanode

import "os"

type DataNode struct {
	ID               string
	BaseDir          string
	AvailableStorage int64
}

func NewDataNode(id string, baseDir string, availableStorage int64) *DataNode {
	err := os.MkdirAll(baseDir, 0755)

	if err != nil {
		panic("failed to create base directory for datanode: " + err.Error())
	}

	return &DataNode{
		ID:               id,
		BaseDir:          baseDir,
		AvailableStorage: availableStorage,
	}
}
