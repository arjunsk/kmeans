package initializer

import (
	"errors"
	"github.com/arjunsk/go-kmeans/containers"
)

type Initializer interface {
	InitCentroids(vectors [][]float64, clusterCnt int) (containers.Clusters, error)
}

func StdInputChecks(vectors [][]float64, clusterCnt int, inputCnt int) error {
	if vectors == nil || len(vectors[0]) == 0 {
		return errors.New("KMeans: data cannot be nil")
	}
	if clusterCnt <= 0 {
		return errors.New("KMeans: k cannot be less than or equal to zero")
	}
	if inputCnt == 0 {
		return errors.New("KMeans: data cannot be empty")
	}
	return nil
}
