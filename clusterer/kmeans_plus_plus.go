package clusterer

import (
	"github.com/arjunsk/go-kmeans/domain"
	"github.com/arjunsk/go-kmeans/initializer"
)

type KmeansPP struct {
	Lloyd
}

var _ Clusterer = new(KmeansPP)

func NewKmeansPlusPlus(vectors []domain.Vector, clusterCnt int) (Clusterer, error) {

	clusterer, err := newKmeansWithOptions(
		0.01,
		500,
		domain.EuclideanDistance,
		initializer.NewKmeansPlusPlusInitializer(domain.EuclideanDistance))
	if err != nil {
		return nil, err
	}

	clusterer.vectors = vectors
	clusterer.clusterCnt = clusterCnt

	return &KmeansPP{
		Lloyd: clusterer,
	}, nil

}
