package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindProjectRoot(t *testing.T) {

	rootDir := os.Getenv("ASSERT_ROOT_DIR")

	root, err := FindProjectRoot("")
	assert.NoError(t, err)
	assert.Equal(t, rootDir, root)
}
