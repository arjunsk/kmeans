package clusterer

import (
	"github.com/arjunsk/kmeans/containers"
	"github.com/arjunsk/kmeans/initializer"
	"math/rand"
)

type Lloyd struct {
	deltaThreshold     float64
	iterationThreshold int

	distFn      containers.DistanceFunction
	initializer initializer.Initializer

	// local state
	vectors    [][]float64
	clusterCnt int
}

var _ Clusterer = new(Lloyd)

func NewKmeans(vectors [][]float64, clusterCnt int) (Clusterer, error) {
	// Iteration count: 500 Ref: https://github.com/pgvector/pgvector/blob/8d7abb659070259e78f5c0974dde26c9e1cda8d3/src/ivfkmeans.c#L261
	// delta threshold: .01 Ref:https://github.com/muesli/kmeans/blob/06e72b51dbf15ea9e20146451e2c523389633707/kmeans.go#L44
	deltaThreshold := 0.01
	iterationThreshold := 500

	err := validateArgs(vectors, clusterCnt, deltaThreshold, iterationThreshold)
	if err != nil {
		return nil, err
	}

	m, err := newKmeansWithOptions(
		deltaThreshold,
		iterationThreshold,
		containers.EuclideanDistance,
		initializer.NewRandomInitializer())
	if err != nil {
		return nil, err
	}

	m.vectors = vectors
	m.clusterCnt = clusterCnt

	return m, nil
}

func newKmeansWithOptions(
	deltaThreshold float64,
	iterationThreshold int,
	distFn containers.DistanceFunction,
	init initializer.Initializer) (Lloyd, error) {

	return Lloyd{
		deltaThreshold:     deltaThreshold,
		iterationThreshold: iterationThreshold,
		distFn:             distFn,
		initializer:        init,
	}, nil
}

func (ll Lloyd) Cluster() (containers.Clusters, error) {

	clusters, err := ll.initializer.InitCentroids(ll.vectors, ll.clusterCnt)
	if err != nil {
		return nil, err
	}

	err = ll.kmeans(clusters)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

// kmeans Complexity := O(n*k*e*d); n = number of vectors, k = number of clusters, e = number of iterations, d = number of dimensions
func (ll Lloyd) kmeans(clusters containers.Clusters) (err error) {

	assignments := make([]int, len(ll.vectors))
	movement := 1

	for i := 0; ; i++ {
		//1. Reset the state
		movement = 0
		clusters.Reset()

		// 2. Assign vectors to the nearest cluster
		movement, err = ll.assignData(ll.vectors, clusters, assignments, movement)
		if err != nil {
			return err
		}

		// 3.b Update the cluster centroids for Empty clusters
		for clusterId := 0; clusterId < len(clusters); clusterId++ {
			if len(clusters[clusterId].GetMembers()) == 0 {
				//vecIdx represents an index of a vector from a "cluster with more than one member"
				var vecIdx int
				for {
					vecIdx = rand.Intn(len(ll.vectors))
					if len(clusters[assignments[vecIdx]].GetMembers()) > 1 {
						break
					}
				}
				clusters[clusterId].AddMember(ll.vectors[vecIdx])
				assignments[vecIdx] = clusterId
				movement = len(ll.vectors)
			}
		}

		// 4. Recenter the clusters
		if movement > 0 {
			err = clusters.Recenter()
			if err != nil {
				return err
			}
		}

		if ll.isConverged(i, movement) {
			break
		}
	}

	return nil
}

func (ll Lloyd) assignData(vectors [][]float64, clusters containers.Clusters, clusterIds []int, movement int) (int, error) {
	// 2. Assign each vector to the nearest cluster
	for vecIdx, vec := range vectors {
		clusterId, _, err := clusters.Nearest(vec, ll.distFn)
		if err != nil {
			return 0, err
		}
		clusters[clusterId].AddMember(vec)

		// 3.a Update the cluster id of the vector
		if clusterIds[vecIdx] != clusterId {
			clusterIds[vecIdx] = clusterId
			movement++
		}
	}
	return movement, nil
}

func (ll Lloyd) isConverged(i int, movement int) bool {
	vectorCnt := float64(len(ll.vectors))
	if i == ll.iterationThreshold || movement < int(vectorCnt*ll.deltaThreshold) || movement == 0 {
		return true
	}
	return false
}
