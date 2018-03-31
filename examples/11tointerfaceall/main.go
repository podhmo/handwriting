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
		g := transform.GeneratorForStructAsInterfaceNew(f)
		pkg, err := f.Use(pkgname)
		if err != nil {
			return err
		}

		var rerr error
		typesutil.IterateAllObjects(pkg.Package, func(ob types.Object) {
			if _, ok := ob.Type().Underlying().(*types.Struct); ok {
				exportedOnly := false
				if err := g.Generate(pkgname, ob.Name(), exportedOnly); err != nil {
					rerr = err
				}
			}
		})
		return nil
	})

	if err := p.Emit(); err != nil {
		log.Fatal(err)
	}
}
