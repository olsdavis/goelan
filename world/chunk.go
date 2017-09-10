package world

import (
	"github.com/olsdavis/goelan/world/val"
	"github.com/olsdavis/goelan/protocol"
)

type Chunk struct {
	Blocks     [][][]*Block
	StartPoint Location3i
}

func (chunk *Chunk) WriteChunk(packet protocol.Response) {
}

func NewChunk(startPoint Location3i) *Chunk {
	return &Chunk{
		Blocks:     make([][][]*Block, val.ChunkSize),
		StartPoint: startPoint,
	}
}
