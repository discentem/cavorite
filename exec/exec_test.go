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

func TestRealExecutor(t *testing.T) {
	tests := []struct {
		name           string
		exec           *RealExecutor
		getWriteCloser func(buf *bytes.Buffer) io.WriteCloser
		expectedErr    error
		expected       string
	}{
		{
			name: "stream to BufferCloser",
			exec: &RealExecutor{
				Cmd: exec.Command("bash", "test/artifacts/long_running.sh"),
			},
			getWriteCloser: func(buf *bytes.Buffer) io.WriteCloser {
				return &NopBufferCloser{Buffer: buf}
			},
			expectedErr: nil,
			expected:    "1 \n2 \n3 \n",
		},
		{
			name: "RealExecutor initialized with call to .Command()",
			getWriteCloser: func(buf *bytes.Buffer) io.WriteCloser {
				return &NopBufferCloser{Buffer: buf}
			},
			exec: func() *RealExecutor {
				re := &RealExecutor{}
				re.Command("bash", "test/artifacts/long_running.sh")
				return re
			}(),
			expectedErr: nil,
			expected:    "1 \n2 \n3 \n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if test.getWriteCloser == nil {
				t.Error("test.getWriteCloser must not be nil")
				t.Fail()
			}
			wc := test.getWriteCloser(out)
			err := test.exec.Stream(wc)
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expected, out.String())
		})
	}
}
func TestRealExecutorStreamNoPosters(t *testing.T) {
	orig := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = w
	re := RealExecutor{
		Cmd: exec.Command("bash", "test/artifacts/long_running.sh"),
	}
	// not passing anything to stream, rely on default behavior
	err = re.Stream()
	assert.NoError(t, err)
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	// confirm default behavior of re.Stream works
	assert.Equal(t, "1 \n2 \n3 \n", string(out))
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
