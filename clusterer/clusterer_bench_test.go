package clusterer

import (
	"github.com/arjunsk/kmeans/containers"
	"github.com/arjunsk/kmeans/initializer"
	"math/rand"
	"testing"
)

/*
date : 2023-10-2
goos: darwin
goarch: arm64
pkg: github.com/arjunsk/kmeans/clusterer
cpu: Apple M2 Pro
rows: 5_000
dims: 1024
Benchmark_kmeans/KMEANS-10         	       1	32952745208 ns/op
Benchmark_kmeans/KMEANS++-10       	       1	89310750250 ns/op
Benchmark_kmeans/ELKAN-10          	       1	21269951792 ns/op
*/
func Benchmark_kmeans(b *testing.B) {
	rowCnt := 5_000
	dims := 1024
	data := make([][]float64, rowCnt)
	loadData(rowCnt, dims, data)

	kmeans, _ := NewKmeans(data, 100,
		0.01, 500,
		containers.EuclideanDistance, initializer.NewRandomInitializer())
	kmeanspp, _ := NewKmeansPlusPlus(data, 100,
		0.01, 500,
		containers.EuclideanDistance)
	elkan, _ := NewKmeansElkan(data, 100,
		0.01, 500,
		containers.EuclideanDistance, initializer.NewRandomInitializer())

	b.Run("KMEANS", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			clusters, err := kmeans.Cluster()
			if err != nil {
				panic(err)
			}
			b.Logf("SSE %v", clusters.SSE())
		}
	})

	b.Run("KMEANS++", func(b *testing.B) {
		b.Skipf("KMEANS++ will take alot of time for k=100. Hence skipping")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			clusters, err := kmeanspp.Cluster()
			if err != nil {
				panic(err)
			}
			b.Logf("SSE %v", clusters.SSE())
		}
	})

	b.Run("ELKAN", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			clusters, err := elkan.Cluster()
			if err != nil {
				panic(err)
			}
			b.Logf("SSE %v", clusters.SSE())
		}
	})
}

func loadData(nb int, d int, xb [][]float64) {
	for r := 0; r < nb; r++ {
		xb[r] = make([]float64, d)
		for c := 0; c < d; c++ {
			xb[r][c] = float64(rand.Float32() * 1000)
		}
	}
}
