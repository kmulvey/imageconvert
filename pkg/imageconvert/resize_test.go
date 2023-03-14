package imageconvert

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResize(t *testing.T) {
	t.Parallel()

	// setup
	var testdir = makeTestDir(t)
	var testImage = moveImage(t, testdir, testPair{Name: "./testimages/realjpg.jpg", Type: "jpeg"})
	var ic, err = NewWithDefaults(testImage, "", 0)
	assert.NoError(t, err)
	ic.WithResize(200, 100, 300, 200)

	resized, err := ic.Resize(testImage)
	assert.NoError(t, err)
	assert.True(t, resized)

	resized, err = ic.Resize("noexist")
	assert.Error(t, err)
	assert.False(t, resized)

	resized, err = ic.Resize("./compress.go")
	assert.Error(t, err)
	assert.False(t, resized)

	ic.WithResize(2000, 1000, 3000, 2000)
	resized, err = ic.Resize(testImage)
	assert.NoError(t, err)
	assert.False(t, resized)

	assert.NoError(t, os.RemoveAll(testdir))
}
