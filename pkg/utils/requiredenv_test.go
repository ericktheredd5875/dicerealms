package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiredEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	RequiredEnv("TEST_KEY")
}

func TestObtainEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	os.Setenv("TEST_KEY_2", "")

	assert.Equal(t, "test_value", ObtainEnv("TEST_KEY", "default_value"))
	assert.Equal(t, "default_value", ObtainEnv("TEST_KEY_2", "default_value"))
}
