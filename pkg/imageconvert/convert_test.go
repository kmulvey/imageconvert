package imageconvert

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testImage struct {
	inputPath      string
	outputPath     string
	imageType      string
	shouldConvert  bool
	wouldOverwrite bool
	err            bool
}

var testImages = []testImage{
	{inputPath: "test.png", outputPath: "test.jpg", imageType: "png", shouldConvert: true, wouldOverwrite: false, err: false},
	{inputPath: "fakejpg.jpg", outputPath: "fakejpg.jpg", imageType: "png", shouldConvert: true, wouldOverwrite: true, err: false},
	{inputPath: "realjpg-portrait.jpg", outputPath: "realjpg-portrait.jpg", imageType: "jpeg", shouldConvert: false, wouldOverwrite: true, err: false},
	{inputPath: "testwebp.webp", outputPath: "testwebp.jpg", imageType: "webp", shouldConvert: true, wouldOverwrite: false, err: false},
	{inputPath: "realjpg.jpg", outputPath: "realjpg.jpg", imageType: "jpeg", shouldConvert: false, wouldOverwrite: true, err: false},
	{inputPath: "realjpg.png", outputPath: "realjpg.png", imageType: "png", shouldConvert: true, wouldOverwrite: true, err: true},
}

func TestConvert(t *testing.T) {
	t.Parallel()

	var testdir = makeTestDir(t)

	for _, image := range testImages {

		var testImage = filepath.Join(testdir, image.inputPath)

		// convert
		convertedImage, format, err := Convert(testImage)
		assert.Equal(t, image.err, err != nil, image.inputPath)
		assert.Equal(t, image.outputPath, filepath.Base(convertedImage), image.inputPath)
		assert.Equal(t, image.imageType, format)

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

func TestConvertJpeg(t *testing.T) {
	t.Parallel()

	var testdir = makeTestDir(t)

	var image = testPair{"testimages/realjpg.jpg", "jpeg"}

	var testImage = moveImage(t, testdir, image)

	// convert
	convertedImage, format, err := Convert(testImage)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(testdir, "realjpg.jpg"), convertedImage)
	assert.Equal(t, image.Type, format)

	// make sure input file was deleted by Convert()
	if _, err := os.Stat(testImage); !errors.Is(err, os.ErrNotExist) {
		assert.NoError(t, err)
	}

	// make sure the converted file really exists
	_, err = os.Stat(convertedImage)
	assert.NoError(t, err)

	// clean up
	err = os.Remove(convertedImage)
	assert.NoError(t, err)

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

	var testdir = makeTestDir(t)

	for _, image := range testImages {
		var overwrite = WouldOverwrite(filepath.Join(testdir, image.inputPath))
		assert.Equal(t, image.wouldOverwrite, overwrite, image.inputPath)
	}

	assert.NoError(t, os.RemoveAll(testdir))
}
