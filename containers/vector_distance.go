package containers

import (
	"fmt"
	"math"
)

// DistanceFunction is a function to find distance between 2 vectors
type DistanceFunction func(v1, v2 Vector) (float64, error)

// EuclideanDistance returns the Euclidean distance between two vectors
// Ref: https://mathworld.wolfram.com/L2-Norm.html
func EuclideanDistance(v1, v2 Vector) (float64, error) {
	if len(v1) != len(v2) {
		return 0, fmt.Errorf("vectors must have the same length")
	}

	distance := 0.0
	for c := range v1 {
		distance += math.Pow(v1[c]-v2[c], 2)
	}
	return math.Sqrt(distance), nil
}
