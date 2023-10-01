package initializer

import (
	"errors"
	"go-kmeans/domain"
)

type Initializer interface {
	InitCentroids(vectors []domain.Vector, clusterCnt int) (domain.Clusters, error)
}

func StdInputChecks(vectors []domain.Vector, clusterCnt int, inputCnt int) error {
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
