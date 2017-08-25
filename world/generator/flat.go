package generator

import "github.com/olsdavis/goelan/world"

type FlatGenerator struct {}

func (generator FlatGenerator) GenerateChunk(x, y, z int, world world.World) *world.Chunk {
	ret := &world.Chunk{}
	return ret
}
