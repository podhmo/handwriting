package lookup

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/loader"
)

// PackageInfo :
func PackageInfo(prog *loader.Program, pkgpath string) (*loader.PackageInfo, error) {
	info := prog.Package(pkgpath)
	if info == nil {
		return nil, &lookupError{Type: Type("packageinfo"), Msg: fmt.Sprintf("%q is not found", pkgpath)}
	}
	return info, nil
}

// Package :
func Package(prog *loader.Program, pkgpath string) (*types.Package, error) {
	info := prog.Package(pkgpath)
	if info == nil {
		return nil, &lookupError{Type: Type("package"), Msg: fmt.Sprintf("%q is not found", pkgpath)}
	}
	return info.Pkg, nil
}
