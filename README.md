# Go Kmeans

This is a simple implementation of the [Elkan's Kmeans](https://cdn.aaai.org/ICML/2003/ICML03-022.pdf) 
algorithm in Go.

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

### Why not Kmeans++ initialization?

Based on the excerpt
from [FAISS discussion](https://github.com/facebookresearch/faiss/issues/268#issuecomment-348184505), it was observed
that Kmeans++ overhead computation cost is not worth for large dimension vectors.

> Scikitlearn uses k-means++ initialization by default (you can also use random points), which is good in the specific
> corner-case you consider. It should actually gives you perfect result even without any iteration with high probability,
> because the kind of evaluation you consider is exactly what k-means++ has be designed to better handle.
> We have not implemented it in Faiss, because with our former Yael library, which implements both k-means++ and regular
> random initialization, we observed that the overhead computational cost was not worth the saving (negligible) in all
> large-scale settings we have considered.
