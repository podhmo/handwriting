package fakestruct

import (
	"io"
)

// fakeByteReader is fake struct of ByteReader
type fakeByteReader struct {
	readByte func() (byte, error)
}
// ReadByte :
func (x *fakeByteReader) ReadByte () (byte, error) {
	return x.readByte()
}


// fakeByteScanner is fake struct of ByteScanner
type fakeByteScanner struct {
	readByte func() (byte, error)
	unreadByte func() error
}
// ReadByte :
func (x *fakeByteScanner) ReadByte () (byte, error) {
	return x.readByte()
}
// UnreadByte :
func (x *fakeByteScanner) UnreadByte () error {
	return x.unreadByte()
}


// fakeByteWriter is fake struct of ByteWriter
type fakeByteWriter struct {
	writeByte func(c byte) error
}
// WriteByte :
func (x *fakeByteWriter) WriteByte (c byte) error {
	return x.writeByte(c)
}


// fakeCloser is fake struct of Closer
type fakeCloser struct {
	close func() error
}
// Close :
func (x *fakeCloser) Close () error {
	return x.close()
}


// fakeEOF is fake struct of EOF
type fakeEOF struct {
	error func() string
}
// Error :
func (x *fakeEOF) Error () string {
	return x.error()
}


// fakeReadCloser is fake struct of ReadCloser
type fakeReadCloser struct {
	close func() error
	read func(p []byte) (n int, err error)
}
// Close :
func (x *fakeReadCloser) Close () error {
	return x.close()
}
// Read :
func (x *fakeReadCloser) Read (p []byte) (n int, err error) {
	return x.read(p)
}


// fakeReadSeeker is fake struct of ReadSeeker
type fakeReadSeeker struct {
	read func(p []byte) (n int, err error)
	seek func(offset int64, whence int) (int64, error)
}
// Read :
func (x *fakeReadSeeker) Read (p []byte) (n int, err error) {
	return x.read(p)
}
// Seek :
func (x *fakeReadSeeker) Seek (offset int64, whence int) (int64, error) {
	return x.seek(offset, whence)
}


// fakeReadWriteCloser is fake struct of ReadWriteCloser
type fakeReadWriteCloser struct {
	close func() error
	read func(p []byte) (n int, err error)
	write func(p []byte) (n int, err error)
}
// Close :
func (x *fakeReadWriteCloser) Close () error {
	return x.close()
}
// Read :
func (x *fakeReadWriteCloser) Read (p []byte) (n int, err error) {
	return x.read(p)
}
// Write :
func (x *fakeReadWriteCloser) Write (p []byte) (n int, err error) {
	return x.write(p)
}


// fakeReadWriteSeeker is fake struct of ReadWriteSeeker
type fakeReadWriteSeeker struct {
	read func(p []byte) (n int, err error)
	seek func(offset int64, whence int) (int64, error)
	write func(p []byte) (n int, err error)
}
// Read :
func (x *fakeReadWriteSeeker) Read (p []byte) (n int, err error) {
	return x.read(p)
}
// Seek :
func (x *fakeReadWriteSeeker) Seek (offset int64, whence int) (int64, error) {
	return x.seek(offset, whence)
}
// Write :
func (x *fakeReadWriteSeeker) Write (p []byte) (n int, err error) {
	return x.write(p)
}


// fakeReadWriter is fake struct of ReadWriter
type fakeReadWriter struct {
	read func(p []byte) (n int, err error)
	write func(p []byte) (n int, err error)
}
// Read :
func (x *fakeReadWriter) Read (p []byte) (n int, err error) {
	return x.read(p)
}
// Write :
func (x *fakeReadWriter) Write (p []byte) (n int, err error) {
	return x.write(p)
}


// fakeReader is fake struct of Reader
type fakeReader struct {
	read func(p []byte) (n int, err error)
}
// Read :
func (x *fakeReader) Read (p []byte) (n int, err error) {
	return x.read(p)
}


// fakeReaderAt is fake struct of ReaderAt
type fakeReaderAt struct {
	readAt func(p []byte, off int64) (n int, err error)
}
// ReadAt :
func (x *fakeReaderAt) ReadAt (p []byte, off int64) (n int, err error) {
	return x.readAt(p, off)
}


// fakeReaderFrom is fake struct of ReaderFrom
type fakeReaderFrom struct {
	readFrom func(r io.Reader) (n int64, err error)
}
// ReadFrom :
func (x *fakeReaderFrom) ReadFrom (r io.Reader) (n int64, err error) {
	return x.readFrom(r)
}


// fakeRuneReader is fake struct of RuneReader
type fakeRuneReader struct {
	readRune func() (r rune, size int, err error)
}
// ReadRune :
func (x *fakeRuneReader) ReadRune () (r rune, size int, err error) {
	return x.readRune()
}


// fakeRuneScanner is fake struct of RuneScanner
type fakeRuneScanner struct {
	readRune func() (r rune, size int, err error)
	unreadRune func() error
}
// ReadRune :
func (x *fakeRuneScanner) ReadRune () (r rune, size int, err error) {
	return x.readRune()
}
// UnreadRune :
func (x *fakeRuneScanner) UnreadRune () error {
	return x.unreadRune()
}


// fakeSeeker is fake struct of Seeker
type fakeSeeker struct {
	seek func(offset int64, whence int) (int64, error)
}
// Seek :
func (x *fakeSeeker) Seek (offset int64, whence int) (int64, error) {
	return x.seek(offset, whence)
}


// fakeWriteCloser is fake struct of WriteCloser
type fakeWriteCloser struct {
	close func() error
	write func(p []byte) (n int, err error)
}
// Close :
func (x *fakeWriteCloser) Close () error {
	return x.close()
}
// Write :
func (x *fakeWriteCloser) Write (p []byte) (n int, err error) {
	return x.write(p)
}


// fakeWriteSeeker is fake struct of WriteSeeker
type fakeWriteSeeker struct {
	seek func(offset int64, whence int) (int64, error)
	write func(p []byte) (n int, err error)
}
// Seek :
func (x *fakeWriteSeeker) Seek (offset int64, whence int) (int64, error) {
	return x.seek(offset, whence)
}
// Write :
func (x *fakeWriteSeeker) Write (p []byte) (n int, err error) {
	return x.write(p)
}


// fakeWriter is fake struct of Writer
type fakeWriter struct {
	write func(p []byte) (n int, err error)
}
// Write :
func (x *fakeWriter) Write (p []byte) (n int, err error) {
	return x.write(p)
}


// fakeWriterAt is fake struct of WriterAt
type fakeWriterAt struct {
	writeAt func(p []byte, off int64) (n int, err error)
}
// WriteAt :
func (x *fakeWriterAt) WriteAt (p []byte, off int64) (n int, err error) {
	return x.writeAt(p, off)
}


// fakeWriterTo is fake struct of WriterTo
type fakeWriterTo struct {
	writeTo func(w io.Writer) (n int64, err error)
}
// WriteTo :
func (x *fakeWriterTo) WriteTo (w io.Writer) (n int64, err error) {
	return x.writeTo(w)
}


