package generator

import (
	"github.com/olsdavis/goelan/world"
	"github.com/olsdavis/goelan/material"
	"github.com/olsdavis/goelan/world/val"
)

type FlatGenerator struct{}

func (generator FlatGenerator) GenerateChunkColumn(x, z int, w *world.World) *world.Chunk {
	ret := &world.Chunk{}
	for y := 0; y < 4; y++ {
		var mat material.Material
		switch y {
		case 0:
			mat = material.Bedrock
		case 1:
			mat = material.Dirt
		case 2:
			fallthrough
		case 3:
			mat = material.Grass
		}
		for x1 := 0; x1 < val.ChunkSize; x1++ {
			for z1 := 0; z1 < val.ChunkSize; z1++ {
				ret.Blocks[x1][y][z1] = world.NewBlock(world.NewLocation3i(int32(x+x1), int32(y), int32(z+z1), w), mat, 0)
			}
		}
	}
	return ret
}
