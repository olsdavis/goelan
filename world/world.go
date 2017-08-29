package world

import "github.com/olsdavis/goelan/world/val"

type Chunk struct {
	Blocks [][][]*Block
}

func NewChunk() *Chunk {
	return &Chunk{
		Blocks: make([][][]*Block, val.ChunkSize),
	}
}

type World struct {
	Name   string
	Chunks [][][]Chunk
}

func NewWorld(name string) *World {
	return &World{
		Name: name,
		Chunks: make([][][]Chunk, 16),
	}
}
