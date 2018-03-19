package handwriting

import (
	"fmt"
	"go/types"
	"io"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting/indent"
	"github.com/podhmo/handwriting/multifile"
	"golang.org/x/tools/go/loader"
)

// Emitter :
type Emitter struct {
	Prog   *loader.Program
	Pkg    *loader.PackageInfo
	Opener multifile.Opener
	*indent.Output
	*File
}

// Emit :
func (e *Emitter) Emit(file *File) error {
	return multifile.WriteFile(e.Opener, file.filename, func(w io.Writer) error {
		e.File = file
		e.Output = indent.New(w)

		for _, ac := range e.File.Setups {
			if err := ac(e); err != nil {
				return errors.Wrap(err, fmt.Sprintf("setup in %q", e.File.Name))
			}
		}

		for _, s := range e.File.Headers {
			e.Output.Println(s)
		}
		e.Output.Printf("package %s\n", e.Pkg.Pkg.Name())
		if len(e.File.imports) > 0 {
			e.Output.Println("")
			e.Output.Println("import (")
			e.Output.Indent()
			// todo : sort
			for _, im := range e.File.imports {
				if im.Name == "" {
					e.Output.Printfln(`%q`, im.Path)
				} else {
					e.Output.Printfln(`%s %q`, im.Name, im.Path)
				}
			}
			e.Output.UnIndent()
			e.Output.Println(")")
			e.Output.Println("")
		}

		for _, ac := range e.File.Actions {
			if err := ac(e); err != nil {
				return errors.Wrap(err, fmt.Sprintf("action in %q", e.File.Name))
			}
		}
		return nil
	})
}

// Lookup :
func (e *Emitter) Lookup(name string) *types.Package {
	info := e.Prog.Package(name)
	if info == nil {
		return nil
	}
	return info.Pkg
}
