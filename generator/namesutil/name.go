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

// ToLowerCamel :
func ToLowerCamel(name string) string {
	if strings.HasSuffix(name, "ID") {
		name = name[:len(name)-2] + "Id"
	}
	for k, lower := range commonInitialisms {
		if strings.HasPrefix(name, k) {
			return lower + ToExported(name[len(k):])
		}
	}
	return ToUnexported(name)
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
		if len(s) > i-1 {
			if _, ok := commonInitialisms[s[:i]]; ok {
				initialism = s[:i]
			}
		}
	}
	return initialism
}

// copy from https://github.com/golang/lint
var commonInitialisms = map[string]string{
	"ACL":   "acl",
	"API":   "api",
	"ASCII": "ascii",
	"CPU":   "cpu",
	"CSS":   "css",
	"DNS":   "dns",
	"EOF":   "eof",
	"GUID":  "guid",
	"HTML":  "html",
	"HTTP":  "http",
	"HTTPS": "https",
	"ID":    "id",
	"IP":    "ip",
	"JSON":  "json",
	"LHS":   "lhs",
	"QPS":   "qps",
	"RAM":   "ram",
	"RHS":   "rhs",
	"RPC":   "rpc",
	"SLA":   "sla",
	"SMTP":  "smtp",
	"SQL":   "sql",
	"SSH":   "ssh",
	"TCP":   "tcp",
	"TLS":   "tls",
	"TTL":   "ttl",
	"UDP":   "udp",
	"UI":    "ui",
	"UID":   "uid",
	"UUID":  "uuid",
	"URI":   "uri",
	"URL":   "url",
	"UTF8":  "utf8",
	"VM":    "vm",
	"XML":   "xml",
	"XMPP":  "xmpp",
	"XSRF":  "xsrf",
	"XSS":   "xss",
}
