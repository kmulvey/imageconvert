package imageconvert

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/kmulvey/imageconvert/v2/testimages"
	"github.com/stretchr/testify/assert"
)

type convertTestCase struct {
	testimages.TestCase
	ShouldConvert    bool
	PartialErrString string
}

var convertTestCases []convertTestCase // nolint: gochecknoglobals

func TestConvert(t *testing.T) {
	t.Parallel()

	testdir, err := testimages.MakeTestDir()
	assert.NoError(t, err)

	for _, testCase := range testimages.TestCases {
		var testImage = filepath.Join(testdir, testCase.InputPath)

		// convert
		convertedImage, format, err := Convert(testImage)
		assert.Equal(t, testCase.Err, err != nil, testCase.InputPath)
		assert.Equal(t, testCase.OutputPath, filepath.Base(convertedImage), testCase.InputPath)
		assert.Equal(t, testCase.ImageType, format)

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
	assert.Contains(t, err.Error(), "error opening file for conversion, image: testImage, error: open testImage:")

	assert.NoError(t, os.WriteFile("testImage", make([]byte, 100), 0600))
	convertedImage, format, err = Convert("testImage")
	assert.Equal(t, "", convertedImage)
	assert.Equal(t, "", format)
	assert.Equal(t, "error decoding image: testImage, error: image: unknown format", err.Error())
	assert.NoError(t, os.RemoveAll("testImage"))
}

func TestWouldOverwrite(t *testing.T) {
	t.Parallel()

	testdir, err := testimages.MakeTestDir()
	assert.NoError(t, err)

	for _, image := range testimages.TestCases {
		var overwrite = WouldOverwrite(filepath.Join(testdir, image.InputPath))
		assert.Equal(t, image.WouldOverwrite, overwrite, image.InputPath)
	}

	assert.NoError(t, os.RemoveAll(testdir))
}
