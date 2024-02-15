# Go Kmeans

[![Go Reference](https://pkg.go.dev/badge/github.com/arjunsk/kmeans/kmeans.svg)](https://pkg.go.dev/github.com/arjunsk/kmeans)
[![Go Report Card](https://goreportcard.com/badge/github.com/arjunsk/kmeans)](https://goreportcard.com/report/github.com/arjunsk/kmeans)
[![Codecov](https://codecov.io/gh/arjunsk/kmeans/branch/master/graph/badge.svg)](https://codecov.io/gh/arjunsk/kmeans)


This is a simple implementation of the [Elkan's Kmeans](https://cdn.aaai.org/ICML/2003/ICML03-022.pdf)
algorithm in Go.

### Installing

```sh
$ go get github.com/arjunsk/kmeans
```

### Usage

```go
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
	/*
	[1 2 3.6666666666666665 4.666666666666666]
	[10 3.333333333333333 4 5]
	*/
}
```

### FAQ
<details>
<summary> Read More </summary>

#### What should be the ideal Centroids Count?
Based on the recommendations from [PGVector](https://github.com/pgvector/pgvector/tree/master#ivfflat) IVF INDEX, 
the idea K should 

> Choose an appropriate number of K - a good place to start is rows / 1000 for up to 1M rows and 
> sqrt(rows) for over 1M rows



</details>
