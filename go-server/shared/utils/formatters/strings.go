package formatters

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ToPascalCase(s string) string {
	re := regexp.MustCompile(`[\s_-]+`)
	words := re.Split(s, -1)

	caser := cases.Title(language.English)

	for i, word := range words {
		if len(word) > 0 {
			words[i] = caser.String(word)
		}
	}
	return strings.Join(words, "")
}

func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	return strings.ToLower(string(pascal[0])) + pascal[1:]
}

func FormatArray(arr []string) string {
	if len(arr) == 0 {
		return "None"
	}
	return strings.Join(arr, ", ")
}
