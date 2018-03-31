package main

import (
	"log"

	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/transform"
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
	c := loader.Config{}
	astf, err := c.ParseFile("f.go", source)
	if err != nil {
		log.Fatal(err)
	}
	c.CreateFromFiles("p", astf)

	p, err := handwriting.New("x", handwriting.WithConsoleOutput(), handwriting.WithConfig(&c))
	if err != nil {
		log.Fatal(err)
	}

	f := p.File("f.go")
	transform.GenerateFakeStruct(f, "p/I", false)
	if err := p.Emit(); err != nil {
		log.Fatalf("%+v", err)
	}
}
