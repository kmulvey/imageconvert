package imageconvert

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Parallel()

	var testdir = makeTestDir(t)
	var testImage = moveImage(t, testdir, testPair{Name: "./testimages/realjpg.jpg", Type: "jpeg"})
	var skipFile = filepath.Join(testdir, "processed.log")

	var config = &ImageConverterConfig{
		OriginalImages: []string{testImage},
		Threads:        1,
		Quality:        30,
		SkipMapFile:    skipFile,
	}
	ic, err := NewImageConverter(config)
	assert.NoError(t, err)

	processedTotal, resizedTotal, err := ic.Start()
	assert.NoError(t, err)
	assert.Equal(t, 1, processedTotal)
	assert.Equal(t, 0, resizedTotal)
	assert.NoError(t, ic.Shutdown())

	// verify
	processedLogContents, err := os.ReadFile(skipFile)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(processedLogContents), testImage))
	stat, err := os.Stat(filepath.Join(testdir, "realjpg.avif"))
	assert.NoError(t, err)
	assert.True(t, stat.Size() > 100)

	// do it again, should skip
	ic, err = NewImageConverter(config)
	assert.NoError(t, err)
	processedTotal, resizedTotal, err = ic.Start()
	assert.NoError(t, err)
	assert.Equal(t, 0, processedTotal)
	assert.Equal(t, 0, resizedTotal)
	assert.NoError(t, ic.Shutdown())

	assert.NoError(t, os.RemoveAll(testdir))
}
