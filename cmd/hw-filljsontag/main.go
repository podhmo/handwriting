package main

import (
	"fmt"
	"go/build"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/lookup"
	"github.com/podhmo/handwriting/generator/namesutil"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type opt struct {
	fromPkg  string
	toPkg    string
	names    []string
	namefunc string
}

func guessPkg() (string, error) {
	curdir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(curdir)
	if err != nil {
		return "", err
	}
	for _, srcdir := range build.Default.SrcDirs() {
		if strings.HasPrefix(path, srcdir) {
			pkgname := strings.TrimLeft(strings.Replace(path, srcdir, "", 1), "/")
			return pkgname, nil
		}
	}
	return "", errors.Errorf("%q is not subdir of srcdirs(%q)", path, build.Default.SrcDirs())
}

func main() {
	var opt opt
	app := kingpin.New("filljsontag", "fill json tag")

	app.Flag("from", "from package").StringVar(&opt.fromPkg)
	// app.Flag("to", "to package").StringVar(&opt.toPkg)
	app.Flag("namefunc", "naming namefunc of json tag (camelcase,snakecase)").EnumVar(&opt.namefunc, "snakecase", "camelcase")
	app.Arg("name", "name").Required().StringsVar(&opt.names)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		app.FatalUsage(err.Error())
	}

	if opt.fromPkg == "" || opt.fromPkg == "." {
		pkg, err := guessPkg()
		if err != nil {
			app.FatalUsage(fmt.Sprintf("%v", err))
		}
		opt.fromPkg = pkg
		log.Printf("guess pkg name .. %q\n", opt.fromPkg)
	}

	var namefunc func(string) string
	switch opt.namefunc {
	case "snakecase":
		namefunc = namesutil.CamelToSnake // FooBar -> foo_bar
	default:
		namefunc = namesutil.ToLowerCamel // FooBar -> fooBar
	}

	if err := run(opt, namefunc); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(opt opt, namefunc func(string) string) error {
	var fnopts []func(*handwriting.Planner)
	if opt.toPkg == "" {
		fnopts = append(fnopts, handwriting.WithConsoleOutput())
		opt.toPkg = opt.fromPkg
	}
	p, err := handwriting.New(opt.toPkg, fnopts...)
	if err != nil {
		return err
	}

	f := p.File("snakecase.go")
	for _, name := range opt.names {
		name := name
		f.Code(func(f *handwriting.File) error {
			pkgref, err := f.Use(f.PkgInfo.Pkg.Path())
			if err != nil {
				return err
			}

			if name == "*" {
				for _, k := range pkgref.Scope().Names() {
					s, err := pkgref.LookupStruct(k)
					if err != nil {
						continue
					}
					if err := generateStruct(f, s, k, namefunc); err != nil {
						return err
					}
					f.Out.Newline()
				}
				return nil
			}

			s, err := pkgref.LookupStruct(name)
			if err != nil {
				return err
			}
			return generateStruct(f, s, name, namefunc)
		})
	}

	return p.Emit()
}

func generateStruct(f *handwriting.File, s *lookup.StructRef, name string, namefunc func(string) string) error {
	f.Out.WithBlock(fmt.Sprintf(`type %s struct`, name), func() {
		generateStructBody(f, s, namefunc)
	})
	return nil
}

func generateStructBody(f *handwriting.File, s *lookup.StructRef, namefunc func(string) string) {
	d := f.CreateCaptureImportDetector()

	i := 0
	s.IterateAllFields(func(field *types.Var) {
		d.Detect(field.Type())

		tag := s.Underlying.Tag(i)
		i++

		if !field.Exported() {
			f.Out.Printfln("%s %s `%s`", field.Name(), f.Resolver.TypeName(field.Type()), mergeTag(tag, "-"))
			return
		}

		if _, isFunc := field.Type().(*types.Signature); isFunc {
			f.Out.Printfln("%s %s `%s`", field.Name(), f.Resolver.TypeName(field.Type()), mergeTag(tag, "-"))
			return
		}
		if _, isisChan := field.Type().(*types.Chan); isisChan {
			f.Out.Printfln("%s %s `%s`", field.Name(), f.Resolver.TypeName(field.Type()), mergeTag(tag, "-"))
			return
		}

		if _, isRawStruct := field.Type().(*types.Struct); isRawStruct {
			f.Out.WithIndent(fmt.Sprintf("%s struct {", field.Name()), func() {
				sref := &lookup.StructRef{Underlying: field.Type().(*types.Struct)}
				generateStructBody(f, sref, namefunc)
			})
			f.Out.Printfln("} `%s`", mergeTag(tag, namefunc(field.Name())))
            return
		}

		f.Out.Printfln("%s %s `%s`", field.Name(), f.Resolver.TypeName(field.Type()), mergeTag(tag, namefunc(field.Name())))
	})
}

func mergeTag(tag, defaultname string) string {
	if tag == "" {
		return fmt.Sprintf(`json:"%s"`, defaultname)
	}
	v, ok := reflect.StructTag(tag).Lookup("json")
	if !ok {
		return fmt.Sprintf(`%s json:"%s"`, tag, defaultname)
	}
	return strings.Replace(tag, fmt.Sprintf(`json:"%s"`, v), fmt.Sprintf(`json:"%s"`, defaultname), 1)
}
