## Go Kmeans

This is a simple implementation of the Elkan's Kmeans algorithm in Go. It is based on
the [Kmeans++](https://en.wikipedia.org/wiki/K-means%2B%2B) algorithm for the initial centroids
and the [Elkan's](https://cdn.aaai.org/ICML/2003/ICML03-022.pdf) algorithm for the clustering.

### Usage

```go
package main

import (
	"fmt"
	"github.com/arjunsk/go-kmeans/clusterer"
	"github.com/arjunsk/go-kmeans/domain"
)

func main() {
	vectors := []domain.Vector{
		{1, 2, 3, 4},
		{0, 3, 4, 1},
		{130, 200, 343, 224},
		{100, 200, 300, 400},
	}

	kmeans, err := clusterer.NewKmeansElkan(vectors, 2)
	if err != nil {
		panic(err)
	}

	clusterGroup, err := kmeans.Cluster()
	if err != nil {
		panic(err)
	}

	for _, cluster := range clusterGroup {
		fmt.Println(cluster.Center)
	}
}

```