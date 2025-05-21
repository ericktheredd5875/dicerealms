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
