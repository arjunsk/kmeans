package clusterer

import (
	"github.com/arjunsk/kmeans/containers"
	"github.com/arjunsk/kmeans/initializer"
)

type KmeansPP struct {
	Lloyd
}

var _ Clusterer = new(KmeansPP)

func NewKmeansPlusPlus(vectors [][]float64, clusterCnt int,
	deltaThreshold float64,
	iterationThreshold int,
	distFn containers.DistanceFunction) (Clusterer, error) {

	clusterer := Lloyd{
		deltaThreshold:     deltaThreshold,
		iterationThreshold: iterationThreshold,
		distFn:             distFn,
		initializer:        initializer.NewKmeansPlusPlusInitializer(distFn),
		vectors:            vectors,
		clusterCnt:         clusterCnt,
	}

	return &KmeansPP{
		Lloyd: clusterer,
	}, nil

}

func (kpp KmeansPP) Cluster() (containers.Clusters, error) {
	return kpp.Lloyd.Cluster()
}
