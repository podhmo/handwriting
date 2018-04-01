package struct2interface

// Encoder :
type Encoder interface {
	Encode(v interface{}) error
	SetIndent(prefix string, indent string)
	SetEscapeHTML(on bool)
}
