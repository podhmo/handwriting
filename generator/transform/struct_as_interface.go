package transform

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/lookup"
	"github.com/podhmo/handwriting/generator/typesutil"
)

// GenerateStructAsInterface :
func GenerateStructAsInterface(f *handwriting.PlanningFile, path string, exportedOnly bool) error {
	// path = <package path>/<name>
	elems := strings.Split(path, "/")
	pkgpath := strings.Join(elems[:len(elems)-1], "/")
	name := elems[len(elems)-1]

	f.Root.Import(pkgpath)
	f.Code(func(f *handwriting.File) error {
		g := GeneratorForStructAsInterfaceNew(f)
		pkg, err := g.f.Use(pkgpath)
		if err != nil {
			return errors.Wrap(err, "lookup package")
		}
		strct, err := pkg.LookupStruct(name)
		if err != nil {
			return errors.Wrap(err, "lookup struct")
		}
		return g.Generate(strct, name, exportedOnly)
	})
	return nil
}

// GeneratorForStructAsInterface :
type GeneratorForStructAsInterface struct {
	f *handwriting.File
	d *typesutil.PackageDetector
}

// GeneratorForStructAsInterfaceNew :
func GeneratorForStructAsInterfaceNew(f *handwriting.File) *GeneratorForStructAsInterface {
	return &GeneratorForStructAsInterface{
		f: f,
		d: f.CreateCaptureImportDetector(),
	}
}

// Generate :
func (g *GeneratorForStructAsInterface) Generate(strct *lookup.StructRef, name string, exportedOnly bool) error {
	o := g.f.Out
	// todo : comment
	o.Printfln("// %s :", name)
	o.WithBlock(fmt.Sprintf("type %s interface", name), func() {
		strct.IterateMethods(typesutil.IterateModeFromBool(exportedOnly), func(method *types.Func) {
			g.d.Detect(method.Type())
			o.Printfln("%s%s", method.Name(), strings.TrimPrefix(g.f.Resolver.TypeName(method.Type()), "func"))
		})
	})
	return nil
}
