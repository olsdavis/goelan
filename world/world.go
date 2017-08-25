package world

const (
	ChunkSize = 16 // 16^3
)

type Chunk struct {
	blocks [][][]*Block
}

func NewChunk() *Chunk {
	return &Chunk{
		make([][][]*Block, ChunkSize),
	}
}

type World struct {
	Name   string
	Chunks [][][]Chunk
}

func NewWorld(name string) *World {
	return &World{
		Name: name,
	}
}
