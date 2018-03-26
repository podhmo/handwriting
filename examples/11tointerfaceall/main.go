package main

import (
	"go/types"
	"log"

	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/codegen/transform"
)

func main() {
	p, err := handwriting.NewFromPackagePath("struct2interface", handwriting.WithDryRun())
	if err != nil {
		log.Fatal(err)
	}
	f := p.File("f.go")

	pkgname := "net/http/httptest"
	f.Import(pkgname)

	f.Code(func(e *handwriting.Emitter) error {
		ioPkginfo := e.Prog.Package(pkgname)
		s := ioPkginfo.Pkg.Scope()

		for _, name := range s.Names() {
			ob := s.Lookup(name)
			if !ob.Exported() {
				continue
			}
			if _, ok := ob.Type().Underlying().(*types.Struct); ok {
				exportedOnly := false
				transform.AsInterface(f, ioPkginfo, ob.Name(), e.Output, exportedOnly)
			}
		}
		return nil
	})

	if err := p.Emit(); err != nil {
		log.Fatal(err)
	}
}
