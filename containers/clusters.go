package containers

import (
	"fmt"
	"sync"
)

// Clusters is a collection of Cluster.
// None of the methods are pointer receivers, since we don't want to mutate the Clusters.
// We mutate the Cluster instead of Clusters. Cluster has a pointer receiver.
type Clusters []Cluster

// Nearest returns the index, distance of the cluster nearest to point
func (c Clusters) Nearest(point Vector, distFn DistanceFunction) (minClusterIdx int, minDistance float64, err error) {
	if distFn == nil {
		panic(fmt.Errorf("distance function is nil"))
	}

	minClusterIdx = 0

	var currDistance = 0.0
	minDistance, err = distFn(point, c[0].GetCenter())
	if err != nil {
		return 0, 0, err
	}

	for i := 1; i < len(c); i++ {
		currDistance, err = distFn(point, c[i].GetCenter())
		if err != nil {
			return 0, 0, err
		}
		if currDistance < minDistance {
			minDistance = currDistance
			minClusterIdx = i
		}
	}

	return minClusterIdx, minDistance, nil
}

func (c Clusters) Recenter() error {
	clusterCnt := len(c)
	var wg sync.WaitGroup
	for i := 0; i < clusterCnt; i++ {
		wg.Add(1)
		go (func(i int) {
			defer wg.Done()
			c[i].Recenter()
		})(i)
	}
	wg.Wait()
	return nil
}

func (c Clusters) RecenterWithDeltaDistance(distFn DistanceFunction) (moveDistances []float64, err error) {
	if distFn == nil {
		return nil, fmt.Errorf("distance function is nil")
	}

	clusterCnt := len(c)
	moveDistances = make([]float64, clusterCnt)

	var wg sync.WaitGroup
	for i := 0; i < clusterCnt; i++ {
		wg.Add(1)
		//TODO: parallelize this
		go (func(i int) {
			defer wg.Done()
			moveDistances[i], _ = c[i].RecenterWithMovedDistance(distFn)
			// NOTE: ignoring error here since RecenterReturningMovedDistance() will return an error
			// only if the distance function returns an error. We are not returning the unhandled error from the
			// distance function.
		})(i)

	}
	wg.Wait()
	return moveDistances, nil
}

func (c Clusters) Reset() {
	for i := 0; i < len(c); i++ {
		c[i].Reset()
	}
}

func (c Clusters) String() string {
	var s = ""
	for i := 0; i < len(c); i++ {
		s += fmt.Sprintf("%d: %s\n", i, c[i].String())
	}
	return s
}

func (c Clusters) SSE() float64 {
	var sse = 0.0
	for i := 0; i < len(c); i++ {
		sse += c[i].SSE()
	}
	return sse
}
