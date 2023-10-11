package sampler

// This pkg implements sub-sampling algorithms to pick a subset of items from the input

type Sampling[T any] func(input []T, percent float64) []T
