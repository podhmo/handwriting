package handwriting

import (
	"go/build"
	"go/types"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting/multifile"
	"github.com/podhmo/handwriting/nameresolve"
	"github.com/podhmo/handwriting/shorthand"
	"golang.org/x/tools/go/loader"
)

// Planner :
type Planner struct {
	Pkg    *types.Package
	Config *loader.Config

	Files  map[string]*PlanningFile
	Opener multifile.Opener
	// options
	TypeCheck bool
}

// New :
func New(path string, ops ...func(*Planner)) (*Planner, error) {
	elems := strings.Split(path, "/")
	pkg := types.NewPackage(path, elems[len(elems)-1])
	return createPackage(pkg, ops...)
}

// WithConfig :
func WithConfig(c *loader.Config) func(*Planner) {
	return func(h *Planner) {
		h.Config = c
		h.importSelf()
	}
}

// WithOpener :
func WithOpener(o multifile.Opener) func(*Planner) {
	return func(h *Planner) {
		h.Opener = o
	}
}

// WithConsoleOutput :
func WithConsoleOutput() func(*Planner) {
	return WithOpener(multifile.Stdout())
}

// createPackage :
func createPackage(pkg *types.Package, ops ...func(*Planner)) (*Planner, error) {
	h := &Planner{
		Pkg:   pkg,
		Files: map[string]*PlanningFile{},
	}

	for _, op := range ops {
		op(h)
	}

	if h.Config == nil {
		if h.TypeCheck {
			h.Config = &loader.Config{}
		} else {
			h.Config = shorthand.NewUncheckConfig()
		}
		h.importSelf()
	}
	if h.Opener == nil {
		createIfNotExists := true
		opener, err := multifile.Package(pkg, createIfNotExists)
		if err != nil {
			return nil, err
		}
		h.Opener = opener
	}
	return h, nil
}

// Import :
func (h *Planner) Import(path string) {
	skipimport := false
	for _, pkgspec := range h.Config.CreatePkgs {
		if pkgspec.Path == path {
			skipimport = true
			break
		}
	}
	if !skipimport {
		h.Config.Import(path)
	}
}

func (h *Planner) importSelf() {
	// check package, if existed, import as initial package (tentative)
	bpkg, err := build.Default.Import(h.Pkg.Path(), ".", build.FindOnly)
	if err != nil || len(bpkg.GoFiles) == 0 {
		h.Config.CreateFromFiles(h.Pkg.Path())
		return
	}
	h.Config.Import(h.Pkg.Path())
}

// File :
func (h *Planner) File(name string) *PlanningFile {
	f, ok := h.Files[name]
	if !ok {
		f = &PlanningFile{Filename: name, Root: h, used: map[string]struct{}{}}
		h.Files[name] = f
	}
	return f
}

// createEmitter
func (h *Planner) createEmitter(prog *loader.Program, pkg *types.Package) (*Emitter, error) {
	emitter := &Emitter{
		Prog:    prog,
		PkgInfo: prog.Package(h.Pkg.Path()),
		Opener:  h.Opener,
	}
	if emitter.PkgInfo == nil {
		return nil, errors.Errorf("%q package is not found", pkg.Path())
	}
	emitter.Resolver = nameresolve.New(emitter.PkgInfo.Pkg)
	if emitter.PkgInfo.Pkg.Name() == "" {
		emitter.PkgInfo.Pkg.SetName(pkg.Name())
	}

	// dummy to concreate package (tentative)
	h.Pkg = emitter.PkgInfo.Pkg

	return emitter, nil
}

// Emit :
func (h *Planner) Emit() error {
	prog, err := h.Config.Load()

	if err != nil {
		return errors.Wrap(err, "emit, load")
	}

	emitter, err := h.createEmitter(prog, h.Pkg)
	if err != nil {
		return err
	}

	files := make([]*PlanningFile, 0, len(h.Files))
	for k := range h.Files {
		files = append(files, h.Files[k])
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Filename < files[j].Filename })

	for i := range files {
		if err := emitter.Emit(files[i]); err != nil {
			return err
		}
	}
	return nil
}
