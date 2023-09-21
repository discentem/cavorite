package objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPrefixModify(t *testing.T) {
	tests := []struct {
		name        string
		originalKey string
		prefix      string
		expected    string
	}{
		{
			name:        "single level key",
			originalKey: "thing",
			prefix:      "stuff",
			expected:    "stuff/thing",
		},
		{
			name:        "multi level key",
			originalKey: "apple/banana/strawberry",
			prefix:      "cake",
			expected:    "cake/apple/banana/strawberry",
		},
		{
			name:        "empty prefix",
			originalKey: "whatever/thing/stuff",
			prefix:      "",
			expected:    "whatever/thing/stuff",
		},
	}

	for _, test := range tests {
		modder := AddPrefixToKey{Prefix: test.prefix}
		actual := modder.Modify(test.originalKey)
		assert.Equal(t, test.expected, actual)
	}

}
