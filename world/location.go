package world

import (
	"github.com/olsdavis/goelan/util"
	"math"
)

// Represents rotations (yaw & pitch)
type Orientation struct {
	Yaw   float32
	Pitch float32
}

type Location3f struct {
	X     float32
	Y     float32
	Z     float32
	World *World
}

// Location3i is used for blocks.
type Location3i struct {
	X     int32
	Y     int32
	Z     int32
	World *World
}

// Used for entities
type Location struct {
	Orientation
	Location3f
}

// NewLocation3i creates a new Location3i and returns its pointer.
// Panics if the world is nil.
func NewLocation3i(x, y, z int32, world *World) *Location3i {
	if world == nil {
		panic("world cannot be nil")
	}

	return &Location3i{
		X: x,
		Y: y,
		Z: z,
		World: world,
	}
}

// Creates a new location with the given (x,y,z) coordinates, world
// and default 0 yaw and pitch orientation.
// Panics if the world is nil.
func NewLocation(x, y, z float32, world *World) *Location {
	return NewFullLocation(0, 0, x, y, z, world)
}

// Creates a new location with the given yaw, pitch and (x,y,z) coordinates.
// Panics if world is nil.
func NewFullLocation(yaw, pitch, x, y, z float32, world *World) *Location {
	return &Location{Orientation{yaw, pitch}, NewSimpleLocation(x, y, z, world)}
}

// Creates a Location3f with the given (x,y,z) coordinates and the world.
// In contrast with the Location type, the Location3f has not got an orientation.
// Panics if the world is nil.
func NewSimpleLocation(x, y, z float32, world *World) Location3f {
	if world == nil {
		panic("world cannot be nil")
	}

	return Location3f{x, y, z, world}
}

// Calculates the squared distance between the two locations.
// (Useful for avoiding the math.Sqrt() call, which can be quite expensive.)
// Returns 0 if the other location is nil.
func (l *Location) DistanceSquared(other *Location) float32 {
	if other == nil {
		return 0
	}

	return util.SquareFloat32(l.X-other.X) + util.SquareFloat32(l.Y-other.Y) + util.SquareFloat32(l.Z-other.Z)
}

// Calculates the distance between the two locations.
// Returns 0 if the other location is nil.
func (l *Location) Distance(other *Location) float32 {
	return float32(math.Sqrt(float64(l.DistanceSquared(other))))
}
