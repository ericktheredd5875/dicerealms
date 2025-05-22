package mcp

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	line := `#$#mcp-emote: text="waves" mood="happy"`
	msg, err := Parse(line)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := &Message{
		Tag: "mcp-emote",
		Args: map[string]string{
			"text": "waves",
			"mood": "happy",
		},
	}

	if !reflect.DeepEqual(msg, expected) {
		t.Errorf("Parsed message did not match:\nGot: %#v\nExpected: %#v", msg, expected)
	}
}

func TestParse_Invalid(t *testing.T) {
	cases := []string{
		"",
		"hello",
		"#$#badformat",
		"#$#incomplete:",
	}

	for _, line := range cases {
		msg, err := Parse(line)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if msg != nil {
			t.Errorf("Expected nil message for line: %q", line)
		}
	}
}
