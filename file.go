package handwriting

import (
	"go/types"

	"github.com/podhmo/handwriting/generator/lookup"
	"github.com/podhmo/handwriting/generator/typesutil"
	"github.com/podhmo/handwriting/indent"
	"github.com/podhmo/handwriting/nameresolve"
	"golang.org/x/tools/go/loader"
)

// File :
type File struct {
	sourcefile *PlanningFile
	Prog       *loader.Program
	PkgInfo    *loader.PackageInfo
	Resolver   *nameresolve.File
	Out        *indent.Output
}

// TODO : import mapping

// Use :
func (f *File) Use(name string) (*types.Package, error) {
	if name != "" {
		for _, im := range f.sourcefile.imports {
			if im.Name == name {
				return lookup.Package(f.Prog, im.Path)
			}
		}
	}
	return lookup.Package(f.Prog, name)
}

// MustUse :
func (f *File) MustUse(name string) *types.Package {
	pkg, err := f.Use(name)
	if err != nil {
		panic(err.Error())
	}
	return pkg
}

// FileName :
func (f *File) FileName() string {
	return f.sourcefile.Filename
}

// CreateCaptureImportDetector :
func (f *File) CreateCaptureImportDetector() *typesutil.PackageDetector {
	return typesutil.NewPackageDetector(func(pkg *types.Package) {
		if pkg != nil {
			f.sourcefile.Import(pkg.Path())
		}
	})
}
