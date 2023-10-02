## Go Kmeans

This is a simple implementation of the Elkan's Kmeans algorithm in Go. It is based on
the [Kmeans++](https://en.wikipedia.org/wiki/K-means%2B%2B) algorithm for the initial centroids
and the [Elkan's](https://cdn.aaai.org/ICML/2003/ICML03-022.pdf) algorithm for the clustering.

### Usage

```go
package main

import (
	"fmt"
	"github.com/arjunsk/kmeans"
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

	clusterer, err := kmeans.NewKmeansElkan(vectors, 2)
	if err != nil {
		panic(err)
	}

	clusters, err := clusterer.Cluster()
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
```