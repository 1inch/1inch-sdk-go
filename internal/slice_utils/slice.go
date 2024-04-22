package slice_utils

// Contains checks if the slice contains the given value.
func Contains[T comparable](value T, sliceB []T) bool {
	for _, v := range sliceB {
		if v == value {
			return true
		}
	}
	return false
}
