package lookup

import (
	"fmt"
)

// Type :
type Type string

// lookupError :
type lookupError struct {
	Type  Type
	Msg  string
	Where string
}

// Error :
func (e *lookupError) Error() string {
	if e.Where == "" {
		return fmt.Sprintf("lookup %s, %s", e.Type, e.Msg)
	}
	return fmt.Sprintf("lookup %s, %s (where %q)", e.Type, e.Msg, e.Where)
}
