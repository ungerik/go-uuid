package uuid

import "bytes"

func Less(a, b UUID) bool {
	return bytes.Compare(a[:], b[:]) < 0
	// for i := len(a) - 1; i >= 0; i-- {
	// 	if a[i] < b[i] {
	// 		return true
	// 	}
	// 	if a[i] > b[i] {
	// 		return false
	// 	}
	// }
	// return false
}
