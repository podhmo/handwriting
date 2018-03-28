package nameresolve

import (
	"go/constant"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNamedTo(t *testing.T) {
	foo := types.NewPackage("github.com/xxx/go-foo", "foo")
	S := types.NewNamed(
		types.NewTypeName(token.NoPos, foo, "S", types.NewStruct(nil, nil)),
		types.NewStruct(nil, nil),
		nil,
	)
	foo.Scope().Insert(types.NewConst(token.NoPos, foo, "X", types.Universe.Lookup("string").Type(), constant.MakeString("x")))

	p := types.NewPackage("github.com/xxx/p", "p")

	t.Run("see package name", func(t *testing.T) {
		f := New(p).File(nil)
		t.Run("Name", func(t *testing.T) {
			assert.Exactly(t, "foo.X", f.Name(foo.Scope().Lookup("X")))
		})
		t.Run("TypeName", func(t *testing.T) {
			assert.Exactly(t, "*foo.S", f.TypeName(types.NewPointer(S)))
		})
	})

	t.Run("see imported name", func(t *testing.T) {
		source := `package p
	    import xfoo "github.com/xxx/go-foo"
	`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "", source, parser.ImportsOnly)
		require.NoError(t, err)

		f := New(p).File(file)

		t.Run("Name", func(t *testing.T) {
			assert.Exactly(t, "xfoo.X", f.Name(foo.Scope().Lookup("X")))
		})
		t.Run("TypeName", func(t *testing.T) {
			assert.Exactly(t, "*xfoo.S", f.TypeName(types.NewPointer(S)))
		})
	})

	t.Run("duplicated name", func(t *testing.T) {
		f := New(p).File(nil)
		f.Import(foo)

		yfoo := types.NewPackage("github.com/yyy/foo", "foo")
		M := types.NewNamed(
			types.NewTypeName(token.NoPos, yfoo, "M", types.NewStruct(nil, nil)),
			types.NewStruct(nil, nil),
			nil,
		)
		yfoo.Scope().Insert(types.NewConst(token.NoPos, yfoo, "Y", types.Universe.Lookup("string").Type(), constant.MakeString("y")))

		f.Import(yfoo)
		f.Import(foo)

		t.Run("Name", func(t *testing.T) {
			assert.Exactly(t, "foo1.Y", f.Name(yfoo.Scope().Lookup("Y")))
		})
		t.Run("TypeName", func(t *testing.T) {
			assert.Exactly(t, "*foo1.M", f.TypeName(types.NewPointer(M)))
		})
	})

	t.Run("duplicated import", func(t *testing.T) {
		f := New(p).File(nil)

		f.Import(foo)
		f.Import(foo)
		f.Import(foo)

		t.Run("Name", func(t *testing.T) {
			assert.Exactly(t, "foo.X", f.Name(foo.Scope().Lookup("X")))
		})
		t.Run("TypeName", func(t *testing.T) {
			assert.Exactly(t, "*foo.S", f.TypeName(types.NewPointer(S)))
		})
	})

}
