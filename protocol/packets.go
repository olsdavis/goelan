// This file contains the structs of some packets that are
// used in different parts of the project.

package protocol

type (
	PositionAndLookPacket struct {
		X          float64
		Y          float64
		Z          float64
		Yaw        float32
		Pitch      float32
		Flags      int8
		TeleportID int32
	}

	PlayerAbilitiesPacket struct {
		Flags       int8
		FlyingSpeed float32
		FovModifier float32
	}
)
