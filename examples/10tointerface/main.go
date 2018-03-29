package main

import (
	"log"

	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/transform"
)

func main() {
	p, err := handwriting.NewFromPackagePath("github.com/podhmo/handwriting", handwriting.WithDryRun())
	if err != nil {
		log.Fatal(err)
	}
	f := p.File("f.go")

	// adding package import code for indirectly impored
	transform.GenerateInterface(f, "golang.org/x/tools/go/loader/Program", true)

	transform.GenerateInterface(f, "github.com/podhmo/handwriting/File", true)

	if err := p.Emit(); err != nil {
		log.Fatal(err)
	}
}
