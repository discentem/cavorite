package cli

import (
	"testing"

	"github.com/gonuts/go-shellquote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitCmd(t *testing.T) {
	expectedInitCmdArgs := `init ~/some/git/repo --backend_address http://127.0.0.1:9000/test --store_type=s3 --region="us-east-1"`

	initCmd := initCmd()

	// Split the args and handle bash escape characters
	args, err := shellquote.Split(expectedInitCmdArgs)
	require.NoError(t, err)

	// Traverse splits the beginning command and args into separate parts
	subCmd, subArgs, err := initCmd.Traverse(args)
	require.NoError(t, err)
	assert.NotNil(t, subCmd)
	assert.Equal(t, subCmd.UseLine(), "init")

	// Test the the subArgs equal the expected expectedRetrieveCmdArgs and flags
	assert.NoError(t, subCmd.ParseFlags(subArgs))

	// Test that subArgs expects the same args above
	assert.Equal(t, []string{
		"init",
		"~/some/git/repo",
		"--backend_address",
		"http://127.0.0.1:9000/test",
		"--store_type=s3",
		"--region=us-east-1",
	},
		subArgs,
	)

	// Check the required flags exist:
	assert.True(t, subCmd.PersistentFlags().Lookup("backend_address").Changed)
	assert.True(t, subCmd.PersistentFlags().Lookup("store_type").Changed)
	assert.True(t, subCmd.PersistentFlags().Lookup("region").Changed)
}
