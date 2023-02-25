package metadata

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestParsePFile(t *testing.T) {
	fsys := afero.NewMemMapFs()
	f, err := fsys.Create(filepath.Join("repo", "thingy.pfile"))
	fmt.Println(fsys)
	assert.NoError(t, err)

	dateModified, err := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	assert.NoError(t, err)

	_, err = f.Write([]byte(`{
	"name":"a", 
	"checksum":"b",
	"date_modified":"2014-11-12T11:45:26.371Z"}`))
	assert.NoError(t, err)
	expect := ObjectMetaData{
		Name:         "a",
		Checksum:     "b",
		DateModified: dateModified,
	}
	actual, err := ParsePfile(fsys, filepath.Join("repo", "thingy"), "pfile")
	assert.NoError(t, err)
	assert.Equal(t, expect, *actual)
}
