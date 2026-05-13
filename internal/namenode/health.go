package namenode

import "time"

func (n *NameNode) monitorHealth() {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		n.mu.Lock()
		for id, node := range n.DataNodes {
			if time.Since(node.LastHeartbeat) > 15*time.Second {
				delete(n.DataNodes, id)
			}
		}
		n.mu.Unlock()
	}
}
