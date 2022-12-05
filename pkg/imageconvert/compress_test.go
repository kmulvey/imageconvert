package imageconvert

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kmulvey/goutils"
	"github.com/stretchr/testify/assert"
)

func TestQualityCheck(t *testing.T) {
	t.Parallel()

	var testdir = goutils.RandomString(5)
	var err = os.Mkdir(testdir, os.ModePerm)
	assert.NoError(t, err)

	var testImage = moveImage(t, testdir, testPair{Name: "./testimages/realjpg.jpg", Type: "jpeg"})
	aboveThreshold, err := QualityCheck(90, testImage)
	assert.NoError(t, err)
	assert.True(t, aboveThreshold)

	assert.NoError(t, os.WriteFile(filepath.Join(testdir, "test.txt"), make([]byte, 10), os.ModePerm))
	aboveThreshold, err = QualityCheck(90, filepath.Join(testdir, "test.txt"))
	assert.True(t, strings.HasPrefix(err.Error(), "error running identify on image:"))
	assert.False(t, aboveThreshold)
}
