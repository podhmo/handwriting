package nameresolve

import "strings"

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
