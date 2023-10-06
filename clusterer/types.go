package clusterer

import (
	"errors"
	"github.com/arjunsk/kmeans/containers"
	"sync"
	"sync/atomic"
)

type Clusterer interface {
	// Cluster will run the clustering algorithm.
	// NOTE: We will not support adding clusterCnt as argument ie Cluster(clusterCnt).
	// This would require a factory pattern to create states
	// (assignments, lowerBounds, upperBounds, r) etc. for each call.
	Cluster() (containers.Clusters, error)
}

func validateArgs(vectors [][]float64, clusterCnt int, deltaThreshold float64, iterationThreshold int) error {
	if len(vectors) == 0 {
		return errors.New("kmeans: The data set must not be empty")
	}

	if clusterCnt > len(vectors) {
		return errors.New("kmeans: The count of the data set must at least equal k")
	}

	if deltaThreshold <= 0.0 || deltaThreshold >= 1.0 {
		return errors.New("kmeans: threshold is out of bounds (must be >0.0 and <1.0, in percent)")
	}

	if iterationThreshold <= 0 {
		return errors.New("kmeans: iterationThreshold must be > 0")
	}

	vecDim := len(vectors[0])
	var dimMismatch atomic.Bool
	var wg sync.WaitGroup
	for i := 1; i < len(vectors); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if len(vectors[i]) != vecDim {
				dimMismatch.Store(true)
			}
		}(i)
	}
	wg.Wait()
	if dimMismatch.Load() {
		return errors.New("kmeans: The data set must contain vectors of the same dimension")
	}

	return nil
}
