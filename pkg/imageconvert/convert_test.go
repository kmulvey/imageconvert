package imageconvert

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kmulvey/imageconvert/v2/testimages"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	t.Parallel()

	var testdir = testimages.MakeTestDir(t)

	for _, image := range testimages.TestCases {

		var testImage = filepath.Join(testdir, image.InputPath)

		// convert
		convertedImage, format, err := Convert(testImage)
		assert.Equal(t, image.Err, err != nil, image.InputPath)
		assert.Equal(t, image.OutputPath, filepath.Base(convertedImage), image.InputPath)
		assert.Equal(t, image.ImageType, format)

		// make sure input file was deleted by Convert()
		if _, err := os.Stat(testImage); !errors.Is(err, os.ErrNotExist) {
			assert.NoError(t, err)
		}

		// make sure the converted file really exists
		_, err = os.Stat(convertedImage)
		assert.NoError(t, err)
	}
	assert.NoError(t, os.RemoveAll(testdir))
}

func TestConvertErrors(t *testing.T) {
	t.Parallel()

	var convertedImage, format, err = Convert("testImage")
	assert.Equal(t, "", convertedImage)
	assert.Equal(t, "", format)
	assert.True(t, strings.Contains(err.Error(), "error opening file for conversion, image: testImage, error: open testImage:"))

	assert.NoError(t, os.WriteFile("testImage", make([]byte, 100), 0600))
	convertedImage, format, err = Convert("testImage")
	assert.Equal(t, "", convertedImage)
	assert.Equal(t, "", format)
	assert.Equal(t, "error decoding image: testImage, error: image: unknown format", err.Error())
	assert.NoError(t, os.RemoveAll("testImage"))
}

func TestWouldOverwrite(t *testing.T) {
	t.Parallel()

	var testdir = testimages.MakeTestDir(t)

	for _, image := range testimages.TestCases {
		var overwrite = WouldOverwrite(filepath.Join(testdir, image.InputPath))
		assert.Equal(t, image.WouldOverwrite, overwrite, image.InputPath)
	}

	assert.NoError(t, os.RemoveAll(testdir))
}
