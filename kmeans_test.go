package go_kmeans

import (
	"fmt"
	"github.com/arjunsk/go-kmeans/clusterer"
	"testing"
)

var vectors = [][]float64{
	{20.0, 20.0, 20.0, 20.0},
	{21.0, 21.0, 21.0, 21.0},
	{100.5, 100.5, 100.5, 100.5},
	{50.1, 50.1, 50.1, 50.1},
	{64.2, 64.2, 64.2, 64.2},
}

func TestTrain_lloyd(t *testing.T) {
	kmeans, _ := clusterer.NewKmeans(vectors, 2)
	clusters, err := kmeans.Cluster()
	if err != nil || clusters == nil || len(clusters) != 2 {
		t.Log("\nClusters:", clusters)
		t.Fail()
	}
	fmt.Println(clusters)
	fmt.Println(clusters.SSE())
}

func TestTrain_kpp(t *testing.T) {
	kmeans, _ := clusterer.NewKmeansPlusPlus(vectors, 2)
	clusters, err := kmeans.Cluster()
	if err != nil || clusters == nil || len(clusters) != 2 {
		t.Log("\nClusters:", clusters)
		t.Fail()
	}
	fmt.Println(clusters)
	fmt.Println(clusters.SSE())
}

func TestTrain_elkan(t *testing.T) {
	kmeans, _ := clusterer.NewKmeansElkan(vectors, 2)
	clusters, err := kmeans.Cluster()
	if err != nil || clusters == nil || len(clusters) != 2 {
		t.Log("\nError:", err)
		t.Log("\nClusters:", clusters)
		t.Fail()
	}
	fmt.Println(clusters)
	fmt.Println(clusters.SSE())
}
