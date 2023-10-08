package main

import (
	"fmt"
	"github.com/arjunsk/kmeans"
	"github.com/arjunsk/kmeans/containers"
	"github.com/arjunsk/kmeans/initializer"
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

	builder := kmeans.NewClusterBuilder(kmeans.ELKAN, vectors, 2)
	builder.Initializer(initializer.NewKmeansPlusPlusInitializer(containers.EuclideanDistance))
	/*
		To define custom initializer, you can use this syntax:

		type Custom struct {}

		func (c *Custom) InitCentroids(vectors [][]float64, clusterCnt int) (containers.Clusters, error) {
			panic("implement me")
		}

		var init initializer.Initializer = &Custom{}
		builder.Initializer(init)
	*/

	clusterer, err := builder.Build()
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
	// [0.25 5.5 3.5 2.5]
	// [176.66666666666666 266.66666666666663 281 244.66666666666666]
}
