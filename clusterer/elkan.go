package clusterer

import (
	"github.com/arjunsk/kmeans/containers"
	"github.com/arjunsk/kmeans/initializer"
	"golang.org/x/sync/errgroup"
	"math"
)

// KmeansElkan Ref Paper: https://cdn.aaai.org/ICML/2003/ICML03-022.pdf
// Slides: https://slideplayer.com/slide/9088301/
type KmeansElkan struct {
	deltaThreshold     float64
	iterationThreshold int

	distFn      containers.DistanceFunction
	initializer initializer.Initializer

	assignments []int       // maps vector index to cluster index
	lowerBounds [][]float64 // distances for vector and all clusters centroids
	upperBounds []float64   // distance between each point and its assigned cluster centroid ie d(x, c(x))

	// local state
	vectors    [][]float64 // input vectors
	clusterCnt int         // number of clusters ie k
}

var _ Clusterer = new(KmeansElkan)

func NewKmeansElkan(vectors [][]float64, clusterCnt int,
	deltaThreshold float64,
	iterationThreshold int,
	distFn containers.DistanceFunction,
	init initializer.Initializer) (Clusterer, error) {

	err := validateArgs(vectors, clusterCnt, deltaThreshold, iterationThreshold)
	if err != nil {
		return nil, err
	}

	n := len(vectors)

	el := KmeansElkan{
		deltaThreshold:     deltaThreshold,
		iterationThreshold: iterationThreshold,
		distFn:             distFn,
		initializer:        init,
		vectors:            vectors,
		clusterCnt:         clusterCnt,
		assignments:        make([]int, n),
		upperBounds:        make([]float64, n),
	}

	el.lowerBounds = make([][]float64, n)
	for i := range el.lowerBounds {
		el.lowerBounds[i] = make([]float64, clusterCnt)
	}

	return &el, nil
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

// kmeansElkan
// During each iteration of the algorithm, the lower bounds l(x, c) are updated for all points x and centers
// c. These updates take O(nk) time, so the complexity of the algorithm remains at least O(nke), even
// though the number of distance calculations is roughly O(n) only.
// Ref:https://www.cse.iitd.ac.in/~rjaiswal/2015/col870/Project/Nipun.pdf
// This variant needs O(n*k) additional memory to store bounds.
func (el *KmeansElkan) kmeansElkan(clusters containers.Clusters) error {
	for i := 0; ; i++ {
		movement := 0
		el.reset(clusters)
		clusters.Reset()

		// step 1.a
		// For all centers c and c', compute d(c, c').
		centroidSelfDistances := el.calculateCentroidDistances(clusters, el.clusterCnt)

		// step 1.b
		//  For all centers c, compute s(c)=0.5 min {d(c, c') | c'!= c}.
		sc := el.computeSc(centroidSelfDistances, el.clusterCnt)

		// step 2 and 3
		movement, err := el.assignData(centroidSelfDistances, sc, clusters, el.vectors, i)
		if err != nil {
			return err
		}

		// step 4
		// For each center c, let m(c) be the mean of the points assigned to c.
		// step 7
		// Replace each center c by m(c).
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

	//NOTE: We can parallelize this because [i][j] is computed on lower triangle.
	//[i][j] computed don't read any other [r][c] value.
	eg := new(errgroup.Group)
	for i := 0; i < k-1; i++ {
		for j := i + 1; j < k; j++ {
			func(r, c int) {
				eg.Go(func() error {
					var err error
					centroidDistances[r][c], err = el.distFn(clusters[r].Center(), clusters[c].Center())
					if err != nil {
						return err
					}
					centroidDistances[c][r] = centroidDistances[r][c]
					return nil
				})
			}(i, j)
		}
	}
	if err := eg.Wait(); err != nil {
		panic(err)
	}

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

		// c(x) in the paper
		meanIndex := el.assignments[x]

		// step 2.
		// Identify all points x such that u(x) <= s(c(x)).
		if el.upperBounds[x] <= sc[meanIndex] {
			continue
		}

		r := true //indicates that upper bound needs to be recalculated

		// step 3.
		// For all remaining points x and centers c such that
		// (i) c != c(x) and
		// (ii) u(x)>l(x, c) and
		// (iii) u(x)> 0.5 d(c(x), c):
		for c := 0; c < k; c++ {

			if c == meanIndex {
				continue // Pruned because this cluster is already the assignment.
			}

			if el.upperBounds[x] <= el.lowerBounds[x][c] {
				continue // Pruned by triangle inequality on lower bound.
			}

			if el.upperBounds[x] <= centroidDistances[meanIndex][c]*0.5 {
				continue // Pruned by triangle inequality on cluster distances.
			}

			//step 3.a - Bounds update
			// If r(x) then compute d(x, c(x)) and assign r(x)= false. Otherwise, d(x, c(x))=u(x).
			if r {
				distance, err := el.distFn(vectors[x], clusters[meanIndex].Center())
				if err != nil {
					return 0, err
				}
				el.upperBounds[x] = distance
				el.lowerBounds[x][meanIndex] = distance
				r = false
			}

			//step 3.b - Update
			// If d(x, c(x))>l(x, c) or d(x, c(x))> 0.5 d(c(x), c) then
			// Compute d(x, c)
			// If d(x, c)<d(x, c(x)) then assign c(x)=c.
			if el.upperBounds[x] > el.lowerBounds[x][c] ||
				el.upperBounds[x] > centroidDistances[meanIndex][c]*0.5 {
				newDistance, _ := el.distFn(vectors[x], clusters[c].Center())
				el.lowerBounds[x][c] = newDistance
				if newDistance < el.upperBounds[x] {
					meanIndex = c
					el.upperBounds[x] = newDistance
				}
			}

		}

		// Update Assigment/Membership & use this info to later update Mean Centroid
		if meanIndex != el.assignments[x] {
			el.assignments[x] = meanIndex
			moves++
		} else if iterationCount == 0 {
			moves++
		}
		clusters[meanIndex].AddMember(vectors[x])
	}
	return moves, nil
}

func (el *KmeansElkan) updateBounds(moveDistances []float64, data [][]float64) {
	k := len(moveDistances)

	for x := range data {
		for c := 0; c < k; c++ {
			//Step 5
			//For each point x and center c, assign
			// l(x, c)= max{ l(x, c)-d(c, m(c)), 0 }
			el.lowerBounds[x][c] = math.Max(el.lowerBounds[x][c]-moveDistances[c], 0)
		}
		// Step 6
		// For each point x, assign
		// u(x)=u(x)+d(m(c(x)), c(x))
		// r(x)= true
		el.upperBounds[x] += moveDistances[el.assignments[x]]
	}
}

func (el *KmeansElkan) isConverged(i int, movement int) bool {
	vectorCnt := float64(len(el.vectors))
	if i == el.iterationThreshold || movement < int(vectorCnt*el.deltaThreshold) || movement == 0 {
		return true
	}
	return false
}

func (el *KmeansElkan) reset(clusters containers.Clusters) {
	clusters.Reset()
	for i := range el.upperBounds {
		el.upperBounds[i] = math.MaxFloat64
	}
}
