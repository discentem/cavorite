package stores

import (
	"context"
	"testing"
	"time"

	"github.com/discentem/cavorite/internal/testutils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type simpleStore struct {
	fsys afero.Fs
}

func (s simpleStore) Upload(ctx context.Context, objects ...string) error {
	return nil
}
func (s simpleStore) Retrieve(ctx context.Context, objects ...string) error {
	return nil
}
func (s simpleStore) GetOptions() (Options, error) {
	return Options{
		MetadataFileExtension: "cfile",
	}, nil
}
func (s simpleStore) GetFsys() (afero.Fs, error) {
	return s.fsys, nil
}
func (s simpleStore) Close() error {
	return nil
}

func TestWriteMetadataToFsys(t *testing.T) {
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	memfs, _ := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"thing/a/whatever": {
			Content: []byte(`blah`),
			ModTime: &mTime,
		},
	})
	store := simpleStore{
		fsys: *memfs,
	}
	fs := *memfs
	fi, err := fs.Open("thing/a/whatever")
	assert.NoError(t, err)

	_, err = WriteMetadataToFsys(store, "thing/a/whatever", fi)
	assert.NoError(t, err)

	b, _ := afero.ReadFile(*memfs, "thing/a/whatever.cfile")
	assert.Equal(t, string(b), `{
 "name": "whatever",
 "checksum": "8b7df143d91c716ecfa5fc1730022f6b421b05cedee8fd52b1fc65a96030ad52",
 "date_modified": "2014-11-12T11:45:26.371Z"
}`)

}
