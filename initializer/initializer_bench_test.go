package initializer

import (
	"github.com/arjunsk/kmeans/containers"
	"math/rand"
	"testing"
)

/*
date : 2023-10-1
goos: darwin
goarch: arm64
pkg: github.com/arjunsk/kmeans/initializer
cpu: Apple M2 Pro
Benchmark_kmeans/RANDOM-10         	  104032	     11292 ns/op
Benchmark_kmeans/KMEANS++-10       	       1	3840350291 ns/op
*/
func Benchmark_kmeans(b *testing.B) {
	rowCnt := 1_000
	dims := 1024
	data := make([][]float64, rowCnt)
	loadData(rowCnt, dims, data)

	random := NewRandomInitializer()
	kmeanspp := NewKmeansPlusPlusInitializer(containers.EuclideanDistance)

	b.Run("RANDOM", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := random.InitCentroids(data, 100)
			if err != nil {
				panic(err)
			}
		}
	})

	b.Run("KMEANS++", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := kmeanspp.InitCentroids(data, 100)
			if err != nil {
				panic(err)
			}
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
