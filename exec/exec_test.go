package exec

import (
	"bytes"
	"io"
	"os"
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

func TestRealExecutorStream(t *testing.T) {
	re := RealExecutor{
		Cmd: exec.Command("bash", "test/artifacts/long_running.sh"),
	}
	out := bytes.Buffer{}
	err := re.Stream(&NopBufferCloser{Buffer: &out})
	assert.NoError(t, err)
	assert.Equal(t, "1 \n2 \n3 \n", out.String())
}

func TestRealExecutorStreamNoPosters(t *testing.T) {
	orig := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = w
	re := RealExecutor{
		Cmd: exec.Command("bash", "test/artifacts/long_running.sh"),
	}
	err = re.Stream()
	assert.NoError(t, err)
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	assert.Equal(t, "1 \n2 \n3 \n", string(out))
}

func TestRealExecutorNilCmdStream(t *testing.T) {
	re := RealExecutor{}
	// if caller doesn't explicitly set RealExecutor.Cmd to some non-nil value in advance,
	// re.Command still correctly sets the path/args
	re.Command("bash", "test/artifacts/long_running.sh")
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
