package util

// ArrayCloneAndReverse returns a copy of the given byte slice in reversed order.
func ArrayCloneAndReverse(b []byte) []byte {
	dest := make([]byte, len(b))
	for i, j := 0, len(b)-1; i <= j; i, j = i+1, j-1 {
		dest[i], dest[j] = b[j], b[i]
	}
	return dest
}
