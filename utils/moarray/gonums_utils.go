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
	"github.com/arjunsk/kmeans/utils/moerr"
	"golang.org/x/exp/constraints"
	"gonum.org/v1/gonum/mat"
)

func ToGonumVector[T constraints.Float](arr1 []T) *mat.VecDense {

	n := len(arr1)
	_arr1 := make([]float64, n)

	//TODO: @arjun optimize this cast to retain float32 precision in float64 array
	// if float64, just copy
	// if float32, convert to float64 without losing precision
	for i := 0; i < n; i++ {
		_arr1[i] = float64(arr1[i])
	}

	return mat.NewVecDense(n, _arr1)
}

func ToGonumVectors[T constraints.Float](arrays ...[]T) (res []*mat.VecDense, err error) {

	n := len(arrays)
	if n == 0 {
		return res, nil
	}

	array0Dim := len(arrays[0])
	for i := 1; i < n; i++ {
		if len(arrays[i]) != array0Dim {
			return nil, moerr.NewArrayInvalidOpNoCtx(array0Dim, len(arrays[i]))
		}
	}

	res = make([]*mat.VecDense, n)

	for i, arr := range arrays {
		res[i] = ToGonumVector[T](arr)
	}
	return res, nil
}

func ToMoArray[T constraints.Float](vec *mat.VecDense) (arr []T) {
	n := vec.Len()
	arr = make([]T, n)
	for i := 0; i < n; i++ {
		//TODO: @arjun optimize this cast
		arr[i] = T(vec.AtVec(i))
	}
	return
}

func ToMoArrays[T constraints.Float](vecs []*mat.VecDense) [][]T {
	moVectors := make([][]T, len(vecs))
	for i, vec := range vecs {
		moVectors[i] = ToMoArray[T](vec)
	}
	return moVectors
}
