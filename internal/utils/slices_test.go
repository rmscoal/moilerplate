package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testInsertAtCases = []struct {
	idx           int
	slice         []int
	expectedSlice []int
	val           int
	expectedLen   int
}{
	{
		idx:           10,
		val:           1,
		slice:         []int{1, 2, 3, 4, 5, 6, 7, 7, 8, 10, 9, 2, 12, 32, 83, 14, 34, 6, 1, 45, 6},
		expectedSlice: []int{1, 2, 3, 4, 5, 6, 7, 7, 8, 10, 1, 9, 2, 12, 32, 83, 14, 34, 6, 1, 45, 6},
		expectedLen:   21 + 1,
	},
	{
		idx:           5,
		val:           10,
		slice:         nil,
		expectedSlice: []int{0, 0, 0, 0, 0, 10},
		expectedLen:   5 + 1,
	},
	{
		idx: 29,
		val: 255,
		slice: []int{
			72, 60, 140, 222, 32, 133, 201, 4, 87, 122, 19, 134, 185, 62, 242, 66, 120, 191, 122, 92, 75, 140, 195, 29, 161, 196, 87, 158, 45, 89, 199, 54, 218, 178, 160, 15, 42, 170, 64, 2,
			40, 41, 87, 40, 120, 11, 212, 73, 88, 76, 61, 219, 81, 113, 183, 47, 186, 106, 77, 59, 20, 197, 52, 123, 228, 86, 215, 87, 128, 255, 145, 220, 129, 102, 170, 166, 238, 29, 179, 4,
			6, 142, 164, 52,
		},
		expectedSlice: []int{
			72, 60, 140, 222, 32, 133, 201, 4, 87, 122, 19, 134, 185, 62, 242, 66, 120, 191, 122, 92, 75, 140, 195, 29, 161, 196, 87, 158, 45, 255, 89, 199, 54, 218, 178, 160, 15, 42, 170, 64, 2,
			40, 41, 87, 40, 120, 11, 212, 73, 88, 76, 61, 219, 81, 113, 183, 47, 186, 106, 77, 59, 20, 197, 52, 123, 228, 86, 215, 87, 128, 255, 145, 220, 129, 102, 170, 166, 238, 29, 179, 4,
			6, 142, 164, 52,
		},
		expectedLen: 84 + 1,
	},
}

func TestInsertAt(t *testing.T) {
	for idx, test := range testInsertAtCases {
		t.Run(fmt.Sprintf("Test-Case-%d", idx+1), func(t *testing.T) {
			res := InsertAt(test.slice, test.idx, test.val)
			assert.Equal(t, test.expectedSlice, res)
			assert.Equal(t, res[test.idx], test.val)
			assert.Equal(t, test.expectedLen, len(res))
		})
	}
}

func BenchmarkInsertAt(b *testing.B) {
	for idx, test := range testInsertAtCases {
		b.Run(fmt.Sprintf("Test-Case-%d", idx+1), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				InsertAt(test.slice, test.idx, test.val)
			}
		})
	}
}
