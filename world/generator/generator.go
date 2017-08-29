package generator

import "github.com/olsdavis/goelan/world"

type WorldGenerator interface {
	GenerateChunkColumn(x, z int, world *world.World) *world.Chunk
}
