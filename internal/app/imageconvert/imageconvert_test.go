package imageconvert

import (
	"os"
	"testing"
	"time"

	"github.com/kmulvey/path"
	"github.com/stretchr/testify/assert"
)

func TestNewImageConverterGood(t *testing.T) {
	t.Parallel()

	var testdir = makeTestDir(t)
	var testImage = moveImage(t, testdir, testPair{Name: "./testimages/realjpg.jpg", Type: "jpeg"})
	var testImageEntry, err = path.NewEntry(testImage, 0)
	assert.NoError(t, err)

	var config = &ImageConverterConfig{
		OriginalImages: []string{testImage},
		Threads:        1,
		Quality:        30,
	}
	ic, err := NewImageConverter(config)
	assert.NoError(t, err)
	assert.Equal(t, []path.Entry{testImageEntry}, ic.OriginalImagesEntries)
	assert.Equal(t, "", ic.WatchDir)
	assert.Equal(t, 1, ic.Threads)
	assert.Equal(t, "processed.log", ic.SkipMapFileHandle.Name())
	assert.Equal(t, 0, len(ic.SkipMap))
	assert.Equal(t, 30, ic.Quality)
	assert.False(t, ic.Force)
	assert.False(t, ic.DeleteOriginal)
	assert.Equal(t, time.Time{}, ic.TimeRange.From)
	assert.Equal(t, time.Time{}, ic.TimeRange.To)
	assert.Equal(t, uint16(0), ic.ResizeWidth)
	assert.Equal(t, uint16(0), ic.ResizeHeight)
	assert.Equal(t, uint16(0), ic.ResizeWidthThreshold)
	assert.Equal(t, uint16(0), ic.ResizeHeightThreshold)

	assert.NoError(t, ic.Shutdown())
	assert.NoError(t, os.RemoveAll("processed.log"))
	assert.NoError(t, os.RemoveAll(testdir))
}
