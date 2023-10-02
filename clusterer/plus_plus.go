package clusterer

import (
	"github.com/arjunsk/kmeans/containers"
	"github.com/arjunsk/kmeans/initializer"
)

type KmeansPP struct {
	Lloyd
}

var _ Clusterer = new(KmeansPP)

func NewKmeansPlusPlus(vectors [][]float64, clusterCnt int) (Clusterer, error) {

	clusterer, err := newKmeansWithOptions(
		0.01,
		500,
		containers.EuclideanDistance,
		initializer.NewKmeansPlusPlusInitializer(containers.EuclideanDistance))
	if err != nil {
		return nil, err
	}

	clusterer.vectors = vectors
	clusterer.clusterCnt = clusterCnt

	return &KmeansPP{
		Lloyd: clusterer,
	}, nil

}

func (kpp KmeansPP) Cluster() (containers.Clusters, error) {
	return kpp.Lloyd.Cluster()
}