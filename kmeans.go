package go_kmeans

import "github.com/arjunsk/go-kmeans/clusterer"

//TODO: Add more options using Builder pattern

func NewKmeans(vectors [][]float64, clusterCnt int) (clusterer.Clusterer, error) {
	return clusterer.NewKmeans(vectors, clusterCnt)
}

func NewKmeansPlusPlus(vectors [][]float64, clusterCnt int) (clusterer.Clusterer, error) {
	return clusterer.NewKmeansPlusPlus(vectors, clusterCnt)
}

func NewKmeansElkan(vectors [][]float64, clusterCnt int) (clusterer.Clusterer, error) {
	return clusterer.NewKmeansElkan(vectors, clusterCnt)
}
