package testutils

import (
	"bytes"
	"io/fs"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemMapFsWith(t *testing.T) {
	files := map[string]MapFile{
		"thing/whatever": {
			Content: []byte(`stuff`),
		},
		"thing/in/other/folder/thing": {
			Content: []byte(`otherStuff`),
		},
	}
	actual, err := MemMapFsWith(files)
	assert.NoError(t, err)

	expectedWalkMap := map[string][]byte{
		"thing/whatever":              []byte(`stuff`),
		"thing/in/other/folder/thing": []byte(`otherStuff`),
	}

	actualWalkMap := map[string][]byte{}
	err = afero.Walk(*actual, "", func(path string, _ fs.FileInfo, _ error) error {
		b, readErr := afero.ReadFile(*actual, path)
		assert.NoError(t, readErr)
		if bytes.Equal(b, []byte(``)) {
			return nil
		}
		actualWalkMap[path] = b
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, expectedWalkMap, actualWalkMap)

}
