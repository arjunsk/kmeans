package initializer

import (
	"errors"
	"github.com/arjunsk/kmeans/containers"
)

// Initializer is an interface and not a function type because initializer like kmeans++
// requires extra arguments.
type Initializer interface {
	InitCentroids(vectors [][]float64, clusterCnt int) (containers.Clusters, error)
}

func validateArgs(vectors [][]float64, clusterCnt int) error {
	if clusterCnt <= 0 {
		return errors.New("KMeans: k cannot be less than or equal to zero")
	}
	if len(vectors) == 0 {
		return errors.New("KMeans: data cannot be empty")
	}
	if vectors == nil || len(vectors[0]) == 0 {
		return errors.New("KMeans: data cannot be nil")
	}
	return nil
}
