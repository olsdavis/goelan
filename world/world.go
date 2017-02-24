package world

type World struct {
	ChunkManager ChunkManager
	Seed         int
}

func (w *World) GetBlock(x, y, z int) *Block {
	return nil
}
