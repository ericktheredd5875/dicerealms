package mcp

import (
	"regexp"
	"strings"
)

var argPattern = regexp.MustCompile(`(\w+)=(".*?"|\S+)`)

type Message struct {
	Tag  string
	Args map[string]string
}

func Parse(line string) (*Message, error) {
	if !strings.HasPrefix(line, "#$#mcp-") {
		return nil, nil // Not a valid MCP message
	}

	// Remove the #$# prefix
	line = line[3:]

	// Check if there's a colon
	if strings.Contains(line, ":") {
		parts := strings.SplitN(line, ":", 2)
		tag := strings.TrimSpace(parts[0])
		args := parseArgs(parts[1])
		return &Message{Tag: tag, Args: args}, nil
	} else {
		// No colon, treat the entire line as the tag
		tag := strings.TrimSpace(line)
		return &Message{Tag: tag, Args: make(map[string]string)}, nil
	}
}

func parseArgs(input string) map[string]string {

	args := make(map[string]string)
	matches := argPattern.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		key := match[1]
		value := match[2]
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			value = strings.Trim(value, `"`)
		}
		args[key] = value
	}

	return args
}
