package util

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var typeCasePattern = regexp.MustCompile(`{{#(\w+)}}(\w+){{/(\w+)}}`)

const (
	Pascal string = "pascal_case"
	Snake  string = "snake_case"
	Camel  string = "camel_case"
	Kebab  string = "kebab_case"
	Title  string = "title_case"
)

var TypeCases = []string{
	Pascal,
	Snake,
	Camel,
	Kebab,
	Title,
}

// PascalCase Convert value into PascalCase
func PascalCase(value string) string {
	split := splitWord(value)
	titleCaser := cases.Title(language.English)

	for i, chara := range split {
		split[i] = titleCaser.String(chara)
	}

	return strings.Join(split, "")
}

// CamelCase Convert `value` into camelCase
func CamelCase(value string) string {
	split := splitWord(value)
	titleCaser := cases.Title(language.English)

	for i, chara := range split {
		// Convert the first `chara` of word into lowercase
		// the rest of it could be in Title case
		if i == 0 {
			split[i] = strings.ToLower(chara)
		} else {
			split[i] = titleCaser.String(chara)
		}
	}

	return strings.Join(split, "")
}

// SnakeCase Convert `value` into snake_case
func SnakeCase(value string) string {
	return strings.ToLower(strings.Join(splitWord(value), "_"))
}

// KebabCase Convert `value` into kebab-case
func KebabCase(value string) string {
	return strings.ToLower(strings.Join(splitWord(value), "-"))
}

// TitleCase Convert `value` into Title case
func TitleCase(value string) string {
	// Split `value` based on whitespace
	words := strings.Fields(value)
	titleCaser := cases.Title(language.English)

	for index := range words {
		words[index] = titleCaser.String(words[index])
	}

	return strings.Join(words, " ")
}

func splitWord(word string) []string {
	// All words
	var words []string
	// Current Word
	var currentWord []rune

	// Loop the `word`
	for _, chara := range word {
		// Check if `chara` is either letter or digit
		// If so, then add them to `currentWord`
		if unicode.IsLetter(chara) || unicode.IsDigit(chara) {
			currentWord = append(currentWord, chara)
		} else if len(currentWord) > 0 {
			// else `chara` considered as a new word
			// then appending `currentWord` into `words`
			words = append(words, string(currentWord))
			// Don't forget to reset the `currentWord`
			currentWord = nil
		}
	}

	// Handle if `currentWord` doesn't get reset inside the loop
	if len(currentWord) > 0 {
		words = append(words, string(currentWord))
	}

	return words
}
