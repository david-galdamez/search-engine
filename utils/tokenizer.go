package utils

import (
	"iter"
	"strings"
	"unicode"
)

func Tokenizer(text string) iter.Seq[string] {
	//trims punctuations and split by spaces
	cleanText := strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}, text)

	wordsIterator := strings.FieldsSeq(strings.ToLower(cleanText))

	return wordsIterator
}
