package clusterer

import (
	"errors"
	"github.com/arjunsk/go-kmeans/containers"
	"sync"
	"sync/atomic"
)

type Clusterer interface {
	Cluster() (containers.Clusters, error)
}

func validateArgs(vectors []containers.Vector, clusterCnt int) error {
	if len(vectors) == 0 {
		return errors.New("kmeans: The data set must not be empty")
	}

	if clusterCnt > len(vectors) {
		return errors.New("kmeans: The count of the data set must at least equal k")
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
