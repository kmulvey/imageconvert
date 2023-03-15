package imageconvert

import (
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/kmulvey/goutils"
	"github.com/stretchr/testify/assert"
)

type testPair struct {
	Name string
	Type string
}

func makeTestDir(t *testing.T) string {
	var testdir = "testdir_" + goutils.RandomString(5)

	var cmd = exec.Command("cp", "-r", "testimages/", testdir)
	var err = cmd.Run()
	assert.NoError(t, err)

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
