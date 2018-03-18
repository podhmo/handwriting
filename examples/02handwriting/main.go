package main

import (
	"os"

	"github.com/podhmo/handwriting"
)

func main() {
	h := handwriting.NewFromPath("github.com/podhmo/f")
	defer h.Commit(os.Stdout)
	f := h.File("fo.go")

	f.ImportWithName("fmt", "xfmt")
	f.Code(func(s *handwriting.State) error {
		o := s.Output

		// todo: nil safe (not panic)
		println := s.Lookup("fmt").Scope().Lookup("Println")

		o.Println("// F :")
		o.WithBlock("func F(x int)", func() {
			o.WithIfAndElse(
				"x % 2 == 0",
				func() { o.Printfln(`%s("even")`, s.File.Name(println)) },
				func() { o.Printfln(`%s("odd")`, s.File.Name(println)) },
			)
		})
		return nil
	})
}
