package main

import (
	"os"

	"github.com/podhmo/handwriting/output"
)

func main() {
	o := output.New(os.Stdout)
	o.Println("// F :")
	o.WithBlock("func F(x int)", func() {
		o.WithIfAndElse(
			"x % 2 == 0",
			func() { o.Println(`fmt.Println("even")`) },
			func() { o.Println(`fmt.Println("odd")`) },
		)
	})
}
