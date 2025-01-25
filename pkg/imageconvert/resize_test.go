package imageconvert

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kmulvey/imageconvert/v2/testimages"
	"github.com/stretchr/testify/assert"
)

func TestResize(t *testing.T) {
	t.Parallel()

	var testdir = testimages.MakeTestDir(t)
	var testImage = filepath.Join(testdir, "realjpg.jpg")

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

	testImage = filepath.Join(testdir, "realjpg-portrait.jpg")
	resized, err = Resize(testImage, 300, 200, 200, 100)
	assert.NoError(t, err)
	assert.True(t, resized)

	assert.NoError(t, os.RemoveAll(testdir))
}
