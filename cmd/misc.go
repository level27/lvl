package cmd

// Finds the index of an element within a slice. Returns -1 if the element is not present.
func indexOf[E comparable](slice []E, value E) int {
	// Replace with golang.org/x/exp/slices/Index if/when that ever becomes non-experimental.
	for idx, v := range slice {
		if v == value {
			return idx
		}
	}

	return -1
}
