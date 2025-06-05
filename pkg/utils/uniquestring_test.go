package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUniqueString(t *testing.T) {
	result, err := GenerateUniqueString(16)
	assert.NoError(t, err)
	assert.Len(t, result, 32) // 16 bytes = 32 hex characters
}

func TestGenerateUniqueString_InvalidLength(t *testing.T) {
	result, err := GenerateUniqueString(0)
	assert.NoError(t, err)
	assert.Equal(t, "", result) // hex.EncodeToString on 0 bytes is ""
}

func TestGenerateUniqueID(t *testing.T) {
	result, err := GenerateUniqueID()
	assert.NoError(t, err)
	assert.Len(t, result, 32)
}

func TestGenerateKHash_WithInput(t *testing.T) {
	result, err := GenerateKHash("testinput", "")
	assert.NoError(t, err)
	assert.Len(t, result, 8)
}

func TestGenerateKHash_WithoutInput(t *testing.T) {
	result, err := GenerateKHash("", "")
	assert.NoError(t, err)
	assert.Len(t, result, 8)
}

func TestGenerateKHash_DateType(t *testing.T) {
	result, err := GenerateKHash("testinput", "date")
	assert.NoError(t, err)
	assert.Len(t, result, 8)
}
