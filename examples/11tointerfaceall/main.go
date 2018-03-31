package main

import (
	"go/types"
	"log"

	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/transform"
	"github.com/podhmo/handwriting/generator/typesutil"
)

func main() {
	p, err := handwriting.New("struct2interface", handwriting.WithConsoleOutput())
	if err != nil {
		log.Fatal(err)
	}
	f := p.File("f.go")

	pkgname := "net/http/httptest"

	f.Import(pkgname)
	f.Code(func(f *handwriting.File) error {
		ioPkg := f.Prog.Package(pkgname).Pkg
		typesutil.IterateAllObjects(ioPkg, func(ob types.Object) {
			if _, ok := ob.Type().Underlying().(*types.Struct); ok {
				exportedOnly := false
				transform.AsInterface(f, ioPkg, ob.Name(), f.Out, exportedOnly)
			}
		})
		return nil
	})

	if err := p.Emit(); err != nil {
		log.Fatal(err)
	}
}
