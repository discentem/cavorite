package exec

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type NopBufferCloser struct {
	*bytes.Buffer
}

func (b *NopBufferCloser) Close() error {
	return nil
}

func TestWriteOutput(t *testing.T) {
	tests := []struct {
		name     string
		in       io.ReadCloser
		expected string
	}{
		{
			name:     "test it works lol",
			in:       io.NopCloser(bytes.NewReader([]byte(`a b c`))),
			expected: "a b c\n",
		},
	}
	for _, test := range tests {
		out := bytes.Buffer{}
		err := WriteOutput(test.in, &NopBufferCloser{Buffer: &out})
		assert.NoError(t, err)
		assert.Equal(t, test.expected, out.String())
	}
}
