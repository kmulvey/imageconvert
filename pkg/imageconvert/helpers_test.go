package imageconvert

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	cp "github.com/otiai10/copy"

	"github.com/kmulvey/goutils"
	"github.com/stretchr/testify/assert"
)

type testPair struct {
	Name string
	Type string
}

func makeTestDir(t *testing.T) string {
	var testdir = "testdir_" + goutils.RandomString(5)
	assert.NoError(t, os.MkdirAll(testdir, os.ModePerm))

	assert.NoError(t, cp.Copy("testimages", testdir))

	return testdir
}

func moveImage(t *testing.T, testdir string, image testPair) string {

	// copy test image into this dir
	from, err := os.Open(image.Name)
	assert.NoError(t, err)

	var testImage = filepath.Join(testdir, path.Base(image.Name))
	to, err := os.Create(testImage)
	assert.NoError(t, err)

	_, err = io.Copy(to, from)
	assert.NoError(t, err)

	// close from and to before converting
	err = from.Close()
	assert.NoError(t, err)
	err = to.Close()
	assert.NoError(t, err)

	return testImage
}
