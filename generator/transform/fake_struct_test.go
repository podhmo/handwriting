package transform

import (
	"go/types"
	"testing"

	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/lookup"
	"github.com/podhmo/handwriting/multifile"
	"github.com/stretchr/testify/require"
)

func TestFakeStruct(t *testing.T) {
	candidates := []struct {
		pkgname  string
		filename string
	}{
		{
			pkgname: "io", filename: "fake_io.go",
		},
		{
			pkgname: "encoding/json", filename: "fake_encoding_json.go",
		},
	}

	for _, c := range candidates {
		c := c
		t.Run(c.pkgname, func(t *testing.T) {
			p, err := handwriting.New("fakestruct", handwriting.WithOpener(multifile.Must(multifile.Dir("testdata/output/fakestruct"))))
			require.NoError(t, err)

			exportedOnly := true
			arrived := map[*types.Interface]struct{}{}

			p.Import(c.pkgname)
			p.File(c.filename).Code(func(f *handwriting.File) error {
				g := GeneratorForFakeStructNew(f)
				iopkg, err := f.Use(c.pkgname)
				require.NoError(t, err)

				for _, name := range iopkg.Scope().Names() {
					ob := iopkg.Scope().Lookup(name)
					if !ob.Exported() {
						continue
					}
					if iface, ok := ob.Type().Underlying().(*types.Interface); ok {
						if _, ok := arrived[iface]; ok {
							continue
						}
						arrived[iface] = struct{}{}
						ref, err := lookup.AsInterface(ob)
						require.NoError(t, err)

						g.Generate(ref, name, "fake"+name, exportedOnly)
						f.Out.Newline()
						f.Out.Newline()
					}
				}
				return nil
			})
			require.NoError(t, p.Emit())
		})
	}
}
