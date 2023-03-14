package imageconvert

import (
	"os"
	"strings"
	"testing"

	"github.com/kmulvey/path"
	"github.com/stretchr/testify/assert"
)

func TestConvertImage(t *testing.T) {
	t.Parallel()

	// setup
	var testdir = makeTestDir(t)
	var testImage = moveImage(t, testdir, testPair{Name: "./testimages/test.webp", Type: "jpeg"})
	var ic, err = NewWithDefaults(testImage, "", 0)
	ic.WithCompression()
	ic.WithResize(200, 100, 300, 200)
	assert.NoError(t, err)

	originalFile, err := path.NewEntry(testImage, 0)
	assert.NoError(t, err)

	var cr = ic.convertImage(originalFile)
	assert.NoError(t, cr.Error)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(cr.OriginalFileName, "test.webp"))
	assert.True(t, strings.HasSuffix(cr.ConvertedFileName, "test.jpg"))
	assert.Equal(t, "webp", cr.ImageType)
	assert.True(t, cr.Compressed)
	assert.False(t, cr.Renamed)
	assert.True(t, cr.Resized)

	assert.NoError(t, os.RemoveAll(testdir))
}
