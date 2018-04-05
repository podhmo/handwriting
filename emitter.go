package handwriting

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"go/printer"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting/indent"
	"github.com/podhmo/handwriting/multifile"
	"github.com/podhmo/handwriting/nameresolve"
	"golang.org/x/tools/go/loader"
)

// Emitter :
type Emitter struct {
	Prog     *loader.Program
	Resolver *nameresolve.Resolver
	PkgInfo  *loader.PackageInfo
	Opener   multifile.Opener

	// options
	Format func(r io.Reader, w io.Writer) error
}

// EmitCreated :
func (e *Emitter) EmitCreated(info *loader.PackageInfo) error {
	for _, f := range info.Files {
		name := e.Prog.Fset.File(f.Pos()).Name()
		if err := multifile.WriteFile(e.Opener, name, func(w io.Writer) error {
			return printer.Fprint(w, e.Prog.Fset, f)
		}); err != nil {
			return err
		}
	}
	return nil
}

// EmitFile :
func (e *Emitter) EmitFile(file *PlanningFile) error {
	return multifile.WriteFile(e.Opener, file.Filename, func(w io.Writer) error {
		var body bytes.Buffer
		f := &File{
			Prog:       e.Prog,
			PkgInfo:    e.PkgInfo,
			sourcefile: file,
			Resolver:   e.Resolver.File(nil),
			Out:        indent.New(&body),
		}
		for _, ac := range file.Setups {
			if err := ac(f); err != nil {
				return errors.Wrap(err, fmt.Sprintf("setup in %q", file.Filename))
			}
		}

		for _, ac := range file.Actions {
			if err := ac(f); err != nil {
				return errors.Wrap(err, fmt.Sprintf("action in %q", file.Filename))
			}
		}

		if e.Format == nil {
			// emitting import clause, lazily
			e.emitPrologue(f, w)
			io.Copy(w, &body)
			return nil
		}

		var out bytes.Buffer
		// emitting import clause, lazily
		e.emitPrologue(f, &out)
		io.Copy(&out, &body)
		return e.Format(&out, w)
	})
}

func (e *Emitter) emitPrologue(f *File, w io.Writer) {
	// todo : header comment (prologue note)

	sourcefile := f.sourcefile
	o := indent.New(w)
	o.Printf("package %s\n", f.PkgInfo.Pkg.Name())

	if len(sourcefile.imports) > 0 {
		o.Println("")
		o.Println("import (")
		o.Indent()

		// todo : sort
		for _, im := range sourcefile.imports {
			name := f.Resolver.Imported[im.Path]
			if im.Path == name || strings.HasSuffix(im.Path, "/"+name) {
				o.Printfln(`%q`, im.Path)
			} else {
				o.Printfln(`%s %q`, name, im.Path)
			}
		}
		o.UnIndent()
		o.Println(")")
		o.Println("")
	}
}
