package bundle

import "io"

// Bundle :
func Bundle(opener Opener, filename string, fn func(w io.Writer) error) error {
	w, err := opener.Open(filename)
	if err != nil {
		return err
	}
	defer w.Close()
	return fn(w)
}
