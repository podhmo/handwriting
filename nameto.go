package handwriting

import (
	"go/ast"
	"go/types"
	"sync"
)

// Prefixer :
type Prefixer struct {
	Pkg      *types.Package
	F        *ast.File
	Imported map[string]string // path -> prefix
	once     sync.Once
}

// NewPrefixer :
func NewPrefixer(pkg *types.Package, f *ast.File) *Prefixer {
	return &Prefixer{
		Pkg: pkg,
		F:   f,
	}
}

// NameTo :
func (p *Prefixer) NameTo(other *types.Package) string {
	if p.Pkg == other {
		return "" // same package; unqualified
	}
	p.once.Do(func() {
		// todo: conflict
		cache := map[string]string{}
		for _, is := range p.F.Imports {
			path := is.Path.Value[1 : len(is.Path.Value)-1]
			if is.Name != nil {
				cache[path] = is.Name.String()
			}
		}
		p.Imported = cache
	})
	if name, ok := p.Imported[other.Path()]; ok {
		return name
	}
	return other.Name() // todo: add import
}
