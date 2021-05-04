package main

import "os"

const (
	ansiReset  = "\033[0m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiCyan   = "\033[36m"
	ansiGray   = "\033[90m"
	ansiBold   = "\033[1m"
)

func colorEnabled(noColor bool) bool {
	if noColor {
		return false
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func colorize(enabled bool, code, s string) string {
	if !enabled {
		return s
	}
	return code + s + ansiReset
}

func priorityColor(p Priority) string {
	switch p {
	case PriorityHigh:
		return ansiRed
	case PriorityLow:
		return ansiGreen
	default:
		return ansiYellow
	}
}
