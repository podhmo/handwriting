package shorthand

import (
	"golang.org/x/tools/go/loader"
)

// NewUncheckConfig :
func NewUncheckConfig() *loader.Config {
	c := &loader.Config{
		TypeCheckFuncBodies: func(path string) bool {
			return false
		},
	}
	c.TypeChecker.DisableUnusedImportCheck = true
	return c
}
