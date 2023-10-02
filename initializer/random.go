package initializer

import (
	"github.com/arjunsk/kmeans/containers"
	"math/rand"
	"time"
)

type Random struct{}

func NewRandomInitializer() Initializer {
	return &Random{}
}

func (k *Random) InitCentroids(vectors [][]float64, clusterCnt int) (containers.Clusters, error) {
	err := validateArgs(vectors, clusterCnt)
	if err != nil {
		return nil, err
	}

	inputCnt := len(vectors)
	var clusters containers.Clusters = make([]containers.Cluster, clusterCnt)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < clusterCnt; i++ {
		randIdx := random.Intn(inputCnt)
		clusters[i] = containers.Cluster{
			Center: vectors[randIdx],
		}
	}

	return clusters, nil
}
