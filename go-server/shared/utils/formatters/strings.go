package formatters

import "strings"

func ToPascalCase(s string) string {
	words := strings.Split(s, " ")
	for i, word := range words {
		words[i] = strings.Title(word)
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
