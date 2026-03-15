package utils

import "math"

// RoundToTwo rounds a float64 to 2 decimal places.
func RoundToTwo(f float64) float64 {
	return math.Round(f*100) / 100
}
