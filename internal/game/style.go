package game

import (
	"strings"

	"github.com/ericktheredd5875/dicerealms/config"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Reverse   = "\033[7m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Purple    = "\033[35m"
	Cyan      = "\033[36m"
	Gray      = "\033[90m"
	White     = "\033[97m"
)

func Colorize(text string, color string) string {

	if !config.SupportsANSI {
		return text
	}

	var colorBuilder strings.Builder
	colorBuilder.WriteString(color)
	colorBuilder.WriteString(text)
	colorBuilder.WriteString(Reset)

	return colorBuilder.String()
}

func ColorizeError(text string) string {
	return Colorize("!! "+text, Red+Underline)
}

func ColorizeSuccess(text string) string {
	return Colorize(text, Green)
}

func ColorizeWarning(text string) string {
	return Colorize(text, Yellow)
}
