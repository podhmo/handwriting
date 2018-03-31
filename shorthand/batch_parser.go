package shorthand

import (
	"go/ast"

	"golang.org/x/tools/go/loader"
)

// BatchParser :
type BatchParser struct {
	config       *loader.Config
	parseActions []func() (*ast.File, error)
	Files        map[string]*ast.File
}

// NewBatchParser :
func NewBatchParser(c *loader.Config) *BatchParser {
	return &BatchParser{
		config: c,
		Files:  map[string]*ast.File{},
	}
}

// Add :
func (bp *BatchParser) Add(filename string, src interface{}) {
	bp.parseActions = append(bp.parseActions, func() (*ast.File, error) {
		return bp.config.ParseFile(filename, src)
	})
}

// ParseFiles :
func (bp *BatchParser) ParseFiles() ([]*ast.File, error) {
	actions := bp.parseActions
	bp.parseActions = nil
	files := make([]*ast.File, 0, len(actions))

	for i := range actions {
		f, err := actions[i]()
		if err != nil {
			return files, err
		}
		files = append(files, f)
	}
	return files, nil
}

// BindCreatePackageFromFiles :
func (bp *BatchParser) BindCreatePackageFromFiles() error {
	files, err := bp.ParseFiles()
	if err != nil {
		return err
	}
	bp.config.CreateFromFiles("p", files...)
	return nil
}
