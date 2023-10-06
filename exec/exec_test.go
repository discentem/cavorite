package exec

import (
	"bytes"
	"io"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

type NopBufferCloser struct {
	*bytes.Buffer
}

func (b *NopBufferCloser) Close() error {
	return nil
}

func TestStream(t *testing.T) {
	re := RealExecutor{
		Cmd: exec.Command("bash", "test/artifacts/long_running.sh"),
	}
	out := bytes.Buffer{}
	err := re.Stream(&NopBufferCloser{Buffer: &out})
	assert.NoError(t, err)
	assert.Equal(t, "1 \n2 \n3 \n", out.String())
}

func TestRealExecutorAsExecutorStream(t *testing.T) {
	re := RealExecutor{
		Cmd: exec.Command("bash", "test/artifacts/long_running.sh"),
	}
	// cast as interface
	e := Executor(&re)
	out := bytes.Buffer{}
	err := e.Stream(&NopBufferCloser{Buffer: &out})
	assert.NoError(t, err)
	assert.Equal(t, "1 \n2 \n3 \n", out.String())
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
