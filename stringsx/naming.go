package stringsx

import (
	"strings"
	"unicode"
)

// LowerSnakeCase e.g. i_am_a_10_years_senior
func LowerSnakeCase(name string) string {
	return rewords(name, func(res, word string, idx int) string {
		lower := strings.ToLower(word)
		if idx == 0 {
			return res + lower
		}
		return res + "_" + lower
	})
}

// UpperSnakeCase e.g. I_AM_A_10_YEARS_SENIOR
func UpperSnakeCase(name string) string {
	return rewords(name, func(res, word string, idx int) string {
		if idx == 0 {
			return res + strings.ToUpper(word)
		}
		return res + "_" + strings.ToUpper(word)
	})
}

// LowerCamelCase e.g. iAmA10YearsSenior
func LowerCamelCase(name string) string {
	return rewords(name, func(res, word string, idx int) string {
		lower := strings.ToLower(word)
		runes := []rune(lower)
		if idx > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		return res + string(runes)
	})
}

// UpperCamelCase e.g. IAmA10YearsSenior
func UpperCamelCase(name string) string {
	return rewords(name, func(res, word string, idx int) string {
		upper := strings.ToUpper(word)
		if _, ok := initialism[upper]; ok {
			return res + upper
		}
		word = strings.ToLower(word)
		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])
		return res + string(runes)
	})
}

// LowerDashJoint e.g. i-am-a-10-years-senior
func LowerDashJoint(name string) string {
	return rewords(name, func(res, word string, idx int) string {
		lower := strings.ToLower(word)
		if idx == 0 {
			return res + lower
		}
		return res + "-" + lower
	})
}

type joint func(result, word string, index int) string

func rewords(s string, fn joint) string {
	words := SplitToWords(s)
	ret := ""

	for i, word := range words {
		ret = fn(ret, word, i)
	}
	return ret
}

var initialism = map[string]struct{}{
	"ACL":   {},
	"API":   {},
	"ASCII": {},
	"CPU":   {},
	"CSS":   {},
	"DNS":   {},
	"EOF":   {},
	"GUID":  {},
	"HTML":  {},
	"HTTP":  {},
	"HTTPS": {},
	"ID":    {},
	"IP":    {},
	"JSON":  {},
	"LHS":   {},
	"QPS":   {},
	"RAM":   {},
	"RHS":   {},
	"RPC":   {},
	"SLA":   {},
	"SMTP":  {},
	"SQL":   {},
	"SSH":   {},
	"TCP":   {},
	"TLS":   {},
	"TTL":   {},
	"UDP":   {},
	"UI":    {},
	"UID":   {},
	"UUID":  {},
	"URI":   {},
	"URL":   {},
	"UTF8":  {},
	"VM":    {},
	"XML":   {},
	"XMPP":  {},
	"XSRF":  {},
	"XSS":   {},
	"QOS":   {},
}
