package metadata

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestParsePFile(t *testing.T) {
	fsys := afero.NewMemMapFs()
	fname := filepath.Join("repo", "thingy.pfile")
	f, err := fsys.Create(fname)
	assert.NoError(t, err)

	modTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	// fake file modified time
	fsys.Chtimes(f.Name(), time.Time{}, modTime)
	assert.NoError(t, err)

	pfile := `{
		"name":"a", 
		"checksum":"b",
		"date_modified":"2014-11-12T11:45:26.371Z"}`

	_, err = f.Write([]byte(pfile))
	assert.NoError(t, err)
	expect := Object{
		Name:         "a",
		Checksum:     "b",
		DateModified: modTime,
	}
	actual, err := ParsePfile(fsys, filepath.Join("repo", "thingy"), "pfile")
	assert.NoError(t, err)
	assert.Equal(t, expect, *actual)
}
