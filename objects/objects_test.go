package objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestAddPrefixOriginal(t *testing.T) {
	tests := []struct {
		name        string
		modifiedKey string
		prefix      string
		expected    string
	}{
		{
			name:        "single level key",
			modifiedKey: "stuff/thing",
			prefix:      "stuff",
			expected:    "thing",
		},
		{
			name:        "multi level key",
			modifiedKey: "apple/banana/strawberry",
			prefix:      "cake",
			expected:    "apple/banana/strawberry",
		},
		{},

		{
			name:        "empty prefix",
			modifiedKey: "whatever/thing/stuff",
			prefix:      "",
			expected:    "whatever/thing/stuff",
		},
	}

	for _, test := range tests {
		t.Log(test.name)
		modder := AddPrefixToKey{Prefix: test.prefix}
		actual := modder.Original(test.modifiedKey)
		assert.Equal(t, test.expected, actual)
	}

}

func TestModifyMultipleKeys(t *testing.T) {
	tests := []struct {
		name         string
		originalKeys []string
		expected     []string
		modifier     KeyModifier
		expectErr    bool
	}{
		{
			name:         "test modification on multiple keys",
			originalKeys: []string{"a", "b"},
			expected:     []string{"c/a", "c/b"},
			modifier:     AddPrefixToKey{Prefix: "c"},
			expectErr:    false,
		},
		{
			name:         "test nil modifier",
			originalKeys: []string{"a", "b"},
			expected:     []string{},
			modifier:     nil,
			expectErr:    true,
		},
		{
			name:         "test no-op modifier",
			originalKeys: []string{"a", "b"},
			expected:     []string{"a", "b"},
			modifier:     AddPrefixToKey{},
			expectErr:    false,
		},
	}
	for _, test := range tests {
		actual, err := ModifyMultipleKeys(test.modifier, test.originalKeys...)
		if test.expectErr {
			require.NotNil(t, err)
		}
		require.Equal(t, test.expected, actual)
	}
}
