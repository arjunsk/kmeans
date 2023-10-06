package kmeans

import (
	"github.com/arjunsk/kmeans/clusterer"
	"github.com/arjunsk/kmeans/containers"
	"github.com/arjunsk/kmeans/initializer"
)

type ClustererType int

const (
	ELKAN ClustererType = iota
	KMEANS
	KMEANS_PLUS_PLUS
)

type Options struct {
	clusterType        ClustererType
	vectors            [][]float64
	clusterCnt         int
	distanceFn         containers.DistanceFunction
	initializer        initializer.Initializer
	deltaThreshold     float64
	iterationThreshold int
}

type Builder struct {
	options *Options
}

// NewClusterBuilder Create a new Options builder
func NewClusterBuilder(clustererType ClustererType, vectors [][]float64, clusterCnt int) *Builder {
	return &Builder{
		options: &Options{
			clusterType: clustererType,
			vectors:     vectors,
			clusterCnt:  clusterCnt,
		},
	}
}

func (b *Builder) DistanceFn(fn containers.DistanceFunction) *Builder {
	b.options.distanceFn = fn
	return b
}

func (b *Builder) Initializer(init initializer.Initializer) *Builder {
	b.options.initializer = init
	return b
}

func (b *Builder) DeltaThreshold(delta float64) *Builder {
	b.options.deltaThreshold = delta
	return b
}

func (b *Builder) IterationThreshold(iter int) *Builder {
	b.options.iterationThreshold = iter
	return b
}

// Build the Options instance
func (b *Builder) Build() (clusterer.Clusterer, error) {
	if b.options.initializer == nil && b.options.clusterType != KMEANS_PLUS_PLUS {
		b.options.initializer = initializer.NewRandomInitializer()
	}
	if b.options.distanceFn == nil {
		b.options.distanceFn = containers.EuclideanDistance
	}
	if b.options.deltaThreshold == 0 {
		b.options.deltaThreshold = 0.01
	}
	if b.options.iterationThreshold == 0 {
		b.options.iterationThreshold = 500
	}

	switch b.options.clusterType {
	case ELKAN:
		return clusterer.NewKmeansElkan(b.options.vectors, b.options.clusterCnt,
			b.options.deltaThreshold, b.options.iterationThreshold,
			b.options.distanceFn, b.options.initializer)
	case KMEANS:
		return clusterer.NewKmeans(b.options.vectors, b.options.clusterCnt,
			b.options.deltaThreshold, b.options.iterationThreshold,
			b.options.distanceFn, b.options.initializer)
	case KMEANS_PLUS_PLUS:
		if b.options.initializer != nil {
			panic("can't override initializer for KMEANS_PLUS_PLUS")
		}
		return clusterer.NewKmeansPlusPlus(b.options.vectors, b.options.clusterCnt,
			b.options.deltaThreshold, b.options.iterationThreshold,
			b.options.distanceFn)
	default:
		panic("invalid cluster type")
	}
}
