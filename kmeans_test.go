package go_clustering

import (
	"fmt"
	"go-kmeans/clusterer"
	"go-kmeans/domain"
	"testing"
)

var vectors = []domain.Vector{
	{20.0, 20.0, 20.0, 20.0},
	{21.0, 21.0, 21.0, 21.0},
	{100.5, 100.5, 100.5, 100.5},
	{50.1, 50.1, 50.1, 50.1},
	{64.2, 64.2, 64.2, 64.2},
}

func TestTrain_lloyd(t *testing.T) {
	kmeans, _ := clusterer.NewKmeans(vectors, 2)
	clusterGroup, err := kmeans.Cluster()
	if err != nil || clusterGroup == nil || len(clusterGroup) != 2 {
		t.Log("\nClusters:", clusterGroup)
		t.Fail()
	}
	fmt.Println(clusterGroup)
}

func TestTrain_kpp(t *testing.T) {
	kmeans, _ := clusterer.NewKmeansPlusPlus(vectors, 2)
	clusterGroup, err := kmeans.Cluster()
	if err != nil || clusterGroup == nil || len(clusterGroup) != 2 {
		t.Log("\nClusters:", clusterGroup)
		t.Fail()
	}
	fmt.Println(clusterGroup)
}

func TestTrain_elkan(t *testing.T) {
	kmeans, _ := clusterer.NewKmeansElkan(vectors, 2)
	clusterGroup, err := kmeans.Cluster()
	if err != nil || clusterGroup == nil || len(clusterGroup) != 2 {
		t.Log("\nError:", err)
		t.Log("\nClusters:", clusterGroup)
		t.Fail()
	}
	fmt.Println(clusterGroup)
}
