package cli

import (
	"testing"

	"github.com/gonuts/go-shellquote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadCmd(t *testing.T) {
	expectedUploadCmdArgs := "upload ./test_file_one ./test_file_two"

	uploadCmd := uploadCmd()

	// Split the args and handle bash escape characters
	args, err := shellquote.Split(expectedUploadCmdArgs)
	require.NoError(t, err)

	// Traverse splits the beginning command and args into separate parts
	subCmd, subArgs, err := uploadCmd.Traverse(args)
	require.NoError(t, err)
	assert.NotNil(t, subCmd)
	assert.Equal(t, subCmd.UseLine(), "upload")

	// Test the the subArgs equal the expected expectedUploadCmdArgs and flags
	assert.NoError(t, subCmd.ParseFlags(subArgs))

	// Test that subArgs expects the same args above
	assert.Equal(t, subArgs, []string{"upload", "./test_file_one", "./test_file_two"})
}
