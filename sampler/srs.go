package sampler

import (
	"math/rand"
	"time"
)

// SrsSampling ie Simple Random Sampling, sample input with equal probability.
func SrsSampling[T any](input []T, percent float64) []T {
	if percent >= 100.0 {
		return input
	}
	if percent <= 0.0 {
		return nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	n := int(float64(len(input)) * percent / 100.0) // calculate how many items based on percentage
	subset := make([]T, 0, n)
	selected := make(map[int]bool)

	for i := 0; i < n; i++ {
		idx := r.Intn(len(input))
		for selected[idx] {
			idx = r.Intn(len(input))
		}
		selected[idx] = true
		subset = append(subset, input[idx])
	}

	return subset
}
