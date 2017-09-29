package world

import (
	"github.com/olsdavis/goelan/world/val"
)

type Chunk struct {
	Blocks     [][][]*Block
	StartPoint Location3i
}

func NewChunk(startPoint Location3i) *Chunk {
	return &Chunk{
		Blocks:     make([][][]*Block, val.ChunkSize),
		StartPoint: startPoint,
	}
}
