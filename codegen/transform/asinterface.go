package transform

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/codegen/typesutil"
	"github.com/podhmo/handwriting/indent"
	"golang.org/x/tools/go/loader"
)

// EmitAsInterface :
func EmitAsInterface(f *handwriting.File, path string, exportedOnly bool) func(e *handwriting.Emitter) error {
	// <package path>/<name>
	elems := strings.Split(path, "/")
	pkgpath := strings.Join(elems[:len(elems)-1], "/")
	name := elems[len(elems)-1]
	f.Import(pkgpath)
	f.Code(func(e *handwriting.Emitter) error {
		return AsInterface(f, e.Prog.Package(pkgpath), name, e.Output, exportedOnly)
	})
	return nil
}

// AsInterface :
func AsInterface(f *handwriting.File, info *loader.PackageInfo, name string, o *indent.Output, exportedOnly bool) error {
	target := info.Pkg.Scope().Lookup(name)
	if target == nil {
		return errors.Errorf("%q is not found from package %q", name, info.Pkg.Path())
	}

	named, _ := target.Type().(*types.Named)
	if named == nil {
		return errors.Errorf("%q is not struct", name)
	}

	// todo : comment
	o.Printfln("// %s :", name)
	o.WithBlock(fmt.Sprintf("type %s interface", name), func() {
		n := named.NumMethods()

		// import pkg, if not imported yet.
		d := typesutil.NewPackageDetector(func(pkg *types.Package) {
			if pkg != nil {
				f.Import(pkg.Path())
			}
		})

		for i := 0; i < n; i++ {
			method := named.Method(i)
			if exportedOnly && !method.Exported() {
				continue
			}
			d.Detect(method.Type())
			o.Printfln("%s%s", method.Name(), strings.TrimPrefix(f.TypeName(method.Type()), "func"))
		}
	})
	return nil
}
