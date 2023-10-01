package containers

type Vector []float64

// Add adds a vector to the current vector
func (v *Vector) Add(vec Vector) {
	for c, val := range vec {
		(*v)[c] += val
	}
}

// Mul multiplies the vector by a scalar
func (v *Vector) Mul(scalar float64) {
	for c := range *v {
		(*v)[c] *= scalar
	}
}

// Compare 2 vectors on a lexicographical order
// Ref: https://stackoverflow.com/a/23907444/1609570
func (v *Vector) Compare(v2 Vector) int {
	lenV1 := len(*v)
	lenV2 := len(v2)

	minLen := lenV1
	if lenV2 < lenV1 {
		minLen = lenV2
	}

	for c := 0; c < minLen; c++ {
		if (*v)[c] < v2[c] {
			return -1
		} else if (*v)[c] > v2[c] {
			return +1
		}
	}

	if lenV1 < lenV2 {
		return -1
	} else if lenV1 > lenV2 {
		return +1
	}

	return 0
}
