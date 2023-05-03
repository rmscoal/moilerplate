package utils

// Thanks to:
// https://stackoverflow.com/questions/46128016/insert-a-value-in-a-slice-at-a-given-index
func InsertAt[T any](s []T, index int, val T) []T {
	n := len(s)
	if index < 0 {
		index = index%n + n
	}

	switch {
	case index == n:
		return append(s, val)
	case index < n:
		s = append(s[:index+1], s[index:]...)
		s[index] = val
		return s
	case index < cap(s):
		s = s[:index+1]
		var zero T
		for i := n; i < index; i++ {
			s[i] = zero
		}
		s[index] = val
		return s
	default: // index > cap(s)
		t := make([]T, index+1)
		if n > 0 {
			copy(t, s)
		}
		t[index] = val
		return t
	}
}
