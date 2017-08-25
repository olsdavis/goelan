package generator

import "github.com/olsdavis/goelan/world"

type WorldGenerator interface {
	GenerateChunk(x, y, z int, world *world.World) *world.Chunk
}
