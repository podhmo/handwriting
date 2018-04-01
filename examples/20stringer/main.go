package main

import (
	"log"

	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/deriving"
	"golang.org/x/tools/go/loader"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	source := `
package p

type S string

const (
	X = S("x")
	Y = S("y")
	Z = S("z")
)
`
	// TODO : loading AST via planner
	c := &loader.Config{}
	astf, err := c.ParseFile("s.go", source)
	if err != nil {
		return err
	}
	c.CreateFromFiles("p", astf)

	p, err := handwriting.New(
		"p",
		handwriting.WithConsoleOutput(),
		handwriting.WithConfig(c),
	)
	if err != nil {
		return err
	}

	if err := deriving.GenerateStringer(p.File("s.go"), "S"); err != nil {
		return err
	}

	if err := p.Emit(); err != nil {
		return err
	}
	return nil
}
