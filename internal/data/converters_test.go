package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertBoolToInt(t *testing.T) {
	assert.Equal(t, 0, ConvertBoolToInt(false))
	assert.Equal(t, 1, ConvertBoolToInt(true))
}
