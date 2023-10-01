package initializer

import (
	"github.com/arjunsk/go-kmeans/containers"
	"math/rand"
	"time"
)

type Kmeans struct{}

func NewKmeansInitializer() Initializer {
	return &Kmeans{}
}

func (k *Kmeans) InitCentroids(vectors []containers.Vector, clusterCnt int) (containers.Clusters, error) {
	inputCnt := len(vectors)

	err := StdInputChecks(vectors, clusterCnt, inputCnt)
	if err != nil {
		return nil, err
	}

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
