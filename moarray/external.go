// Copyright 2023 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package moarray

import (
	"github.com/arjunsk/kmeans/moerr"
	"golang.org/x/exp/constraints"
	"gonum.org/v1/gonum/mat"
	"math"
)

// Compare returns an integer comparing two arrays/vectors lexicographically.
// TODO: this function might not be correct. we need to compare using tolerance for float values.
// TODO: need to check if we need len(v1)==len(v2) check.
func Compare[T constraints.Float](v1, v2 []T) int {
	minLen := len(v1)
	if len(v2) < minLen {
		minLen = len(v2)
	}

	for i := 0; i < minLen; i++ {
		if v1[i] < v2[i] {
			return -1
		} else if v1[i] > v2[i] {
			return 1
		}
	}

	if len(v1) < len(v2) {
		return -1
	} else if len(v1) > len(v2) {
		return 1
	}
	return 0
}

/* ------------ [START] Performance critical functions. ------- */

func InnerProduct[T constraints.Float](v1, v2 []T) (float64, error) {

	vec, err := ToGonumVectors[T](v1, v2)
	if err != nil {
		return 0, err
	}

	return mat.Dot(vec[0], vec[1]), nil
}

func L2Distance[T constraints.Float](v1, v2 []T) (float64, error) {
	vec, err := ToGonumVectors[T](v1, v2)
	if err != nil {
		return 0, err
	}

	diff := mat.NewVecDense(vec[0].Len(), nil)
	diff.SubVec(vec[0], vec[1])

	return math.Sqrt(mat.Dot(diff, diff)), nil
}

func CosineDistance[T constraints.Float](v1, v2 []T) (float64, error) {
	cosineSimilarity, err := CosineSimilarity[T](v1, v2)
	if err != nil {
		return 0, err
	}

	return 1 - cosineSimilarity, nil
}

func CosineSimilarity[T constraints.Float](v1, v2 []T) (float64, error) {

	vec, err := ToGonumVectors[T](v1, v2)
	if err != nil {
		return 0, err
	}

	dotProduct := mat.Dot(vec[0], vec[1])

	normVec1 := mat.Norm(vec[0], 2)
	normVec2 := mat.Norm(vec[1], 2)

	if normVec1 == 0 || normVec2 == 0 {
		return 0, moerr.NewInternalErrorNoCtx("cosine_similarity: one of the vectors is zero")
	}

	cosineSimilarity := dotProduct / (normVec1 * normVec2)

	// Handle precision issues. Clamp the cosine_similarity to the range [-1, 1].
	if cosineSimilarity > 1.0 {
		cosineSimilarity = 1.0
	} else if cosineSimilarity < -1.0 {
		cosineSimilarity = -1.0
	}

	// NOTE: Downcast the float64 cosine_similarity to float32 and check if it is
	// 1.0 or -1.0 to avoid precision issue.
	//
	//  Example for corner case:
	// - cosine_similarity(a,a) = 1:
	// - Without downcasting check, we get the following results:
	//   cosine_similarity( [0.46323407, 23.498016, 563.923, 56.076736, 8732.958] ,
	//					    [0.46323407, 23.498016, 563.923, 56.076736, 8732.958] ) =   0.9999999999999998
	// - With downcasting, we get the following results:
	//   cosine_similarity( [0.46323407, 23.498016, 563.923, 56.076736, 8732.958] ,
	//					    [0.46323407, 23.498016, 563.923, 56.076736, 8732.958] ) =   1
	//
	//  Reason:
	// The reason for this check is
	// 1. gonums mat.Dot, mat.Norm returns float64. In other databases, we mostly do float32 operations.
	// 2. float64 operations are not exact.
	// mysql> select 76586261.65813679/(8751.35770370157 *8751.35770370157);
	//+-----------------------------------------------------------+
	//| 76586261.65813679 / (8751.35770370157 * 8751.35770370157) |
	//+-----------------------------------------------------------+
	//|                                            1.000000000000 |
	//+-----------------------------------------------------------+
	//mysql> select cast(76586261.65813679 as double)/(8751.35770370157 * 8751.35770370157);
	//+---------------------------------------------------------------------------+
	//| cast(76586261.65813679 as double) / (8751.35770370157 * 8751.35770370157) |
	//+---------------------------------------------------------------------------+
	//|                                                        0.9999999999999996 |
	//+---------------------------------------------------------------------------+
	// 3. We only need to handle the case for 1.0 and -1.0 with float32 precision.
	//    Rest of the cases can have float64 precision.
	cosineSimilarityF32 := float32(cosineSimilarity)
	if cosineSimilarityF32 == 1 {
		cosineSimilarity = 1
	} else if cosineSimilarityF32 == -1 {
		cosineSimilarity = -1
	}

	return cosineSimilarity, nil
}

func NormalizeL2[T constraints.Float](v1 []T) ([]T, error) {

	vec := ToGonumVector[T](v1)

	norm := mat.Norm(vec, 2)
	if norm == 0 {
		// NOTE: don't throw error here. If you throw error, then when a zero vector comes in the Vector Index
		// Mapping Query, the query will fail. Instead, return the same zero vector.
		// This is consistent with FAISS:https://github.com/facebookresearch/faiss/blob/0716bde2500edb2e18509bf05f5dfa37bd698082/faiss/utils/distances.cpp#L97
		return v1, nil
	}

	vec.ScaleVec(1/norm, vec)

	return ToMoArray[T](vec), nil
}
