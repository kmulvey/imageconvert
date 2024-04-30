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

	resized, err := Resize(testImage, 300, 200, 200, 100)
	assert.NoError(t, err)
	assert.True(t, resized)

	resized, err = Resize("noexist", 300, 200, 200, 100)
	assert.Error(t, err)
	assert.False(t, resized)

	resized, err = Resize("./compress.go", 300, 200, 200, 100)
	assert.Error(t, err)
	assert.False(t, resized)

	resized, err = Resize(testImage, 3000, 2000, 2000, 1000)
	assert.NoError(t, err)
	assert.False(t, resized)

	assert.NoError(t, os.RemoveAll(testdir))
}
