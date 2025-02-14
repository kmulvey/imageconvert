package imageconvert

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kmulvey/imageconvert/v2/testimages"
	"github.com/stretchr/testify/assert"
)

type compressTestCase struct {
	testimages.TestCase
	ShouldCompress   bool
	PartialErrString string
}

var compressTestCases []compressTestCase // nolint: gochecknoglobals

func TestQualityCheck(t *testing.T) {
	t.Parallel()

	var testdir = testimages.MakeTestDir(t)
	for _, testCase := range compressTestCases {
		var testImage = filepath.Join(testdir, testCase.InputPath)

		aboveThreshold, err := QualityCheck(90, testImage)
		assert.NoError(t, err, testCase.InputPath)
		assert.True(t, aboveThreshold, testCase.InputPath)
	}

	assert.NoError(t, os.WriteFile(filepath.Join(testdir, "test.txt"), make([]byte, 10), 0600))
	aboveThreshold, err := QualityCheck(90, filepath.Join(testdir, "test.txt"))
	assert.True(t, strings.HasPrefix(err.Error(), "error running identify on image:"))
	assert.False(t, aboveThreshold)

	aboveThreshold, err = QualityCheck(90, "not a file")
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.False(t, aboveThreshold)

	assert.NoError(t, os.RemoveAll(testdir))
}

func TestCompressJPEG(t *testing.T) {
	t.Parallel()

	var testdir = testimages.MakeTestDir(t)
	for _, testCase := range compressTestCases {
		var testImage = filepath.Join(testdir, testCase.InputPath)

		var compressed, _, err = CompressJPEG(90, testImage)
		if testCase.PartialErrString != "" {
			assert.Error(t, err, testCase.InputPath)
			assert.Contains(t, err.Error(), testCase.PartialErrString, testCase.InputPath)
		} else {
			assert.NoError(t, err, testCase.InputPath)
		}
		assert.Equal(t, testCase.ShouldCompress, compressed, testCase.InputPath)
	}

	assert.NoError(t, os.WriteFile(filepath.Join(testdir, "test.txt"), make([]byte, 10), 0600))
	compressed, _, err := CompressJPEG(90, filepath.Join(testdir, "test.txt"))
	assert.True(t, strings.HasPrefix(err.Error(), "error running jpegoptim on image:"))
	assert.False(t, compressed)

	compressed, _, err = CompressJPEG(90, "not a file")
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.False(t, compressed)

	assert.NoError(t, os.RemoveAll(testdir))
}

func TestEscapeFilePath(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "/some/file/\\&name", EscapeFilePath("/some/file/&name"))
	assert.Equal(t, "/some/file/\\(name\\)", EscapeFilePath("/some/file/(name)"))
}
