package util

// ArrayCloneAndReverse returns a copy of the given byte slice in reversed order.
func ArrayCloneAndReverse[T any](a []T) []T {
	dest := make([]T, len(a))
	for i, j := 0, len(a)-1; i <= j; i, j = i+1, j-1 {
		dest[i], dest[j] = a[j], a[i]
	}
	return dest
}

// Returns copy of given array
func ArrayClone[T any](a []T) []T {
	clone := make([]T, len(a))
	copy(clone, a)

	return clone
}
