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

	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"

	BrRed    = "\033[91m"
	BrGreen  = "\033[92m"
	BrYellow = "\033[93m"
	BrBlue   = "\033[94m"
	BrPurple = "\033[95m"
	BrCyan   = "\033[96m"
	White    = "\033[97m"

	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
	BgPurple = "\033[45m"
	BgCyan   = "\033[46m"

	BgLtRed    = "\033[101m"
	BgLtGreen  = "\033[102m"
	BgLtYellow = "\033[103m"
	BgLtBlue   = "\033[104m"
	BgLtPurple = "\033[105m"
	BgLtCyan   = "\033[106m"
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

func ColorizeInfo(text string) string {
	return Colorize(text, Blue)
}
