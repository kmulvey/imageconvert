package imageconvert

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kmulvey/imageconvert/v2/testimages"
	"github.com/kmulvey/path"
	"github.com/stretchr/testify/assert"
)

func TestConvertImage(t *testing.T) {
	t.Parallel()

	// setup
	var testdir = testimages.MakeTestDir(t)
	var testImage = filepath.Join(testdir, "testwebp.webp")

	var ic, err = NewWithDefaults(testImage, "", 0)
	ic.WithCompression(uint8(90))
	ic.WithResize(200, 100, 300, 200)
	assert.NoError(t, err)

	originalFile, err := path.NewEntry(testImage, 0)
	assert.NoError(t, err)

	var cr = ic.convertImage(originalFile)
	assert.NoError(t, cr.Error)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(cr.OriginalFileName, "testwebp.webp"))
	assert.True(t, strings.HasSuffix(cr.ConvertedFileName, "testwebp.jpg"))
	assert.Equal(t, "webp", cr.ImageType)
	assert.True(t, cr.Compressed)
	assert.False(t, cr.Renamed)
	assert.True(t, cr.Resized)

	// copy test image into this dir
	// we cant use moveImage() because we need to change the extension
	from, err := os.Open(filepath.Join(testdir, "realjpg.jpg"))
	assert.NoError(t, err)

	testImage = filepath.Join(testdir, "realjpg.jpeg")
	to, err := os.Create(testImage)
	assert.NoError(t, err)

	_, err = io.Copy(to, from)
	assert.NoError(t, err)

	// close from and to before converting
	err = from.Close()
	assert.NoError(t, err)
	err = to.Close()
	assert.NoError(t, err)

	ic.WithCompression(uint8(90))
	ic.WithResize(200, 100, 300, 200)
	assert.NoError(t, err)

	originalFile, err = path.NewEntry(testImage, 0)
	assert.NoError(t, err)

	cr = ic.convertImage(originalFile)
	assert.NoError(t, cr.Error)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(cr.OriginalFileName, "realjpg.jpeg"))
	assert.True(t, strings.HasSuffix(cr.ConvertedFileName, "realjpg.jpg"))
	assert.Equal(t, "jpeg", cr.ImageType)
	assert.True(t, cr.Compressed)
	assert.True(t, cr.Renamed)
	assert.True(t, cr.Resized)

	assert.NoError(t, os.RemoveAll(testdir))
}
