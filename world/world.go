package world

type World struct {
	Name   string
	Chunks [][][]*Chunk
}

func NewWorld(name string) *World {
	return &World{
		Name: name,
		Chunks: make([][][]*Chunk, 16),
	}
}
