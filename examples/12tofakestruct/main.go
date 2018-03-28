package main

import (
	"go/parser"
	"go/token"
	"log"

	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/codegen/transform"
	"golang.org/x/tools/go/loader"
)

func main() {
	source := `
package p

type I interface {
	F(s string)
	G(s string) string
	H(s string) (string, error)
	K(s string) (xxx string, yyy error)
}
`
	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, "f.go", source, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	c := loader.Config{Fset: fset}
	c.CreateFromFiles("p", astf)

	p, err := handwriting.NewFromPackagePath("x", handwriting.WithDryRun(), handwriting.WithConfig(&c))
	if err != nil {
		log.Fatal(err)
	}

	f := p.File("f.go")
	transform.EmitAsFakeStruct(f, "p/I", false)
	if err := p.Emit(); err != nil {
		log.Fatalf("%+v", err)
	}
}
