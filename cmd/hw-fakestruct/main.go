package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/generator/namesutil"
	"github.com/podhmo/handwriting/generator/transform"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type opt struct {
	fromPkg string
	toPkg   string
	names   []string
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
	app := kingpin.New("fakestruct", "generate fake-struct")

	app.Flag("from", "from package").StringVar(&opt.fromPkg)
	app.Flag("to", "to package").StringVar(&opt.toPkg)
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

	if err := run(opt); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(opt opt) error {
	var fnopts []func(*handwriting.Planner)
	if opt.toPkg == "" {
		fnopts = append(fnopts, handwriting.WithConsoleOutput())
		opt.toPkg = opt.fromPkg
	}
	p, err := handwriting.New(opt.toPkg, fnopts...)
	if err != nil {
		return err
	}

	exportedOnly := true
	for _, name := range opt.names {
		f := p.File(fmt.Sprintf("fake_%s.go", namesutil.CamelToSnake(name)))
		transform.GenerateFakeStruct(f, fmt.Sprintf("%s/%s", opt.fromPkg, name), exportedOnly)
	}

	return p.Emit()
}
