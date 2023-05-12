package cli

import (
	"testing"

	"github.com/gonuts/go-shellquote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCmd(t *testing.T) {
	tests := []string{
		"cavorite",
		"cavorite init",
		"cavorite upload",
		"cavorite retrieve",
	}

	rootCmd := rootCmd()

	for _, tc := range tests {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			// Split the args and handle bash escape characters
			args, err := shellquote.Split(tc)
			require.NoError(t, err)

			// Traverse splits the beginning command and args into separate parts
			subCmd, subArgs, err := rootCmd.Traverse(args)
			require.NoError(t, err)
			assert.NotNil(t, subCmd)
			assert.Equal(t, subCmd.Use, "cavorite")
			assert.NoError(t, subCmd.ParseFlags(subArgs))
		})
	}
}
