package util

// Returns the smallest integer of x and y. If x equals y, returns x.
// This function is the same as math.Min,
// the only difference is that it works with integers,
// so it won't create any bugs with floating points.
//
// See: https://mrekucci.blogspot.fr/2015/07/dont-abuse-mathmax-mathmin.html
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Returns the square of the given integer.
func SquareInt(x int) int {
	return x * x
}

// Returns the square of the given float (32)
func SquareFloat32(x float32) float32 {
	return x * x
}
