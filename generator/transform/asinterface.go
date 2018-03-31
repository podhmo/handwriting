package transform

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/lookup"
	"github.com/podhmo/handwriting/generator/typesutil"
	"github.com/podhmo/handwriting/indent"
)

// GenerateInterface :
func GenerateInterface(f *handwriting.PlanningFile, path string, exportedOnly bool) func(e *handwriting.Emitter) error {
	// path = <package path>/<name>
	elems := strings.Split(path, "/")
	pkgpath := strings.Join(elems[:len(elems)-1], "/")
	name := elems[len(elems)-1]

	f.Import(pkgpath)
	f.Code(func(f *handwriting.File) error {
		info, err := lookup.PackageInfo(f.Prog, pkgpath)
		if info == nil {
			return errors.Wrap(err, "lookup pacakge")
		}
		return AsInterface(f, info.Pkg, name, f.Out, exportedOnly)
	})
	return nil
}

// AsInterface :
func AsInterface(f *handwriting.File, pkg *types.Package, name string, o *indent.Output, exportedOnly bool) error {
	target, err := lookup.Object(pkg, name)
	if err != nil {
		return errors.Wrap(err, "lookup target")
	}
	strct, err := lookup.AsStruct(target)
	if err != nil {
		return errors.Wrap(err, "lookup struct")
	}

	// import pkg, if not imported yet.
	d := f.CreateCaptureImportDetector()

	// todo : comment
	o.Printfln("// %s :", name)
	o.WithBlock(fmt.Sprintf("type %s interface", name), func() {
		strct.IterateMethods(typesutil.IterateModeFromBool(exportedOnly), func(method *types.Func) {
			d.Detect(method.Type())
			o.Printfln("%s%s", method.Name(), strings.TrimPrefix(f.Resolver.TypeName(method.Type()), "func"))
		})
	})
	return nil
}
