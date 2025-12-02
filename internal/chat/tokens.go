package chat

import (
	"strings"

	"github.com/teasec4/ollama-go-cli/internal/constants"
)

// EstimateTokens estimates token count using simple heuristics
// Approximation: ~1 token per 4 characters or ~1.3 tokens per word
func EstimateTokens(text string) int {
	if text == "" {
		return 0
	}

	// Method 1: Count by characters (1 token ≈ 4 chars)
	charCount := len(text)
	charTokens := (charCount + constants.CharactersPerToken - 1) / constants.CharactersPerToken // Round up

	// Method 2: Count by words (1 token ≈ 0.75 words, so 1 word ≈ 1.3 tokens)
	words := strings.Fields(text)
	wordTokens := (len(words) * constants.TokensPerWordFactor) / constants.TokensPerWordDivisor

	// Average both methods
	estimated := (charTokens + wordTokens) / 2

	// Minimum 1 token if there's any text
	if estimated < 1 {
		estimated = 1
	}

	return estimated
}

// CountSessionTokens counts tokens for entire session
func (s *Session) CountSessionTokens() int {
	total := 0
	for _, msg := range s.Messages {
		total += EstimateTokens(msg.Content)
	}
	return total
}
