package world

import (
	"sync"
	opensimplex "github.com/ojrac/opensimplex-go"
	"math"
	"fmt"
)

const (
	ChunkSize = 16
)

type Chunk struct {
	X, Y, Z int // in chunk coordinates
	blocks  map[int8]map[int8]map[int8]*Block
	sync.Mutex
}

// Sets the given block at the given location.
func (chunk *Chunk) PutBlock(x, y, z int8, block *Block) {
	if x > ChunkSize || x < 0 || y > ChunkSize || y < 0 || z > ChunkSize || z < 0 {
		panic("x, y or z coordinate out of bounds")
	}

	chunk.Lock()
	chunk.doSetBlock(x, y, z, block)
	chunk.Unlock()
}

func (chunk *Chunk) doSetBlock(x, y, z int8, block *Block) {
	yMap := chunk.blocks[x]
	if yMap == nil {
		yMap = make(map[int8]map[int8]*Block)
		chunk.blocks[x] = yMap
		zMap := make(map[int8]*Block)
		yMap[y] = zMap
	}
	zMap := yMap[y]
	if zMap == nil {
		zMap = make(map[int8]*Block)
		yMap[y] = zMap
	}
	zMap[z] = block
}

type ChunkManager struct {
	world        *World
	simplex      *opensimplex.Noise
	chunkLock    sync.Mutex
	loadedChunks map[int]map[int]map[int]*Chunk // x to y to z to Chunk
}

// Creates a new manager for the given world.
func NewManager(seed int64) ChunkManager {
	return ChunkManager{
		simplex: opensimplex.NewWithSeed(seed),
		chunkLock: sync.Mutex{},
		loadedChunks: make(map[int]map[int]map[int]*Chunk),
	}
}

// Returns the chunk at the given location.
// Returns nil if the chunk has not been loaded yet.
func (manager ChunkManager) GetChunkAtLocation(location SimpleLocation) *Chunk {
	return nil
}

// Returns the chunk at the given chunk coordinates.
func (manager ChunkManager) GetChunkAtChunkCoordinates(x, y, z int) *Chunk {
	defer manager.chunkLock.Unlock()
	manager.chunkLock.Lock()
	arr := manager.loadedChunks[x]
	if arr == nil {
		return nil
	}
	second := arr[y]
	if second == nil {
		return nil
	}
	return second[z]
}

// Loads the chunk at the given location.
func (manager ChunkManager) LoadChunk(location SimpleLocation) *Chunk {
	x, y, z := ChunkCoords(location)
	chunk := manager.GetChunkAtChunkCoordinates(x, y, z)
	if chunk == nil {
		chunk = &Chunk{
			blocks: make(map[int8]map[int8]map[int8]*Block),
			X: x,
			Y: y,
			Z: z,
		}
		manager.GenerateChunk(chunk)
		manager.putChunk(x, y, z, chunk)
		return chunk
	} else {
		return chunk
	}
}

func (manager ChunkManager) putChunk(x, y, z int, chunk *Chunk) {
	manager.chunkLock.Lock()
	yMap := manager.loadedChunks[x]
	if yMap == nil {
		yMap = make(map[int]map[int]*Chunk)
		manager.loadedChunks[x] = yMap
		zMap := make(map[int]*Chunk)
		yMap[y] = zMap
	}
	zMap := yMap[z]
	if zMap == nil {
		zMap = make(map[int]*Chunk)
		yMap[z] = zMap
	}
	zMap[z] = chunk
	manager.chunkLock.Unlock()
}

// Generates the given chunk.
func (manager ChunkManager) GenerateChunk(chunk *Chunk) {
	if chunk == nil {
		panic("chunk cannot be nil")
	}

	for x := 0; x < ChunkSize; x++ {
		for z := 0; z < ChunkSize; z++ {
			y := manager.smoothNoiseAt(x, z)
			fmt.Println(y)
		}
	}
}

func (manager *ChunkManager) smoothNoiseAt(x, z int) int {
	return 0
}

func (manager *ChunkManager) noiseAt(x, z int) int {
	return int(math.Abs(manager.simplex.Eval2(float64(x), float64(z)) * 132))
}

// Converts the given location to chunk coordinates.
func ChunkCoords(location SimpleLocation) (int, int, int) {
	x := int(math.Floor(float64(location.X / ChunkSize)))
	y := int(math.Floor(float64(location.Y / ChunkSize)))
	z := int(math.Floor(float64(location.Z / ChunkSize)))
	return x, y, z
}
