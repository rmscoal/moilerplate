package utils

import "unsafe"

// ConvertStringToByteSlice converts a string to slice of byte in an unsafe way.
// NOTE: Use this method with caution.
func ConvertStringToByteSlice(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
