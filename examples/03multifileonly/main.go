package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/podhmo/handwriting/multifile"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	// opener := multifile.Must(multifile.Dir("./foo"))
	opener := multifile.Console(os.Stdout)

	if err := multifile.WriteFile(opener, "f0.go", func(w io.Writer) error {
		fmt.Fprintln(w, "f0")
		return nil
	}); err != nil {
		return err
	}

	if err := multifile.WriteFile(opener, "f1.go", func(w io.Writer) error {
		fmt.Fprintln(w, "f1")
		return nil
	}); err != nil {
		return err
	}
	return nil
}
