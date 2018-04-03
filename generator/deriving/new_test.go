package deriving

import (
	"testing"

	"github.com/podhmo/handwriting"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/loader"
)

func TestNew(t *testing.T) {
	c := &loader.Config{}
	source := `
package new

import "time"

type I struct {
}

type J struct {
}

type K struct {
}

type S struct {
	I I
	J *J
	KS []K
	CreatedAt time.Time
}
`
	astf, err := c.ParseFile("f.go", source)
	require.NoError(t, err)

	pkgpath := "./testdata/new"
	c.CreateFromFiles(pkgpath, astf)

	p, err := handwriting.New(pkgpath, handwriting.WithConsoleOutput(), handwriting.WithConfig(c))
	require.NoError(t, err)

	f := p.File("s_output_stringer.go")
	require.NoError(t, GenerateNew(f, "S"))
	require.NoError(t, p.Emit())
}
