package main

import (
	"fmt"
	"go/types"
	"log"

	"github.com/podhmo/handwriting/nameresolve"
	"golang.org/x/tools/go/loader"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	c := loader.Config{TypeCheckFuncBodies: func(path string) bool { return false }}

	{
		src := `package p
import (
	xloader "golang.org/x/tools/go/loader"
)
type C xloader.Config
`
		file, err := c.ParseFile("f0.go", src)
		if err != nil {
			return err
		}

		c.CreateFromFiles("p", file)
	}

	prog, err := c.Load()
	if err != nil {
		return err
	}

	{
		f := nameresolve.New(prog.Package("p").Pkg).File(file)

		ob := prog.Package("golang.org/x/tools/go/loader").Pkg.Scope().Lookup("Config")

		// xloader.Config
		fmt.Println(f.Name(ob))

		// *xloader.Config
		fmt.Println(f.TypeName(types.NewPointer(ob.Type())))
	}
	return nil
}
