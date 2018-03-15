package main

import (
	"os"

	"github.com/podhmo/handwriting"
)

func main() {
	o := handwriting.NewOutput(os.Stdout)
	o.Println("// F :")
	o.WithBlock("func F(x int)", func() {
		o.WithIfAndElse(
			"x % 2 == 0",
			func() { o.Println(`fmt.Println("even")`) },
			func() { o.Println(`fmt.Println("odd")`) },
		)
	})
}
