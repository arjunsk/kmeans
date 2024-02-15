package main

import (
	"fmt"
	"github.com/arjunsk/kmeans"
	"github.com/arjunsk/kmeans/elkans"
)

func main() {
	vectorList := [][]float64{
		{1, 2, 3, 4},
		{1, 2, 4, 5},
		{1, 2, 4, 5},
		{1, 2, 3, 4},
		{1, 2, 4, 5},
		{1, 2, 4, 5},
		{10, 2, 4, 5},
		{10, 3, 4, 5},
		{10, 5, 4, 5},
		{10, 2, 4, 5},
		{10, 3, 4, 5},
		{10, 5, 4, 5},
	}

	clusterer, err := elkans.NewKMeans(vectorList, 2,
		500, 0.5,
		kmeans.L2Distance, kmeans.KmeansPlusPlus, false)
	if err != nil {
		panic(err)
	}

	centroids, err := clusterer.Cluster()
	if err != nil {
		panic(err)
	}

	for _, centroid := range centroids {
		fmt.Println(centroid)
	}
}
