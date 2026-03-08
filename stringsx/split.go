package stringsx

import (
	"unicode"
	"unicode/utf8"
)

type LetterType int

const (
	other LetterType = iota
	lower
	upper
	digit
)

func letterType(c rune) LetterType {
	switch {
	case unicode.IsLower(c):
		return lower
	case unicode.IsUpper(c):
		return upper
	case unicode.IsDigit(c):
		return digit
	default:
		return other
	}
}

func SplitToWords(s string) []string {
	if !utf8.ValidString(s) {
		return []string{s}
	}

	words := make([]string, 0)
	runes := [][]rune{{rune(s[0])}}
	index := 1 // index for runes

	for i := 1; i < len(s); i++ {
		prev := letterType(runes[index-1][0])
		curr := letterType(rune(s[i]))
		if prev == curr {
			runes[index-1] = append(runes[index-1], rune(s[i]))
		} else {
			runes = append(runes, []rune{rune(s[i])})
			index++
		}
	}

	for i := 0; i < len(runes)-1; i++ {
		curr := letterType(runes[i][0])
		next := letterType(runes[i+1][0])
		if curr == upper && next == lower {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}

	for _, word := range runes {
		if len(word) == 0 {
			continue
		}
		kind := letterType(word[0])
		if kind == other {
			continue
		}
		words = append(words, string(word))
	}

	return words
}

// SplitCamelCase splits the camelcase word and returns a list of words. It also
// supports digits. Both lower camel case and upper camel case are supported.
// For more info please check: http://en.wikipedia.org/wiki/CamelCase
// Splitting rules
//
//  1. If string is not valid UTF-8, return it without splitting as
//     single item array.
//  2. Assign all Unicode characters into one of 4 sets: lower case
//     letters, upper case letters, numbers, and all other characters.
//  3. Iterate through characters of string, introducing splits
//     between adjacent characters that belong to different sets.
//  4. Iterate through array of split strings, and if a given string
//     is upper case:
//     if subsequent string is lower case:
//     move last character of upper case string to beginning of
//     lower case string
func SplitCamelCase(src string) (parts []string) {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return []string{src}
	}

	parts = make([]string, 0)

	var (
		runes = make([][]rune, 0, len(src))
		prev  = LetterType(-1)
		curr  LetterType
	)

	// split into fields based on class of Unicode character
	for _, r := range src {
		curr = letterType(r)

		if prev == -1 {
			runes = append(runes, []rune{r})
			prev = curr
			continue
		}

		if curr == prev ||
			(curr == digit && (prev == upper || prev == lower)) ||
			(curr == lower && prev == digit) {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
			prev = curr
			continue
		}

		runes = append(runes, []rune{r})
		prev = curr
	}

	// handle upper case -> lower case sequences, e.g.
	// "HTTPS", "erver" -> "HTTP", "Server"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}

	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			parts = append(parts, string(s))
		}
	}

	return parts
}
