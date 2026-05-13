package namenode

import "time"

func (n *NameNode) RegisterNode(ip string, port int, availableStorage int64) string {
	n.mu.Lock()
	defer n.mu.Unlock()

	id := "xyz"

	n.DataNodes[id] = DataNode{
		ID:               id,
		IP:               ip,
		Port:             port,
		AvailableStorage: availableStorage,
		LastHeartbeat:    time.Now(),
	}

	return id
}

func (n *NameNode) Heartbeat(id string, availableStorage int64) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	node, exists := n.DataNodes[id]
	if !exists {
		return false
	}

	node.LastHeartbeat = time.Now()
	node.AvailableStorage = availableStorage
	n.DataNodes[id] = node

	return true
}

func (n *NameNode) AllocateChunk(fileName string, chunkIDs []string) map[string][]string {
	n.mu.Lock()
	defer n.mu.Unlock()

	healthyNodes := []DataNode{}
	for _, node := range n.DataNodes {
		if time.Since(node.LastHeartbeat) < 15*time.Second && node.AvailableStorage > 0 {
			healthyNodes = append(healthyNodes, node)
		}
	}

	if len(healthyNodes) < 3 {
		return nil
	}

	chunkLocations := make(map[string][]string)

	for _, chunkID := range chunkIDs {
		for i := 0; i < 3; i++ {
			node := healthyNodes[i]
			chunkLocations[chunkID] = append(chunkLocations[chunkID], node.ID)
			node.AvailableStorage -= 64 * 1024 * 1024
			n.DataNodes[node.ID] = node
		}

		n.Chunks[chunkID] = ChunkMetaData{
			ChunkID:   chunkID,
			DataNodes: chunkLocations[chunkID],
		}
	}
	n.Files[fileName] = FileMetaData{
		FileName: fileName,
		ChunkIDs: chunkIDs,
	}

	return chunkLocations
}

func (n *NameNode) ConfirmChunk(chunkID string, dataNodeID string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	chunkMeta, exists := n.Chunks[chunkID]
	if !exists {
		return false
	}

	for _, nodeID := range chunkMeta.DataNodes {
		if nodeID == dataNodeID {
			chunkMeta.ConfirmedNodes = append(chunkMeta.ConfirmedNodes, dataNodeID)
			n.Chunks[chunkID] = chunkMeta
			return true
		}
	}

	return false
}

func (n *NameNode) GetFileLocations(fileName string) map[string][]string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	fileMeta, exists := n.Files[fileName]

	if !exists {
		return nil
	}

	chunkLocations := make(map[string][]string)

	for _, chunkID := range fileMeta.ChunkIDs {
		chunkMeta, exists := n.Chunks[chunkID]
		if exists {
			chunkLocations[chunkID] = chunkMeta.ConfirmedNodes
		}
	}
	return chunkLocations
}

func NewNameNode() *NameNode {

	nn := &NameNode{
		DataNodes: make(map[string]DataNode),
		Files:     make(map[string]FileMetaData),
		Chunks:    make(map[string]ChunkMetaData),
	}

	go nn.monitorHealth()

	return nn
}
