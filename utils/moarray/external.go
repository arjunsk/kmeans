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
	"golang.org/x/exp/constraints"
	"gonum.org/v1/gonum/mat"
)

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
