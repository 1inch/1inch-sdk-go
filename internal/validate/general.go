package validate

// HasDuplicates checks if the provided slice contains any duplicate elements.
// It accepts a slice of any comparable type and returns true if there are duplicates, otherwise it returns false.
func HasDuplicates[T comparable](slice []T) bool {
	seen := make(map[T]bool)
	for _, v := range slice {
		if seen[v] {
			return true
		}
		seen[v] = true
	}
	return false
}

// IsSubset checks if all elements of sliceA are also present in sliceB.
// It returns true if sliceA is a subset of sliceB, otherwise it returns false.
func IsSubset[T comparable](sliceA, sliceB []T) bool {
	setB := make(map[T]bool)
	for _, v := range sliceB {
		setB[v] = true
	}

	for _, v := range sliceA {
		if !setB[v] {
			return false
		}
	}
	return true
}
