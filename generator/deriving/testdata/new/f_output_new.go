package new

import (
	"time"
)

// New :
func S(i I, j *J, kS []K, createdAt time.Time) *S {
	return &S{
		I:         i,
		J:         j,
		KS:        kS,
		CreatedAt: createdAt,
	}
}
