package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func Tokenizer(text string) []string {
	//trims punctuations and split by spaces
	normalizedText := removeAccents(text)
	re := regexp.MustCompile(`[^\p{L}\p{N}\s]+`)
	cleanText := re.ReplaceAllString(normalizedText, "")

	wordsIterator := strings.Fields(strings.ToLower(cleanText))

	return wordsIterator
}

func removeAccents(s string) string {
	t := norm.NFD.String(s)

	var sb strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}

		sb.WriteRune(r)
	}

	return sb.String()
}
