package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// export ASSERT_ROOT_DIR=C:/CodeBases/dicerealms
func TestFindProjectRoot(t *testing.T) {

	rootDir := os.Getenv("ASSERT_ROOT_DIR")
	if rootDir == "" {
		rootDir = "/home/runner/work/dicerealms/dicerealms"
	}

	root, err := FindProjectRoot("")
	assert.NoError(t, err)
	assert.Equal(t, rootDir, root)
}
