package main

import (
	"log"

	"github.com/podhmo/handwriting"
)

// TODO:
// - finding unexisted package
// - finding unexisted member of package

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	h, err := handwriting.New("github.com/podhmo/g", handwriting.WithConsoleOutput())
	if err != nil {
		return err
	}

	f := h.File("fo.go")

	f.ImportWithName("fmt", "xfmt")
	f.Code(func(f *handwriting.File) error {
		// todo: nil safe (not panic)
		println := f.Prog.Package("fmt").Pkg.Scope().Lookup("Println")

		f.Out.Println("// F :")
		f.Out.WithBlock("func F(x int)", func() {
			f.Out.WithIfAndElse(
				"x % 2 == 0",
				func() { f.Out.Printfln(`%s("even")`, f.Resolver.Name(println)) },
				func() { f.Out.Printfln(`%s("odd")`, f.Resolver.Name(println)) },
			)
		})
		return nil
	})
	return h.Emit()
}
