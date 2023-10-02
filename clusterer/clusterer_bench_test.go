package clusterer

import (
	"math/rand"
	"testing"
)

/*
date : 2023-10-2
goos: darwin
goarch: arm64
pkg: github.com/arjunsk/kmeans/clusterer
cpu: Apple M2 Pro
Benchmark_kmeans/KMEANS-10         	       1	4426765042 ns/op
Benchmark_kmeans/KMEANS++-10       	       1	8420956083 ns/op
Benchmark_kmeans/ELKAN-10          	       1	4060144125 ns/op
*/
func Benchmark_kmeans(b *testing.B) {
	rowCnt := 1_000
	dims := 1024
	data := make([][]float64, rowCnt)
	loadData(rowCnt, dims, data)

	kmeans, _ := NewKmeans(data, 100)
	kmeanspp, _ := NewKmeansPlusPlus(data, 100)
	elkan, _ := NewKmeansElkan(data, 100)

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
