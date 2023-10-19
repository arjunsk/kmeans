package initializer

import (
	"errors"
	"github.com/arjunsk/kmeans/containers"
	"math/rand"
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
// Ref (animation):https://www.youtube.com/watch?v=efKGmOH4Y_A
// Complexity: O(k*n*k); n = number of vectors, k = number of clusters
// Reason: For k-1 times, compute the distance of each vector to its nearest center O(nk)
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
	clusters[0] = containers.NewCluster(vectors[randIdx])

	// O(k-1)
	for i := 1; i < clusterCnt; i++ {
		// NOTE: Since Nearest function is called on clusters-1, parallel handling
		// can cause bugs, since all the clusters are not initialized.
		distances := make([]float64, len(vectors))
		sum := 0.0
		minDistance := 0.0
		// 2. for each data point, compute the distance to the existing centers
		// O(n)
		for vecIdx, vec := range vectors {

			// O(k)
			_, minDistance, err = clusters[:i].Nearest(vec, kpp.DistFn)
			if err != nil {
				return nil, err
			}

			distances[vecIdx] = minDistance * minDistance // D(x)^2
			sum += distances[vecIdx]
		}

		// 3. choose the next random center, using a weighted probability distribution
		// where it is chosen with probability proportional to D(x)^2
		// Ref: https://en.wikipedia.org/wiki/K-means%2B%2B#Improved_initialization_algorithm
		// Ref: https://stats.stackexchange.com/a/272133/397621
		target := rand.Float64() * sum
		idx := 0
		for currSum := distances[0]; currSum < target; currSum += distances[idx] {
			idx++
		}

		// Select a cluster center based on a probability distribution where vectors
		//	with larger distances have a higher chance of being chosen as the center.
		clusters[i] = containers.NewCluster(vectors[idx])

	}
	return clusters, nil
}
