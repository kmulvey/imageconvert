package imageconvert

import (
	"errors"
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testPair struct {
	Name string
	Type string
}

func TestCompress(t *testing.T) {
	for _, image := range []testPair{{"testimages/test.png", "png"}, {"testimages/test.webp", "webp"}} {

		// copy test image into this dir
		var from, err = os.Open(image.Name)
		assert.NoError(t, err)

		var testImage = path.Base(image.Name)
		to, err := os.Create(testImage)
		assert.NoError(t, err)

		_, err = io.Copy(to, from)
		assert.NoError(t, err)

		// close from and to before converting
		err = from.Close()
		assert.NoError(t, err)
		err = to.Close()
		assert.NoError(t, err)

		// convert
		convertedImage, format, err := Convert(testImage)
		assert.NoError(t, err)
		assert.Equal(t, "test.jpg", convertedImage)
		assert.Equal(t, image.Type, format)

		// make sure input file was deleted by Convert()
		if _, err := os.Stat(testImage); !errors.Is(err, os.ErrNotExist) {
			assert.NoError(t, err)
		}

		// clean up
		err = os.Remove(convertedImage)
		assert.NoError(t, err)
	}
}
