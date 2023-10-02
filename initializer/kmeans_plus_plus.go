package initializer

import (
	"errors"
	"github.com/arjunsk/go-kmeans/containers"
	"math/rand"
	"sync"
)

type KmeansPlusPlus struct {
	DistFn containers.DistanceFunction
}

func NewKmeansPlusPlusInitializer(distFn containers.DistanceFunction) Initializer {
	return &KmeansPlusPlus{
		DistFn: distFn,
	}
}

// InitCentroids initializes the centroids using kmeans++ algorithm
// Ref: https://www.youtube.com/watch?v=HatwtJSsj5Q
func (kpp *KmeansPlusPlus) InitCentroids(vectors [][]float64, clusterCnt int) (clusters containers.Clusters, err error) {
	err = validateArgs(vectors, clusterCnt)
	if err != nil {
		return nil, err
	}

	if kpp.DistFn == nil {
		return nil, errors.New("KMeans: distance function cannot be nil")
	}

	clusters = make([]containers.Cluster, clusterCnt)

	// 1. start with a random center
	randIdx := rand.Intn(len(vectors))
	clusters[0] = containers.Cluster{
		Center: vectors[randIdx],
	}

	var wg sync.WaitGroup
	for i := 1; i < clusterCnt; i++ {
		//TODO: see if we can parallelize this loop
		wg.Add(1)
		go (func(i int) {
			defer wg.Done()
			distances := make([]float64, len(vectors))
			sum := 0.0
			minDistance := 0.0
			// 2. for each data point, compute the distance to the existing centers
			for vecIdx, vec := range vectors {
				_, minDistance, _ = clusters[:i].Nearest(vec, kpp.DistFn)
				// NOTE: ignoring error here since Nearest() will return an error
				// only if the distance function returns an error. We are not returning the unhandled error from the
				// distance function.

				distances[vecIdx] = minDistance * minDistance // D(x)^2
				sum += distances[vecIdx]
			}

			// 3. choose the next random center, using a weighted probability distribution
			// where it is chosen with probability proportional to D(x)^2
			// Ref: https://en.wikipedia.org/wiki/K-means%2B%2B#Improved_initialization_algorithm
			// Ref: https://stats.stackexchange.com/a/272133/397621
			target := rand.Float64() * sum
			nextClusterCenterIdx := 0
			for sum = distances[0]; sum < target; sum += distances[nextClusterCenterIdx] {
				nextClusterCenterIdx++
			}

			// Select a cluster center based on a probability distribution where vectors
			//	with larger distances have a higher chance of being chosen as the center.
			clusters[i] = containers.Cluster{Center: vectors[nextClusterCenterIdx]}
		})(i)

	}
	wg.Wait()
	return clusters, nil
}
