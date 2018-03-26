package main

import (
	"log"

	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/codegen/transform"
)

func main() {
	p, err := handwriting.NewFromPackagePath("github.com/podhmo/handwriting", handwriting.WithDryRun())
	if err != nil {
		log.Fatal(err)
	}
	f := p.File("f.go")

	transform.BindToInterface(f, "golang.org/x/tools/go/loader/Program") // TODO: support indirect import package
	transform.BindToInterface(f, "github.com/podhmo/handwriting/File")
	f.Code(func(e *handwriting.Emitter) error {
		e.File.Import("go/type")
		e.File.Import("go/ast")
		return nil
	})
	if err := p.Emit(); err != nil {
		log.Fatal(err)
	}
}
