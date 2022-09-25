package codec

func SliceEqual[T comparable](left, right []T) bool {
	if len(left) != len(right) {
		return false
	}

	for i, elem := range left {
		if elem != right[i] {
			return false
		}
	}

	return true
}
