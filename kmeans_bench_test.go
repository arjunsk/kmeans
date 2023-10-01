package go_clustering

import (
	"go-kmeans/clusterer"
	"go-kmeans/domain"
	"math/rand"
	"testing"
)

func Benchmark_kmeans(b *testing.B) {
	rowCnt := 1_000
	dims := 1024
	data := make([]domain.Vector, rowCnt)
	loadData(rowCnt, dims, data)

	kmeans, _ := clusterer.NewKmeans(data, 100)
	kmeanspp, _ := clusterer.NewKmeansPlusPlus(data, 100)
	elkan, _ := clusterer.NewKmeansElkan(data, 100)

	b.Run("KMEANS", func(b *testing.B) {
		b.ResetTimer()
		for i := 1; i < b.N; i++ {
			_, err := kmeans.Cluster()
			if err != nil {
				panic(err)
			}
		}
	})

	b.Run("KMEANS++", func(b *testing.B) {
		b.ResetTimer()
		for i := 1; i < b.N; i++ {
			_, err := kmeanspp.Cluster()
			if err != nil {
				panic(err)
			}
		}
	})

	b.Run("ELKAN", func(b *testing.B) {
		b.ResetTimer()
		for i := 1; i < b.N; i++ {
			_, err := elkan.Cluster()
			if err != nil {
				panic(err)
			}
		}
	})
}

func loadData(nb int, d int, xb []domain.Vector) {
	for r := 0; r < nb; r++ {
		xb[r] = make(domain.Vector, d)
		for c := 0; c < d; c++ {
			xb[r][c] = float64(rand.Float32() * 1000)
		}
	}
}
