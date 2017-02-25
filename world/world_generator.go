package world

import "../material"

type ChunkGenerator interface {
	Generate(*World, *Chunk)
}

type FlatWorldChunkGenerator struct {
	noise NoiseGenerator
}

// Generates the given chunk.
func (flat FlatWorldChunkGenerator) Generate(world *World, chunk *Chunk) {
	if chunk == nil {
		panic("asked to generate a nil chunk")
	}

	for x := 0; x < ChunkSize; x++ {
		for z := 0; z < ChunkSize; z++ {
			y := flat.noise.Eval(x, z)
			chunk.doSetBlock(int8(x), int8(y), int8(z), &Block{
				BlockState: 0,
				location: NewSimpleLocation(float32(x * chunk.X), float32(y * chunk.Y), float32(z * chunk.Z), world),
				material: material.Dirt,
			})
		}
	}
}

type NoiseGenerator interface {
	// Eval takes in the x and z coordinates and returns height of the block.
	Eval(int, int) int
}

type FlatWorldNoiseGenerator struct {}

// The flat world is composed of 4 layers of blocks.
// 1 => bedrock, 2 to 3 => dirt, 4 => grass
func (flat FlatWorldNoiseGenerator) Eval(x, y int) int {
	return 4
}
