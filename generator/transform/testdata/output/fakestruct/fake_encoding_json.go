package fakestruct

// fakeMarshaler is fake struct of Marshaler
type fakeMarshaler struct {
	marshalJSON func() ([]byte, error)
}

// MarshalJSON :
func (x *fakeMarshaler) MarshalJSON() ([]byte, error) {
	return x.marshalJSON()
}

// fakeToken is fake struct of Token
type fakeToken struct {
}

// fakeUnmarshaler is fake struct of Unmarshaler
type fakeUnmarshaler struct {
	unmarshalJSON func([]byte) error
}

// UnmarshalJSON :
func (x *fakeUnmarshaler) UnmarshalJSON(v0 []byte) error {
	return x.unmarshalJSON(v0)
}
