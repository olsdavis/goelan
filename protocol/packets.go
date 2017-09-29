// This file contains the structs of some packets that are
// used in different parts of the project.

package protocol

import "github.com/olsdavis/goelan/world"

type (
	JsonString interface {}

	PositionAndLookPacket struct {
		world.Location
		Flags      int8
		TeleportID int32
	}

	PlayerAbilitiesPacket struct {
		Flags       int8
		FlyingSpeed float32
		FovModifier float32
	}
)
