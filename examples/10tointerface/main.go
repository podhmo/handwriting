package main

import (
	"log"

	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/transform"
)

func main() {
	p, err := handwriting.New("github.com/podhmo/handwriting", handwriting.WithConsoleOutput())
	if err != nil {
		log.Fatal(err)
	}
	f := p.File("f.go")

	// adding package import code for indirectly impored
	transform.GenerateStructAsInterface(f, "golang.org/x/tools/go/loader/Program", true)

	transform.GenerateStructAsInterface(f, "github.com/podhmo/handwriting/File", true)

	if err := p.Emit(); err != nil {
		log.Fatal(err)
	}
}
