package imageconvert

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kmulvey/humantime"
	"github.com/kmulvey/imageconvert/v2/testimages"
	"github.com/stretchr/testify/assert"
)

func TestNewWithDefaults(t *testing.T) {
	t.Parallel()

	// setup
	var testdir, err = testimages.MakeTestDir()
	assert.NoError(t, err)
	var testImage = filepath.Join(testdir, "realjpg.jpg")

	ic, err := New(testImage, "", 0, WithCompression(uint8(90)), WithWatch(), WithForce(), WithResize(200, 100, 300, 200), WithThreads(3), WithTimeRange(humantime.TimeRange{From: time.Time{}, To: time.Now()}))
	assert.NoError(t, err)

	assert.Equal(t, uint8(90), ic.CompressQuality)
	assert.True(t, ic.Watch)
	assert.True(t, ic.Force)
	assert.Equal(t, uint16(200), ic.ResizeWidth)
	assert.Equal(t, uint16(300), ic.ResizeWidthThreshold)
	assert.Equal(t, uint16(100), ic.ResizeHeight)
	assert.Equal(t, uint16(200), ic.ResizeHeightThreshold)
	assert.Equal(t, 3, ic.Threads)
	assert.Equal(t, time.Time{}, ic.TimeRange.From)
	assert.Equal(t, time.Now().Day(), ic.TimeRange.To.Day())

	assert.NoError(t, os.RemoveAll(testdir))
}

func TestStartSlice(t *testing.T) {
	t.Parallel()

	// setup
	var testdir, err = testimages.MakeTestDir()
	assert.NoError(t, err)

	ic, err := New(testdir, "", 1, WithCompression(uint8(90)))
	assert.NoError(t, err)
	assert.Equal(t, 1, ic.Threads)

	compressedTotal, renamedTotal, resizedTotal, totalFiles, conversionTypeTotals, err := ic.Start(nil)
	assert.NoError(t, err)

	assert.Equal(t, 5, compressedTotal)
	assert.Equal(t, 0, renamedTotal)
	assert.Equal(t, 0, resizedTotal)
	assert.Equal(t, 6, totalFiles)
	assert.EqualValues(t, map[string]int{"jpeg": 2, "png": 2, "webp": 1}, conversionTypeTotals)

	assert.NoError(t, os.RemoveAll("processed.log"))
	assert.NoError(t, os.RemoveAll(testdir))
}
