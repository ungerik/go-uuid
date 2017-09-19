package uuid

import (
	"bytes"
	"sort"
	"testing"
)

func Test_Less(t *testing.T) {
	input := make([]UUID, 64)
	for i := range input {
		input[i] = NewV4()
	}
	sortedRef := make([]UUID, len(input))
	for i := range input {
		sortedRef[i] = input[i]
	}
	sort.Slice(sortedRef, func(i, j int) bool { return bytes.Compare(sortedRef[i][:], sortedRef[j][:]) == -1 })

	sorted := make([]UUID, len(input))
	for i := range input {
		sorted[i] = input[i]
	}
	sort.Slice(sorted, func(i, j int) bool { return Less(sorted[i], sorted[j]) })

	for i := range input {
		if sorted[i] != sortedRef[i] {
			t.Fatalf("%v != %v", sorted, sortedRef)
		}
	}
}
