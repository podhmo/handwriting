package handwriting

import (
	"bytes"
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
	return multifile.WriteFile(e.Opener, file.Filename, func(w io.Writer) error {
		var body bytes.Buffer
		e.File = file
		e.Output = indent.New(&body)

		for _, ac := range e.File.Setups {
			if err := ac(e); err != nil {
				return errors.Wrap(err, fmt.Sprintf("setup in %q", e.File.Filename))
			}
		}
		for _, s := range e.File.Headers {
			e.Output.Println(s)
		}
		e.Output.Printf("package %s\n", e.Pkg.Pkg.Name())

		for _, ac := range e.File.Actions {
			if err := ac(e); err != nil {
				return errors.Wrap(err, fmt.Sprintf("action in %q", e.File.Filename))
			}
		}

		// emitting import clause, lazily
		o := indent.New(w)
		if len(e.File.imports) > 0 {
			o.Println("")
			o.Println("import (")
			o.Indent()
			// todo : sort
			for _, im := range e.File.imports {
				if im.Name == "" {
					o.Printfln(`%q`, im.Path)
				} else {
					o.Printfln(`%s %q`, im.Name, im.Path)
				}
			}
			o.UnIndent()
			o.Println(")")
			o.Println("")
		}
		io.Copy(w, &body)
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
