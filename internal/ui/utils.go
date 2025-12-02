package ui

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

// visibleLen returns display width excluding ANSI codes (accounts for wide characters)
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
	return runewidth.StringWidth(visible)
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

// renderContextBar generates a visual context usage bar
func renderContextBar(used, maxSize int, barLength int) string {
	if maxSize <= 0 {
		return ""
	}

	percentage := float64(used) / float64(maxSize)
	if percentage > 1.0 {
		percentage = 1.0
	}

	filledLength := int(float64(barLength) * percentage)
	emptyLength := barLength - filledLength

	filled := strings.Repeat("█", filledLength)
	empty := strings.Repeat("░", emptyLength)

	percentStr := fmt.Sprintf("%.0f%%", percentage*100)
	bar := fmt.Sprintf("[%s%s] %s", filled, empty, percentStr)

	return bar
}

// estimateTokens estimates token count for text
// Rough approximation: ~1 token per 4 characters or ~1.3 tokens per word
func estimateTokens(text string) int {
	words := len(strings.Fields(text))
	if words > 0 {
		return int(float64(words) * 1.3)
	}
	return len(text) / 4
}


