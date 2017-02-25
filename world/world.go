package world

type WorldType int

const (
	Overworld WorldType = iota
	Nether
	TheEnd
)

type World struct {
	ChunkManager   ChunkManager
	Seed           int64
	MaxBuildHeight int64
	Type           WorldType
}

func (w *World) GetBlockAt(location SimpleLocation) *Block {
	return nil
}
