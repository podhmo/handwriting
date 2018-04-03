package deriving

import (
	"fmt"
	"go/types"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/lookup"
	"github.com/podhmo/handwriting/generator/namesutil"
	"github.com/podhmo/handwriting/generator/typesutil"
)

// GenerateNew :
func GenerateNew(f *handwriting.PlanningFile, name string) error {
	f.Code(func(f *handwriting.File) error {
		g := GeneratorDerivingNewNew(f)
		pkgref := &lookup.PackageRef{Package: f.PkgInfo.Pkg}
		strct, err := pkgref.LookupStruct(name)
		if err != nil {
			return errors.Wrap(err, "lookup struct")
		}
		return g.Generate(strct, name)
	})
	return nil
}

// GeneratorDerivingNew :
type GeneratorDerivingNew struct {
	f *handwriting.File
	d *typesutil.PackageDetector
}

// GeneratorDerivingNewNew :
func GeneratorDerivingNewNew(f *handwriting.File) *GeneratorDerivingNew {
	return &GeneratorDerivingNew{
		f: f,
		d: f.CreateCaptureImportDetector(),
	}
}

// Generate :
func (g *GeneratorDerivingNew) Generate(strct *lookup.StructRef, name string) error {
	o := g.f.Out
	var params []*types.Var
	strct.IterateFields(func(x *types.Var) {

		params = append(params, types.NewVar(x.Pos(), x.Pkg(), namesutil.ToUnexported(x.Name()), x.Type()))
		g.d.Detect(x.Type())
	}, typesutil.ExportedOnly)

	o.Println("// New :")
	o.WithBlock(fmt.Sprintf("func %s %s %s", name, g.f.Resolver.TypeName(types.NewTuple(params...)), g.f.Resolver.TypeName(types.NewPointer(strct.Obj().Type()))), func() {
		o.WithBlock(fmt.Sprintf("return &%s", name), func() {
			strct.IterateFields(func(x *types.Var) {
				o.Printfln("%s: %s,", x.Name(), namesutil.ToUnexported(x.Name()))
			}, typesutil.ExportedOnly)

		})
	})
	return nil
}
