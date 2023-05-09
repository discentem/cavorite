package cli

import (
	"testing"

	"github.com/gonuts/go-shellquote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveCommand(t *testing.T) {
	expectedRetrieveCmdArgs := "retrieve ./test_file_one ./test_file_two"

	retrieveCmd := retrieveCommand()

	// Split the args and handle bash escape characters
	args, err := shellquote.Split(expectedRetrieveCmdArgs)
	require.NoError(t, err)

	// Traverse splits the beginning command and args into separate parts
	subCmd, subArgs, err := retrieveCmd.Traverse(args)
	require.NoError(t, err)
	assert.NotNil(t, subCmd)
	assert.Equal(t, subCmd.UseLine(), "retrieve")

	// Test the the subArgs equal the expected expectedRetrieveCmdArgs and flags
	assert.NoError(t, subCmd.ParseFlags(subArgs))

	// Test that subArgs expects the same args above
	assert.Equal(t, subArgs, []string{"retrieve", "./test_file_one", "./test_file_two"})
}
