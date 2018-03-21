package transform

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/indent"
	"golang.org/x/tools/go/loader"
)

// BindToInterface :
func BindToInterface(f *handwriting.File, path string) func(e *handwriting.Emitter) error {
	// <package path>/<name>
	elems := strings.Split(path, "/")
	pkgpath := strings.Join(elems[:len(elems)-1], "/")
	name := elems[len(elems)-1]
	f.Import(pkgpath)
	f.Code(func(e *handwriting.Emitter) error {
		return ToInterface(f, e.Prog.Package(pkgpath), name, e.Output)
	})
	return nil
}

// ToInterface :
func ToInterface(f *handwriting.File, info *loader.PackageInfo, name string, o *indent.Output) error {
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
		for i := 0; i < n; i++ {
			method := named.Method(i)
			if !method.Exported() {
				continue
			}
			o.Printfln("%s %s", method.Name(), f.TypeName(method.Type()))
		}
	})
	return nil
}
