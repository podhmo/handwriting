package main

import (
	"os"

	"github.com/podhmo/handwriting"
)

// TODO:
// - finding unexisted package
// - finding unexisted member of package

func main() {
	h := handwriting.NewFromPath("github.com/podhmo/f")
	defer h.Commit(os.Stdout)
	// todo: defer h.Commit(PackageWriter("github.com/podhmo/f"))
	// todo: defer h.Commit(PhysicalFilePathWriter("./"))

	f := h.File("fo.go")

	f.ImportWithName("fmt", "xfmt")
	f.Code(func(s *handwriting.State) error {
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
}
