package bindetector

import (
	"errors"
	"io"
	"testing"

	shell "github.com/discentem/cavorite/exec"
	"github.com/hashicorp/go-multierror"
	multierr "github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

type fakeExecutor struct {
	output string
}

func (f fakeExecutor) Command(name string, arg ...string) {}

func (f fakeExecutor) Stream(posters ...io.WriteCloser) error {
	var result *multierror.Error
	for _, p := range posters {
		_, err := p.Write([]byte(f.output))
		result = multierr.Append(result, err)
	}
	return result.ErrorOrNil()
}

type brokenExecutor struct{}

func (b brokenExecutor) Command(name string, arg ...string) {}

func (b brokenExecutor) Stream(posters ...io.WriteCloser) error {
	return multierr.Append(nil, errors.New("broken executor"))
}
func TestFileIsABinary(t *testing.T) {
	tests := []struct {
		name        string
		executor    shell.Executor
		filepath    string
		expected    bool
		expectedErr bool
	}{
		{
			name:        "file is a binary",
			executor:    fakeExecutor{output: "application/x-executable; charset=binary\n"},
			filepath:    "whatever",
			expected:    true,
			expectedErr: false,
		},
		{
			name:        "file is not a binary",
			executor:    fakeExecutor{output: "text/plain; charset=us-ascii\n"},
			filepath:    "whatever",
			expected:    false,
			expectedErr: false,
		},
		{
			name:        "no output from binary check command",
			executor:    fakeExecutor{output: ""},
			filepath:    "whatever",
			expected:    false, // no output means we can't determine if the file is a binary
			expectedErr: true,  // we expect an error because the output is empty
		},
		{
			name:        "broken executor",
			executor:    brokenExecutor{},
			filepath:    "whatever",
			expected:    false, // we can't determine if the file is a binary because the executor
			expectedErr: true,  // we expect an error because the executor is broken
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := fileIsABinary(test.executor, test.filepath)
			assert.Equal(t, test.expected, actual, "expected %v, got %v", test.expected, actual)
			assert.Equal(t, test.expectedErr, err != nil, "expected error: %v", test.expectedErr)
		})
	}
}
