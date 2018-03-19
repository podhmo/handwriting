package indent

import (
	"fmt"
	"io"
	"strings"
)

// Output : leveled output writer
type Output struct {
	W      io.Writer
	i      int
	prefix string
}

// New :
func New(w io.Writer) *Output {
	return &Output{W: w}
}

// Indent :
func (w *Output) Indent() {
	w.i++
	w.prefix = strings.Repeat("	", w.i)
}

// UnIndent :
func (w *Output) UnIndent() {
	w.i--
	w.prefix = strings.Repeat("	", w.i)
}

// WithIndent :
func (w *Output) WithIndent(prefix string, callback func()) {
	w.Println(prefix)
	w.Indent()
	callback()
	w.UnIndent()
}

// WithBlock :
func (w *Output) WithBlock(prefix string, callback func()) {
	w.Println(prefix + " {")
	w.Indent()
	callback()
	w.UnIndent()
	w.Println("}")
}

// WithIfAndElse :
func (w *Output) WithIfAndElse(prefix string, callback func(), callback2 func()) {
	w.Println("if " + prefix + " {")
	w.Indent()
	callback()
	w.UnIndent()
	w.Println("} else {")
	w.Indent()
	callback2()
	w.UnIndent()
	w.Println("}")
}

// Newline :
func (w *Output) Newline() (int, error) {
	return io.WriteString(w.W, "\n")
}

// Println :
func (w *Output) Println(s string) (int, error) {
	return fmt.Fprintf(w.W, "%s%s\n", w.prefix, s)
}

// Printf :
func (w *Output) Printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(w.W, (w.prefix + format), args...)
}

// Printfln :
func (w *Output) Printfln(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(w.W, (w.prefix + format + "\n"), args...)
}
