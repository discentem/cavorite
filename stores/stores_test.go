package stores

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferObjectPath(t *testing.T) {
	tests := []struct {
		name      string
		givenPath string
		expected  string
	}{
		{
			name:      "Trim .cfile suffix",
			givenPath: "blah.cfile",
			expected:  "blah",
		},
		{
			name:      "Trime some other suffix",
			givenPath: "blah.vfile",
			expected:  "blah",
		},
	}
	for _, test := range tests {
		t.Log(test.name)
		actual := inferObjPath(test.givenPath)
		assert.Equal(t, test.expected, actual)
	}
}
