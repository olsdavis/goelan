package world

import (
	"../util"
	"math"
)

// Represents rotations (yaw & pitch)
type Orientation struct {
	Yaw   float32
	Pitch float32
}

// Used for blocks
type SimpleLocation struct {
	X     float32
	Y     float32
	Z     float32
	World *World
}

// Used for entities
type Location struct {
	Orientation
	SimpleLocation
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

// Creates a SimpleLocation with the given (x,y,z) coordinates and the world.
// In contrast with the Location type, the SimpleLocation has not got an orientation.
// Panics if the world is nil.
func NewSimpleLocation(x, y, z float32, world *World) *SimpleLocation {
	if world == nil {
		panic("world cannot be nil")
	}

	return &SimpleLocation{x, y, z, world}
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
