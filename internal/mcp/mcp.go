package mcp

import (
	"strings"
)

type Message struct {
	Tag  string
	Args map[string]string
}

func Parse(line string) (*Message, error) {
	if !strings.HasPrefix(line, "#$#") {
		return nil, nil // Not a valid MCP message
	}

	parts := strings.SplitN(line[3:], ":", 2)
	if len(parts) != 2 {
		return nil, nil
	}

	tag := strings.TrimSpace(parts[0])
	args := parseArgs(parts[1])

	return &Message{Tag: tag, Args: args}, nil
}

func parseArgs(input string) map[string]string {

	args := make(map[string]string)
	tokens := strings.Fields(input)

	for _, token := range tokens {
		parts := strings.SplitN(token, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := strings.Trim(parts[1], `"`)
			args[key] = value
		}
	}

	return args
}
