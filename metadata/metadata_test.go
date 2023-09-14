package metadata

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/discentem/cavorite/testutils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

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

func TestWriteMetadataToFsys(t *testing.T) {
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	memfs, _ := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"thing/a/whatever": {
			Content: []byte(`blah`),
			ModTime: &mTime,
		},
	})
	fs := *memfs
	fi, err := fs.Open("thing/a/whatever")
	assert.NoError(t, err)

	err = WriteToFsys(FsysWriteRequest{
		Object:    "thing/a/whatever",
		Fsys:      fs,
		Fi:        fi,
		Extension: "cfile",
	})
	assert.NoError(t, err)

	b, _ := afero.ReadFile(*memfs, "thing/a/whatever.cfile")
	assert.Equal(t, string(b), `{
 "name": "whatever",
 "checksum": "8b7df143d91c716ecfa5fc1730022f6b421b05cedee8fd52b1fc65a96030ad52",
 "date_modified": "2014-11-12T11:45:26.371Z"
}`)

}
