package main

import (
	"fmt"
	"go/types"
	"log"

	"go/parser"
	"go/token"

	"github.com/podhmo/handwriting"
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
	info := prog.Package("p")

	qf := handwriting.NewPrefixer(info.Pkg, info.Files[0]).NameTo

	ob := prog.Package("golang.org/x/tools/go/loader").Pkg.Scope().Lookup("Config")
	fmt.Println(types.TypeString(types.NewPointer(ob.Type()), qf))
	fmt.Println(types.TypeString(types.NewPointer(ob.Type()), qf))
	return nil
}
