package world

import (
	"../util"
	"sync"
)

const (
	ChunkSize = 16
)

type Chunk struct {
	X, Y, Z int // in chunk coordinates
	blocks  [][][]Block
}

type ChunkManager struct {
	world        World
	chunkLock    sync.Mutex
	loadedChunks map[byte]map[byte]map[byte]Chunk // x to y to z to Chunk
}

// Returns the chunk at the given location.
// Returns nil if the chunk has not been loaded yet.
func (manager ChunkManager) GetChunkAtLocation(location SimpleLocation) *Chunk {
	return chunkCoords(location)
}

// Returns the chunk at the given chunk coordinates.
func (manager *ChunkManager) GetChunkAtChunkCoordinates(x, y, z int) *Chunk {
	defer manager.chunkLock.Unlock()
	manager.chunkLock.Lock()
	arr := manager.loadedChunks[x]
	if arr == nil {
		return nil
	}
	arr = arr[y]
	if arr == nil {
		return nil
	}
	return arr[z]
}

// Loads the chunk at the given location.
func (manager *ChunkManager) LoadChunk(location SimpleLocation) *Chunk {
	x, y, z := chunkCoords(location)
	chunk := manager.GetChunkAtChunkCoordinates(x, y, z)
	if chunk == nil {
		chunk = &Chunk{
			blocks: make([][][]byte, 0, ChunkSize),
			X: x,
			Y: y,
			Z: z,
		}
		manager.GenerateChunk(chunk)
		manager.chunkLock.Lock()
		manager.loadedChunks[x][y][z] = chunk
		manager.chunkLock.Unlock()
		return chunk
	} else {
		return chunk
	}
}

// Generates the given chunk.
func (manager *ChunkManager) GenerateChunk(chunk *Chunk) {
	if chunk == nil {
		panic("chunk cannot be nil")
	}
}

func chunkCoords(location SimpleLocation) (int, int, int) {
	x := util.ClosestMultiple(location.X, 8) / 8
	y := util.ClosestMultiple(location.Y, 8) / 8
	z := util.ClosestMultiple(location.Z, 8) / 8

	if location.X < 0 {
		x = -x
	}
	if location.Y < 0 {
		y = -y
	}
	if location.Z < 0 {
		z = -z
	}
	return x, y, z
}
