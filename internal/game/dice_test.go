package game

import (
	"strings"
	"testing"
)

func TestRollValid(t *testing.T) {
	_, detail, err := Roll("2d6+3")
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if !strings.Contains(detail, "Rolled 2d6+3") {
		t.Errorf("Unexpected detail: %s", detail)
	}
}

func TestRollInvalid(t *testing.T) {
	_, _, err := Roll("foo")
	if err == nil {
		t.Error("Expected error for invalid dice expression")
	}
}
