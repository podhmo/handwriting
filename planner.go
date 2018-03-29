package handwriting

import (
	"go/build"
	"go/types"
	"sort"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting/multifile"
	"github.com/podhmo/handwriting/nameresolve"
	"golang.org/x/tools/go/loader"
)

// Planner :
type Planner struct {
	Pkg    *types.Package
	Config *loader.Config

	Resolver *nameresolve.Resolver
	Files    map[string]*File
	Opener   multifile.Opener
	// options
	TypeCheck bool
}

// New :
func New(pkg *types.Package, ops ...func(*Planner)) (*Planner, error) {
	h := &Planner{
		Pkg:      pkg,
		Resolver: nameresolve.New(pkg),
		Files:    map[string]*File{},
	}

	for _, op := range ops {
		op(h)
	}

	if h.Config == nil {
		h.Config = &loader.Config{}
		if !h.TypeCheck {
			h.Config.TypeCheckFuncBodies = func(path string) bool {
				return false
			}
			h.Config.TypeChecker.DisableUnusedImportCheck = true
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

func (h *Planner) importSelf() {
	// check package, if existed, import as initial package (tentative)
	bpkg, err := build.Default.Import(h.Pkg.Path(), ".", build.FindOnly)
	if err != nil || len(bpkg.GoFiles) == 0 {
		h.Config.CreateFromFiles(h.Pkg.Path())
		return
	}
	h.Config.Import(h.Pkg.Path())
}

// Emit :
func (h *Planner) Emit() error {
	prog, err := h.Config.Load()

	if err != nil {
		return errors.Wrap(err, "emit, load")
	}
	r := &Emitter{
		Prog:    prog,
		PkgInfo: prog.Package(h.Pkg.Path()),
		Opener:  h.Opener,
	}

	if r.PkgInfo == nil {
		return errors.Errorf("%q package is not found", h.Pkg.Path())
	}

	if r.PkgInfo.Pkg.Name() == "" {
		r.PkgInfo.Pkg.SetName(h.Pkg.Name())
	}

	// dummy to concreate package (tentative)
	h.Pkg = r.PkgInfo.Pkg
	h.Resolver.Pkg = r.PkgInfo.Pkg

	files := make([]*File, 0, len(h.Files))
	for k := range h.Files {
		files = append(files, h.Files[k])
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Filename < files[j].Filename })

	for i := range files {
		if err := r.Emit(files[i]); err != nil {
			return err
		}
	}
	return nil
}

// File :
func (h *Planner) File(name string) *File {
	f, ok := h.Files[name]
	if !ok {
		f = &File{Filename: name, Root: h, File: h.Resolver.File(nil), used: map[string]struct{}{}}
		h.Files[name] = f
	}
	return f
}
