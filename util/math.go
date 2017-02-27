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

// Returns the square of the given float (32).
func SquareFloat32(x float32) float32 {
	return x * x
}

// ClosestMultiple returns the closest multiple of the given number n.
// Example:
// ClosestMultiple(18, 16) returns 16
func ClosestMultiple(n, multiple int) int {
	if multiple == 0 {
		return 0
	}

	if n >= 0 {
		return n / multiple  * multiple
	}

	return ((n - multiple + 1) / multiple) * multiple
}

func ToRange(limitMax, limitMin, baseMax, baseMin, v float32) float32 {
	return ((limitMax - limitMin) * (v - baseMin) / (baseMax - baseMin)) + limitMin
}
