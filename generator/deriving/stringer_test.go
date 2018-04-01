package deriving

import (
	"testing"

	"github.com/podhmo/handwriting"
	"github.com/stretchr/testify/require"
)

func TestStringer(t *testing.T) {
	pkgpath := "./testdata/enum"
	p, err := handwriting.New(pkgpath)
	require.NoError(t, err)

	f := p.File("s_output_stringer.go")
	require.NoError(t, GenerateStringer(f, "S"))
	require.NoError(t, p.Emit())
}
