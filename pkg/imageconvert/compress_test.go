package imageconvert

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/kmulvey/imageconvert/v2/testimages"
	"github.com/stretchr/testify/assert"
)

func TestQualityCheck(t *testing.T) {
	t.Parallel()

	var testdir = testimages.MakeTestDir(t)
	var testImage = filepath.Join(testdir, "realjpg.jpg")

	aboveThreshold, err := QualityCheck(90, testImage)
	assert.NoError(t, err)
	assert.True(t, aboveThreshold)

	assert.NoError(t, os.WriteFile(filepath.Join(testdir, "test.txt"), make([]byte, 10), 0600))
	aboveThreshold, err = QualityCheck(90, filepath.Join(testdir, "test.txt"))
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
	var testImage = filepath.Join(testdir, "realjpg.jpg")

	var compressed, _, err = CompressJPEG(90, testImage)
	assert.NoError(t, err)
	assert.True(t, compressed)

	// do it til it wont compress anymore
	var skipped bool
	for range 10 {
		compressed, _, err = CompressJPEG(90, testImage)
		assert.NoError(t, err)
		if !compressed {
			skipped = true
			break
		}
	}
	if runtime.GOOS != "windows" { // windows cant do anything right
		assert.True(t, skipped)
	}

	assert.NoError(t, os.WriteFile(filepath.Join(testdir, "test.txt"), make([]byte, 10), 0600))
	compressed, _, err = CompressJPEG(90, filepath.Join(testdir, "test.txt"))
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
