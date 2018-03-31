package struct2interface

import (
	"io"
)

// Buffer :
type Buffer interface {
	Bytes() []byte
	String() string
	Len() int
	Cap() int
	Truncate(n int)
	Reset()
	Grow(n int)
	Write(p []byte) (n int, err error)
	WriteString(s string) (n int, err error)
	ReadFrom(r io.Reader) (n int64, err error)
	WriteTo(w io.Writer) (n int64, err error)
	WriteByte(c byte) error
	WriteRune(r rune) (n int, err error)
	Read(p []byte) (n int, err error)
	Next(n int) []byte
	ReadByte() (byte, error)
	ReadRune() (r rune, size int, err error)
	UnreadRune() error
	UnreadByte() error
	ReadBytes(delim byte) (line []byte, err error)
	ReadString(delim byte) (line string, err error)
}
