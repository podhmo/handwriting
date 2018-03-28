package main

import (
	"fmt"
	"log"

	"go/parser"
	"go/token"
	"go/types"

	"github.com/podhmo/handwriting/nameresolve"
	"golang.org/x/tools/go/loader"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	fset := token.NewFileSet()
	src := `package p
import (
	xloader "golang.org/x/tools/go/loader"
)
type C xloader.Config
`
	file, err := parser.ParseFile(fset, "f0.go", src, parser.ParseComments)
	if err != nil {
		return err
	}

	c := loader.Config{
		Fset:                fset,
		TypeCheckFuncBodies: func(path string) bool { return false },
	}
	c.CreateFromFiles("p", file)

	prog, err := c.Load()
	if err != nil {
		return err
	}
	f := nameresolve.New(prog.Package("p").Pkg).File(file)

	ob := prog.Package("golang.org/x/tools/go/loader").Pkg.Scope().Lookup("Config")

	// xloader.Config
	fmt.Println(f.Name(ob))

	// *xloader.Config
	fmt.Println(f.TypeName(types.NewPointer(ob.Type())))
	return nil
}
