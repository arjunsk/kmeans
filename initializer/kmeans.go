package initializer

import (
	"go-kmeans/domain"
	"math/rand"
	"time"
)

type Kmeans struct{}

func NewKmeansInitializer() Initializer {
	return &Kmeans{}
}

func (k *Kmeans) InitCentroids(vectors []domain.Vector, clusterCnt int) (domain.Clusters, error) {
	inputCnt := len(vectors)

	err := StdInputChecks(vectors, clusterCnt, inputCnt)
	if err != nil {
		return nil, err
	}

	var clusters domain.Clusters = make([]domain.Cluster, clusterCnt)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < clusterCnt; i++ {
		randIdx := random.Intn(inputCnt)
		clusters[i] = domain.Cluster{
			Center: vectors[randIdx],
		}
	}

	return clusters, nil
}
