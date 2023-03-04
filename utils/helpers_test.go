package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSliceHas(t *testing.T) {
	slice := []string{"foo1", "foo2", "foo3"}

	element := "foo2"
	has := IsSliceHas(slice, element)
	assert.Equal(t, has, true)

	element = "foo4"
	notHas := IsSliceHas(slice, element)
	assert.Equal(t, notHas, false)
}

func TestIsFileHasAnyExt_ShouldHas(t *testing.T) {
	filePath := "/home/somedir/test/test.txt"
	has := IsExtHasAnyValue(filePath)
	assert.Equal(t, has, true)

	filePath = "/home/user/somedir/test"
	notHas := IsExtHasAnyValue(filePath)
	assert.Equal(t, notHas, false)
}
