package kmeans

import (
	"fmt"
	"github.com/arjunsk/kmeans/clusterer"
	"github.com/arjunsk/kmeans/containers"
	"github.com/arjunsk/kmeans/initializer"
)

var (
	defaultDeltaThreshold     = 0.01
	defaultIterationThreshold = 500
)

type ClustererType int

const (
	ELKAN ClustererType = iota
	KMEANS
	KMEANSPP // Kmeans Plus Plus
)

type options struct {
	clusterType        ClustererType
	vectors            [][]float64
	clusterCnt         int
	distanceFn         containers.DistanceFunction
	initializer        initializer.Initializer
	deltaThreshold     *float64
	iterationThreshold *int
}

func (o *options) fillDefaults() {
	if o.initializer == nil && o.clusterType != KMEANSPP {
		o.initializer = initializer.NewRandomInitializer()
	}
	if o.distanceFn == nil {
		o.distanceFn = containers.EuclideanDistance
	}
	if o.deltaThreshold == nil {
		o.deltaThreshold = &defaultDeltaThreshold
	}
	if o.iterationThreshold == nil {
		o.iterationThreshold = &defaultIterationThreshold
	}

}

type Option func(options *options) error

// NewCluster constructs an options instance with the provided functional options.
func NewCluster(clustererType ClustererType, vectors [][]float64, clusterCnt int, opts ...Option) (clusterer.Clusterer, error) {
	o := &options{
		clusterType: clustererType,
		vectors:     vectors,
		clusterCnt:  clusterCnt,
	}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	o.fillDefaults()

	switch o.clusterType {
	case ELKAN:
		return clusterer.NewKmeansElkan(o.vectors, o.clusterCnt,
			*o.deltaThreshold, *o.iterationThreshold,
			o.distanceFn, o.initializer)
	case KMEANS:
		return clusterer.NewKmeans(o.vectors, o.clusterCnt,
			*o.deltaThreshold, *o.iterationThreshold,
			o.distanceFn, o.initializer)
	case KMEANSPP:
		if o.initializer != nil {
			panic("can't override initializer for KMEANS_PLUS_PLUS")
		}
		return clusterer.NewKmeansPlusPlus(o.vectors, o.clusterCnt,
			*o.deltaThreshold, *o.iterationThreshold,
			o.distanceFn)
	default:
		return nil, fmt.Errorf("invalid cluster type %v", o.clusterType)
	}
}

// WithDistanceFunction sets the distance function in options.
func WithDistanceFunction(fn containers.DistanceFunction) Option {
	return func(o *options) error {
		o.distanceFn = fn
		return nil
	}
}

// WithInitializer sets the initializer in options.
func WithInitializer(init initializer.Initializer) Option {
	return func(o *options) error {
		o.initializer = init
		return nil
	}
}

// WithDeltaThreshold sets the delta threshold in options.
func WithDeltaThreshold(delta float64) Option {
	return func(o *options) error {
		o.deltaThreshold = &delta
		return nil
	}
}

// WithIterationThreshold sets the iteration threshold in options.
func WithIterationThreshold(iterations int) Option {
	return func(o *options) error {
		o.iterationThreshold = &iterations
		return nil
	}
}
