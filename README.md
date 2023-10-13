# Go Kmeans

[![Go Reference](https://pkg.go.dev/badge/github.com/arjunsk/kmeans/kmeans.svg)](https://pkg.go.dev/github.com/arjunsk/kmeans)
[![Go Report Card](https://goreportcard.com/badge/github.com/arjunsk/kmeans)](https://goreportcard.com/report/github.com/arjunsk/kmeans)

This is a simple implementation of the [Elkan's Kmeans](https://cdn.aaai.org/ICML/2003/ICML03-022.pdf)
algorithm in Go. The library also contains [Kmeans++](https://en.wikipedia.org/wiki/K-means%2B%2B),
[Lloyd's kmeans](https://en.wikipedia.org/wiki/K-means_clustering#Standard_algorithm_(naive_k-means)) and
[Simple Random Sampling](https://en.wikipedia.org/wiki/Simple_random_sample) algorithms.

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

	clusterer, err := kmeans.NewCluster(kmeans.ELKAN, vectors, 2)
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
```

### FAQ
<details>
<summary> Read More </summary>

#### Why not Kmeans++ initialization in Elkan's?

The default settings of Elkan's Kmeans is to use [random initialization](/initializer/random.go)
instead of  [Kmeans++ initialization](/initializer/kmeans_plus_plus.go).

Based on the excerpt
from [FAISS discussion](https://github.com/facebookresearch/faiss/issues/268#issuecomment-348184505), it was observed
that Kmeans++ overhead computation cost is not worth for large scale use case.

> Scikitlearn uses k-means++ initialization by default (you can also use random points), which is good in the specific
> corner-case you consider. It should actually gives you perfect result even without any iteration with high
> probability,
> because the kind of evaluation you consider is exactly what k-means++ has be designed to better handle.
> We have not implemented it in Faiss, because with our former Yael library, which implements both k-means++ and regular
> random initialization, we observed that the overhead computational cost was not worth the saving (negligible) in all
> large-scale settings we have considered.

#### When should you consider sub-sampling?

As mentioned [here](https://github.com/facebookresearch/faiss/wiki/FAQ#can-i-ignore-warning-clustering-xxx-points-to-yyy-centroids),
when the number of vectors is large, it is recommended to use sub-sampling.


> When applying k-means algorithm to cluster n points to k centroids, there are several cases:
>
> - n < k: this raises an exception with an assertion because we cannot do anything meaningful
> - n < min_points_per_centroid * k: this produces the warning above. It means that usually there are too few points to
    reliably estimate the centroids. This may still be ok if the dataset to index is as small as the training set.
> - n < max_points_per_centroid * k: comfort zone
> - n > max_points_per_centroid * k: there are too many points, making k-means unnecessarily slow. Then the training set
    is sampled.
>
>The parameters {min,max}_points_per_centroids (39 and 256 by default) belong to the ClusteringParameters structure.

#### What should be the ideal K?
Based on the recommendations from [PGVector](https://github.com/pgvector/pgvector/tree/master#ivfflat) IVF INDEX, 
the idea K should 

> Choose an appropriate number of K - a good place to start is rows / 1000 for up to 1M rows and 
> sqrt(rows) for over 1M rows



</details>
