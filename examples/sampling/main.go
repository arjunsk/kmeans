package main

import (
	"fmt"
	"github.com/arjunsk/kmeans"
	"github.com/arjunsk/kmeans/sampler"
)

func main() {
	vectors := [][]float64{
		{1, 2, 3, 4},
		{0, 3, 4, 1},
		{0, 9, 3, 1},
		{0, 8, 4, 4},
		{130, 200, 343, 224},
		{100, 200, 300, 400},
		{300, 400, 200, 110},
	}

	// sub-sample
	var sampleFn sampler.Sampling[[]float64] = sampler.SrsSampling[[]float64]
	subset := sampleFn(vectors, 50)

	// run clustering
	clusterer, err := kmeans.NewCluster(kmeans.ELKAN, subset, 2)
	if err != nil {
		panic(err)
	}

	clusters, err := clusterer.Cluster()
	if err != nil {
		panic(err)
	}

	for _, cluster := range clusters {
		fmt.Println(cluster.Center())
	}
	// Output:
	// [1 2 3 4]
	// [130 200 343 224]

}
