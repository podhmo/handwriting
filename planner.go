package handwriting

import (
	"go/build"
	"go/types"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"os"

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
	return createPlanner(pkg, ops...)
}

// Must :
func Must(p *Planner, err error) *Planner {
	if err != nil {
		panic(err)
	}
	return p
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

// createPlanner :
func createPlanner(pkg *types.Package, ops ...func(*Planner)) (*Planner, error) {
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
	importable := false
	path := h.Pkg.Path()

	for _, spec := range h.Config.CreatePkgs {
		if spec.Path == path {
			return
		}
	}

	if build.IsLocalImport(path) {
		_, err := os.Stat(path)
		if err == nil {
			importable = true
		}
	} else {
		for _, srcdir := range build.Default.SrcDirs() {
			if _, err := os.Stat(filepath.Join(srcdir, path)); err == nil {
				importable = true
			}
		}
	}

	if importable {
		h.Config.Import(path)
		return
	}

	log.Printf("package %s is not found, creating.", h.Pkg.Path())
	h.Config.CreateFromFiles(h.Pkg.Path())
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

	for _, info := range prog.Created {
		if err := emitter.EmitCreated(info); err != nil {
			return err
		}
	}

	files := make([]*PlanningFile, 0, len(h.Files))
	for k := range h.Files {
		files = append(files, h.Files[k])
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Filename < files[j].Filename })

	for i := range files {
		if err := emitter.EmitFile(files[i]); err != nil {
			return err
		}
	}
	return nil
}
