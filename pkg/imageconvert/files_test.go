package imageconvert

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSkipMap(t *testing.T) {
	t.Parallel()

	// setup
	var testdir = makeTestDir(t)
	var handle, err = os.OpenFile(filepath.Join(testdir, "skipFile"), os.O_RDWR|os.O_CREATE, 0755)
	assert.NoError(t, err)
	_, err = handle.WriteString("realjpg.jpg")
	assert.NoError(t, err)
	err = handle.Close()
	assert.NoError(t, err)

	ic, err := NewWithDefaults(testdir, filepath.Join(testdir, "skipFile"), 0)
	assert.NoError(t, err)

	skipMap, err := ic.ParseSkipMap()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(skipMap))

	assert.NoError(t, os.RemoveAll(testdir))
}

func TestHasEOI(t *testing.T) {
	t.Parallel()

	// setup
	var testdir = makeTestDir(t)
	var testImage = moveImage(t, testdir, testPair{Name: "./testimages/realjpg.jpg", Type: "jpeg"})
	assert.True(t, hasEOI(testImage))
	assert.False(t, hasEOI("./compress.go"))
	assert.NoError(t, os.RemoveAll(testdir))
}

func TestEscapeFilePath(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "/some/file/\\&name", EscapeFilePath("/some/file/&name"))
	assert.Equal(t, "/some/file/\\(name\\)", EscapeFilePath("/some/file/(name)"))
}
