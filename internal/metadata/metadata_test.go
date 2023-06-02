package metadata

import (
	"bytes"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSHA256FromReader(t *testing.T) {
	r := bytes.NewReader([]byte(`blah`))
	got, err := SHA256FromReader(r)
	require.NoError(t, err)

	require.Equal(t,
		`8b7df143d91c716ecfa5fc1730022f6b421b05cedee8fd52b1fc65a96030ad52`,
		got,
	)
}

func TestGenerateFromReader(t *testing.T) {

}

func TestParseCFile(t *testing.T) {
	fsys := afero.NewMemMapFs()
	fname := filepath.Join("repo", "thingy.cfile")
	f, err := fsys.Create(fname)
	assert.NoError(t, err)

	modTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	// fake file modified time
	err = fsys.Chtimes(f.Name(), time.Time{}, modTime)
	assert.NoError(t, err)

	cfile := `{
		"name":"a",
		"checksum":"b",
		"date_modified":"2014-11-12T11:45:26.371Z"}`

	_, err = f.Write([]byte(cfile))
	assert.NoError(t, err)
	expect := ObjectMetaData{
		Name:         "a",
		Checksum:     "b",
		DateModified: modTime,
	}
	actual, err := ParseCfileWithExtension(fsys, filepath.Join("repo", "thingy"), "cfile")
	assert.NoError(t, err)
	assert.Equal(t, expect, *actual)
}
