package shared

import "math"

const floatEpsilon = 0.0001

func FloatEqual(left float64, right float64) bool {
	return math.Abs(left-right) < floatEpsilon
}

func FloatGreater(left float64, right float64) bool {
	return left-right > floatEpsilon
}

func FloatLess(left float64, right float64) bool {
	return right-left > floatEpsilon
}

func FloatIsZero(value float64) bool {
	return FloatEqual(value, 0)
}
