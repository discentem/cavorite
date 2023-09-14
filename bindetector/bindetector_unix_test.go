package bindetector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsBinary(t *testing.T) {
	require.True(t, isBinary([]byte(`binary`)))
	require.False(t, isBinary([]byte(`blah`)))

}
