package util

import (
	"testing"
)

func TestClosestMultiple(t *testing.T) {
	if c := ClosestMultiple(16, 16); c != 16 {
		t.Error("Closest multiple of 16 to 16 should be 16. Currently returns", c)
	}

	if c := ClosestMultiple(17, 16); c != 16 {
		t.Error("Closest multiple of 17 to 16 should be 16. Currently returns", c)
	}

	// the function closest multiple rounds down to the closest multiple
	if c := ClosestMultiple(31, 16); c != 16 {
		t.Error("Closest multipe of 31 to 16 should be 16. Currently returns", c)
	}

	if c := ClosestMultiple(-16, 16); c != -16 {
		t.Error("Closest multiple of -16 to 16 should be -16. Currently returns", c)
	}
}

func TestSquareFloat32(t *testing.T) {
	if c := SquareFloat32(16); c != 256 {
		t.Error("Square of 16 is 256. Currently returns", c)
	}
}

func TestSquareInt(t *testing.T) {
	if c := SquareInt(16); c != 256 {
		t.Error("Square of 16 if 256. Currently returns", c)
	}
}

func TestMin(t *testing.T) {
	if c := Min(3, 2); c != 2 {
		t.Error("Min of 3 and 2 should return 2. Currently returns", c)
	}

	if c := Min(2, 3); c != 2 {
		t.Error("Min of 2 and 3 should return 2. Currently returns", c)
	}

	if c := Min(-1, 10); c != -1 {
		t.Error("Min of -1 and 10 should return -1. Currently returns", c)
	}
}
