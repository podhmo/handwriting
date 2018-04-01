package enum

import (
	"fmt"
)

// String :
func (x S) String() string {
	switch x {
	case X:
		return "x"
	case Y:
		return "y"
	case Z:
		return "z"
	default:
		return fmt.Sprintf("S(%q)", string(x))
	}
}
