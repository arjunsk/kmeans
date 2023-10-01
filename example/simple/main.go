package main

import (
	"fmt"
	"github.com/arjunsk/go-kmeans/clusterer"
)

func main() {
	vectors := [][]float64{
		{1, 2, 3, 4},
		{0, 3, 4, 1},
		{130, 200, 343, 224},
		{100, 200, 300, 400},
	}

	kmeans, err := clusterer.NewKmeansElkan(vectors, 2)
	if err != nil {
		panic(err)
	}

	clusters, err := kmeans.Cluster()
	if err != nil {
		panic(err)
	}

	for _, cluster := range clusters {
		fmt.Println(cluster.Center)
	}
	// Output:
	// [1 2 3 4]
	// [130 200 343 224]

}
