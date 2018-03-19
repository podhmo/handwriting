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
	h, err := handwriting.NewFromPackagePath("github.com/podhmo/f", handwriting.WithDryRun())
	if err != nil {
		return err
	}

	f := h.File("fo.go")

	f.ImportWithName("fmt", "xfmt")
	f.Code(func(s *handwriting.Emitter) error {
		// todo: nil safe (not panic)
		println := s.Lookup("fmt").Scope().Lookup("Println")

		s.Println("// F :")
		s.WithBlock("func F(x int)", func() {
			s.WithIfAndElse(
				"x % 2 == 0",
				func() { s.Printfln(`%s("even")`, s.Name(println)) },
				func() { s.Printfln(`%s("odd")`, s.Name(println)) },
			)
		})
		return nil
	})
	return h.Run()
}
