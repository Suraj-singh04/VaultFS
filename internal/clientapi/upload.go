package clientapi

import (
	"context"

	pbNameNode "github.com/Suraj-singh04/vaultfs/proto/namenode"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type chunk struct {
	id   string
	data []byte
}

func (c *Coordinator) UploadFile(filePath string, data []byte) error {
	fileSize := len(data)
	chunkSize := 64 * 1024 * 1024
	numChunks := (fileSize + chunkSize - 1) / chunkSize

	chunks := []chunk{}

	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > fileSize {
			end = fileSize
		}
		chunkID := uuid.New().String()
		chunks = append(chunks, chunk{id: chunkID, data: data[start:end]})
	}

	chunkIDs := []string{}
	for _, ch := range chunks {
		chunkIDs = append(chunkIDs, ch.id)
	}

	ctx := context.Background()

	resp, err := c.nameNodeClient.AllocateChunks(ctx, &pbNameNode.AllocateChunksRequest{
		FileName: filePath,
		ChunkIds: chunkIDs,
	})

	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, ch := range chunks {
		ch := ch
		nodeList := resp.ChunkLocations[ch.id]

		for _, nodeID := range nodeList.NodeIds {
			nodeID := nodeID
			g.Go(func() error {
				return c.sendChunkToNode(ctx, nodeID, nodeID, ch.id, ch.data)
			})
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}
	return nil

}
