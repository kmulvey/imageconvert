package imageconvert

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWithDefaults(t *testing.T) {
	t.Parallel()

	// setup
	var testdir = makeTestDir(t)
	var ic, err = NewWithDefaults("noexist", "", 0)
	assert.Error(t, err)
	assert.Equal(t, uint8(1), ic.Threads)

	assert.NoError(t, os.RemoveAll(testdir))
}
