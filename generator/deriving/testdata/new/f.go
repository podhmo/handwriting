package new

import "time"

type I struct {
}

type J struct {
}

type K struct {
}

type S struct {
	I		I
	J		*J
	KS		[]K
	CreatedAt	time.Time
}
