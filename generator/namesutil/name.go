package namesutil

import (
	"strings"
	"unicode"
)

// ToExported :
func ToExported(name string) string {
	if len(name) == 0 {
		return ""
	}
	// todo : see golint's code
	return strings.ToUpper(name[0:1]) + name[1:]
}

// ToUnexported :
func ToUnexported(name string) string {
	if len(name) == 0 {
		return ""
	}
	return strings.ToLower(name[0:1]) + name[1:]
}

// CamelToSnake :
func CamelToSnake(s string) string {
	var result string
	var words []string
	var lastPos int
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if i > 0 && unicode.IsUpper(rs[i]) {
			if initialism := startsWithInitialism(s[lastPos:]); initialism != "" {
				words = append(words, initialism)

				i += len(initialism) - 1
				lastPos = i
				continue
			}
			words = append(words, s[lastPos:i])
			lastPos = i
		}
	}
	if s[lastPos:] != "" {
		words = append(words, s[lastPos:])
	}
	for k, word := range words {
		if k > 0 {
			result += "_"
		}
		result += strings.ToLower(word)
	}
	return result
}

// startsWithInitialism returns the initialism if the given string begins with it
func startsWithInitialism(s string) string {
	var initialism string
	// the longest initialism is 5 char, the shortest 2
	for i := 1; i <= 5; i++ {
		if len(s) > i-1 && commonInitialisms[s[:i]] {
			initialism = s[:i]
		}
	}
	return initialism
}

// copy from https://github.com/golang/lint
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}
