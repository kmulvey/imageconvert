package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResizeValue(t *testing.T) {
	t.Parallel()

	assert.Equal(t, uint16(5120), getResizeValue("5120"))
}
