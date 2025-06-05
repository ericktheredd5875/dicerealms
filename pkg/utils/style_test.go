package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ericktheredd5875/dicerealms/config"
)

func TestColorize_ANSIEnabled(t *testing.T) {
	config.SupportsANSI = true
	result := Colorize("Hello", Red)
	expected := Red + "Hello" + Reset
	assert.Equal(t, expected, result)
}

func TestColorize_ANSIDisabled(t *testing.T) {
	config.SupportsANSI = false
	result := Colorize("Hello", Green)
	assert.Equal(t, "Hello", result)
}

func TestColorizeError(t *testing.T) {
	config.SupportsANSI = true
	result := ColorizeError("Something went wrong")
	expected := Red + Underline + "!! Something went wrong" + Reset
	assert.Equal(t, expected, result)
}

func TestColorizeSuccess(t *testing.T) {
	config.SupportsANSI = true
	result := ColorizeSuccess("All good")
	expected := Green + "All good" + Reset
	assert.Equal(t, expected, result)
}

func TestColorizeWarning(t *testing.T) {
	config.SupportsANSI = true
	result := ColorizeWarning("Be careful")
	expected := Yellow + "Be careful" + Reset
	assert.Equal(t, expected, result)
}

func TestColorizeInfo(t *testing.T) {
	config.SupportsANSI = true
	result := ColorizeInfo("FYI")
	expected := Blue + "FYI" + Reset
	assert.Equal(t, expected, result)
}
