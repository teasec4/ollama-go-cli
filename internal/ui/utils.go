package ui

import (
	"strings"

	"github.com/teasec4/ollama-go-cli/internal/constants"
)

// visibleLen returns visible length excluding ANSI codes
func visibleLen(text string) int {
	visible := text
	for {
		start := strings.Index(visible, "\x1b[")
		if start == -1 {
			break
		}
		end := strings.Index(visible[start:], "m")
		if end == -1 {
			break
		}
		visible = visible[:start] + visible[start+end+1:]
	}
	return len(visible)
}

// wrapText wraps text to fit within max width
func wrapText(text string, maxWidth int) string {
	if maxWidth < 1 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if visibleLen(testLine) <= maxWidth {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// buildMemoryBar creates a visual memory usage bar
func buildMemoryBar(percent int) string {
	if percent > 100 {
		percent = 100
	}
	if percent < 0 {
		percent = 0
	}

	filledLength := (percent * constants.MemoryBarLength) / 100
	bar := "["
	for i := 0; i < constants.MemoryBarLength; i++ {
		if i < filledLength {
			bar += "="
		} else {
			bar += " "
		}
	}
	bar += "]"
	return bar
}
