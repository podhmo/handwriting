package handwriting

import (
	"go/types"

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

// // Package :
// func (f *File) Package(name string) *loader.PackageInfo {
// 	return f.Prog.Package(name)
// }

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
