package imageconvert

import (
	"errors"
	"fmt"
	"io"
	"os"
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
	var err = os.Mkdir(testdir, os.ModePerm)
	assert.NoError(t, err)

	return testdir
}

func TestConvertErrors(t *testing.T) {
	t.Parallel()

	var convertedImage, format, err = Convert("testImage")
	assert.Equal(t, "", convertedImage)
	assert.Equal(t, "", format)
	assert.Equal(t, "error opening file for conversion, image: testImage, error: open testImage: no such file or directory", err.Error())

	assert.NoError(t, os.WriteFile("testImage", make([]byte, 100), os.ModePerm))
	convertedImage, format, err = Convert("testImage")
	assert.Equal(t, "", convertedImage)
	assert.Equal(t, "", format)
	assert.Equal(t, "error decoding image: testImage, error: image: unknown format", err.Error())
	assert.NoError(t, os.RemoveAll("testImage"))
}

func TestConvert(t *testing.T) {
	t.Parallel()

	var testdir = makeTestDir(t)

	for _, image := range []testPair{{"testimages/test.png", "png"}, {"testimages/test.webp", "webp"}} {

		var testImage = moveImage(t, testdir, image)

		// convert
		convertedImage, format, err := Convert(testImage)
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(testdir, "test.jpg"), convertedImage)
		assert.Equal(t, image.Type, format)

		// make sure input file was deleted by Convert()
		if _, err := os.Stat(testImage); !errors.Is(err, os.ErrNotExist) {
			assert.NoError(t, err)
		}

		// make sure the converted file really exists
		_, err = os.Stat(convertedImage)
		assert.NoError(t, err)

		// clean up
		err = os.Remove(convertedImage)
		assert.NoError(t, err)
	}
	assert.NoError(t, os.RemoveAll(testdir))
}

func TestConvertJpeg(t *testing.T) {
	t.Parallel()

	var testdir = makeTestDir(t)

	var image = testPair{"testimages/realjpg.jpg", "jpeg"}

	var testImage = moveImage(t, testdir, image)

	// convert
	convertedImage, format, err := Convert(testImage)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(testdir, "realjpg.jpg"), convertedImage)
	assert.Equal(t, image.Type, format)

	// make sure input file was deleted by Convert()
	if _, err := os.Stat(testImage); !errors.Is(err, os.ErrNotExist) {
		assert.NoError(t, err)
	}

	// make sure the converted file really exists
	_, err = os.Stat(convertedImage)
	assert.NoError(t, err)

	// clean up
	err = os.Remove(convertedImage)
	assert.NoError(t, err)

	assert.NoError(t, os.RemoveAll(testdir))
}

func TestConvertWouldOverwrite(t *testing.T) {
	t.Parallel()

	var testdir = makeTestDir(t)

	var image = testPair{"testimages/test.png", "png"}

	var testImage = moveImage(t, testdir, image)

	// convert
	convertedImage, format, err := Convert(testImage)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(testdir, "test.jpg"), convertedImage)
	assert.Equal(t, image.Type, format)

	// do it again
	moveImage(t, testdir, image)
	convertedImage, format, err = Convert(testImage)
	assert.Equal(t, fmt.Sprintf("converting %s/test.png would overwrite an existing jpeg, skipping", testdir), err.Error())
	assert.Equal(t, filepath.Join(testdir, "test.png"), convertedImage)
	assert.Equal(t, "", format)

	// make sure input file was deleted by Convert()
	if _, err := os.Stat(testImage); !errors.Is(err, os.ErrNotExist) {
		assert.NoError(t, err)
	}

	// make sure the converted file really exists
	_, err = os.Stat(convertedImage)
	assert.NoError(t, err)

	// clean up
	err = os.Remove(convertedImage)
	assert.NoError(t, err)

	assert.NoError(t, os.RemoveAll(testdir))
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
