package bot

import (
	"regexp"
	"strings"
)

var anySeparatorRegexp = regexp.MustCompile(`[^а-яА-ЯёЁ\w\d]+`)

// tokenize splits input by any non-alphabetic/non-numeric characters and lowercases tokens.
func tokenize(message string) []string {
	splitted := anySeparatorRegexp.Split(message, -1)
	result := lowerCase(splitted)
	return result
}

func lowerCase(tokens []string) []string {
	result := make([]string, 0, len(tokens))
	for _, token := range tokens {
		result = append(result, strings.ToLower(token))
	}
	return result
}
