package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStringToByteSlice(t *testing.T) {
	testcases := []struct {
		s string
		b []byte
	}{
		{
			s: "",
			b: []byte(nil),
		},
		{
			s: "hello world",
			b: []byte("hello world"),
		},
	}

	for i, testcase := range testcases {
		t.Run(fmt.Sprintf("Running test case %d", i), func(t *testing.T) {
			result := ConvertStringToByteSlice(testcase.s)
			assert.Equal(t, testcase.b, result)
		})
	}
}

func TestNewStringPointer(t *testing.T) {
	assert.NotNil(t, NewStringPointer("hello"))
	assert.NotNil(t, NewStringPointer(""))
}
