package transform

import (
	"testing"

	"github.com/podhmo/handwriting"
	"github.com/podhmo/handwriting/multifile"
	"github.com/stretchr/testify/require"
)

func TestStructAsInterface(t *testing.T) {
	p, err := handwriting.New("struct2interface", handwriting.WithOpener(multifile.Must(multifile.Dir("testdata/output/struct2interface"))))
	require.NoError(t, err)

	exportedOnly := true

	require.NoError(t, GenerateStructAsInterface(p.File("bytes_buffer.go"), "bytes/Buffer", exportedOnly))
	require.NoError(t, GenerateStructAsInterface(p.File("encoding_json_encoder.go"), "encoding/json/Encoder", exportedOnly))
	require.NoError(t, GenerateStructAsInterface(p.File("encoding_json_decoder.go"), "encoding/json/Decoder", exportedOnly))

	require.NoError(t, p.Emit())
}
