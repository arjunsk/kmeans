package clusterer

import (
	"github.com/arjunsk/go-kmeans/containers"
	"github.com/arjunsk/go-kmeans/initializer"
	"math"
	"sync"
)

// KmeansElkan Ref Paper: https://cdn.aaai.org/ICML/2003/ICML03-022.pdf
type KmeansElkan struct {
	deltaThreshold     float64
	iterationThreshold int

	distFn      containers.DistanceFunction
	initializer initializer.Initializer

	assignments []int
	lowerBounds [][]float64
	upperBounds []float64
	r           []bool

	// local state
	vectors    [][]float64
	clusterCnt int
}

var _ Clusterer = new(KmeansElkan)

func NewKmeansElkan(vectors [][]float64, clusterCnt int) (Clusterer, error) {
	deltaThreshold := 0.01
	iterationThreshold := 500

	err := validateArgs(vectors, clusterCnt, deltaThreshold, iterationThreshold)
	if err != nil {
		return nil, err
	}

	el, err := newKmeansElkanWithOptions(
		deltaThreshold,
		iterationThreshold,
		containers.EuclideanDistance,
		initializer.NewKmeansInitializer())
	if err != nil {
		return nil, err
	}

	el.vectors = vectors
	el.clusterCnt = clusterCnt

	n := len(vectors)
	el.r = make([]bool, n)
	el.assignments = make([]int, n)
	el.upperBounds = make([]float64, n)
	el.lowerBounds = make([][]float64, n)
	for i := range el.lowerBounds {
		el.lowerBounds[i] = make([]float64, clusterCnt)
	}

	return &el, nil
}

func newKmeansElkanWithOptions(
	deltaThreshold float64,
	iterationThreshold int,
	distFn containers.DistanceFunction,
	init initializer.Initializer) (KmeansElkan, error) {

	return KmeansElkan{
		deltaThreshold:     deltaThreshold,
		iterationThreshold: iterationThreshold,
		distFn:             distFn,
		initializer:        init,
	}, nil
}

func (el *KmeansElkan) Cluster() (containers.Clusters, error) {

	clusters, err := el.initializer.InitCentroids(el.vectors, el.clusterCnt)
	if err != nil {
		return nil, err
	}

	err = el.kmeansElkan(clusters)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

// kmeansElkan Complexity := closer to O(n); n = number of vectors
func (el *KmeansElkan) kmeansElkan(clusters containers.Clusters) (err error) {
	for i := 0; ; i++ {
		movement := 0
		clusters.Reset()

		centroidSelfDistances := el.calculateCentroidDistances(clusters, el.clusterCnt)
		sc := el.computeSc(centroidSelfDistances, el.clusterCnt)

		// step 3
		movement, err = el.assignData(centroidSelfDistances, sc, clusters, el.vectors, i)
		if err != nil {
			return err
		}

		// step 4 and 5
		moveDistances, err := clusters.RecenterWithDeltaDistance(el.distFn)
		if err != nil {
			return err
		}

		// step 5 and 6
		el.updateBounds(moveDistances, el.vectors)

		if el.isConverged(i, movement) {
			break
		}
	}

	return nil
}

func (el *KmeansElkan) calculateCentroidDistances(clusters containers.Clusters, k int) [][]float64 {
	centroidDistances := make([][]float64, k)
	for i := 0; i < k; i++ {
		centroidDistances[i] = make([]float64, k)
	}

	var wg sync.WaitGroup
	for i := 0; i < k-1; i++ {
		for j := i + 1; j < k; j++ {
			wg.Add(1)
			//TODO: parallelize this
			go (func(i, j int) {
				defer wg.Done()
				centroidDistances[i][j], _ = el.distFn(clusters[i].Center, clusters[j].Center)
				centroidDistances[j][i] = centroidDistances[i][j]
			})(i, j)
		}
	}
	wg.Wait()

	return centroidDistances
}

// s(c)	= 0.5 * min{d(c, c') | c' != c}
func (el *KmeansElkan) computeSc(centroidDistances [][]float64, k int) []float64 {
	sc := make([]float64, k)
	for i := 0; i < k; i++ {
		scMin := math.MaxFloat64
		for j := 0; j < k; j++ {
			if i == j {
				continue
			}
			scMin = math.Min(centroidDistances[i][j], scMin)
		}
		sc[i] = 0.5 * scMin
	}
	return sc
}

func (el *KmeansElkan) assignData(centroidDistances [][]float64,
	sc []float64,
	clusters containers.Clusters,
	vectors [][]float64,
	iterationCount int) (int, error) {

	moves := 0
	k := len(centroidDistances)

	for x := range vectors {
		//TODO: parallelize this

		// c(x)
		meanIndex := el.assignments[x]

		if el.upperBounds[x] <= sc[meanIndex] {
			continue
		}

		for c := 0; c < k; c++ {

			if c != meanIndex &&
				el.upperBounds[x] > el.lowerBounds[x][c] &&
				el.upperBounds[x] > centroidDistances[meanIndex][c]*0.5 {

				//step3a BoundsUpdate
				if el.r[x] {
					distance, err := el.distFn(vectors[x], clusters[meanIndex].Center)
					if err != nil {
						return 0, err
					}
					el.upperBounds[x] = distance
					el.lowerBounds[x][meanIndex] = distance
					el.r[x] = false
				}

				//step3b Update
				if el.upperBounds[x] > el.lowerBounds[x][c] ||
					el.upperBounds[x] > centroidDistances[meanIndex][c]*0.5 {
					newDistance, _ := el.distFn(vectors[x], clusters[c].Center)
					el.lowerBounds[x][c] = newDistance
					if newDistance < el.upperBounds[x] {
						meanIndex = c
						el.upperBounds[x] = newDistance
					}
				}

			}

		}
		if meanIndex != el.assignments[x] {
			el.assignments[x] = meanIndex
			moves++
		} else if iterationCount == 0 {
			moves++
		}

		clusters[meanIndex].Add(vectors[x])
	}
	return moves, nil
}

func (el *KmeansElkan) updateBounds(moveDistances []float64, data [][]float64) {
	k := len(moveDistances)

	for x := range data {
		for c := 0; c < k; c++ {
			el.lowerBounds[x][c] = math.Max(el.lowerBounds[x][c]-moveDistances[c], 0)
		}
		el.upperBounds[x] += moveDistances[el.assignments[x]]
		el.r[x] = true
	}
}

func (el *KmeansElkan) isConverged(i int, movement int) bool {
	vectorCnt := float64(len(el.vectors))
	if i == el.iterationThreshold || movement < int(vectorCnt*el.deltaThreshold) || movement == 0 {
		return true
	}
	return false
}
