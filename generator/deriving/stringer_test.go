package deriving

import (
	"testing"

	"github.com/podhmo/handwriting"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/loader"
)

func TestStringer(t *testing.T) {
	c := &loader.Config{}
	source := `
package enum

// S :
type S string

//
const (
	X = S("x")
	Y = S("y")
	Z = S("z")
)
`
	astf, err := c.ParseFile("s.go", source)
	require.NoError(t, err)

	c.CreateFromFiles("./testdata/enum", astf)
	p, err := handwriting.New("./testdata/enum", handwriting.WithConfig(c))
	require.NoError(t, err)

	f := p.File("s_output_stringer.go")
	require.NoError(t, GenerateStringer(f, "S"))
	require.NoError(t, p.Emit())
}
