package lookup

import (
	"fmt"

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
