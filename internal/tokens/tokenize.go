package tokens

import (
	"regexp"
	"strings"
)

var anySeparatorRegexp = regexp.MustCompile(`[^а-яА-ЯёЁ\w\d]+`)

// Tokenize splits input by any non-alphabetic/non-numeric characters and lowercases tokens.
func Tokenize(message string) []string {
	splitted := anySeparatorRegexp.Split(message, -1)
	result := tokensToLowerCase(splitted)
	return result
}

func tokensToLowerCase(tokens []string) []string {
	result := make([]string, 0, len(tokens))
	for _, token := range tokens {
		result = append(result, strings.ToLower(token))
	}
	return result
}
