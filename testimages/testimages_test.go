package testimages

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeTestDir(t *testing.T) {
	t.Parallel()

	dir, err := MakeTestDir()
	assert.NoError(t, err)

	assert.DirExists(t, dir)

	assert.NoError(t, os.RemoveAll(dir))
}
