package metadata

import (
	"errors"
	"io"
	"path/filepath"
	"testing"
	"time"

	"github.com/discentem/cavorite/testutils"
	"github.com/google/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestWriteToFsys(t *testing.T) {
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
		Object:       "thing/a/whatever",
		Fsys:         fs,
		Fi:           fi,
		Extension:    "cfile",
		MetadataPath: "thing/a/whatever",
	})
	assert.NoError(t, err)

	b, _ := afero.ReadFile(*memfs, "thing/a/whatever.cfile")
	assert.Equal(t, string(b), `{
 "name": "thing/a/whatever",
 "checksum": "8b7df143d91c716ecfa5fc1730022f6b421b05cedee8fd52b1fc65a96030ad52",
 "date_modified": "2014-11-12T11:45:26.371Z"
}`)
}

type fsWithBrokenOpen struct {
	afero.Fs
}

func (f *fsWithBrokenOpen) Open(name string) (afero.File, error) {
	return nil, errors.New("open is intentionally broken in fsWithBrokenOpen implementation")
}

func TestHashFromCfileMatches(t *testing.T) {
	tests := []struct {
		name          string
		cfile         string
		fsys          afero.Fs
		expectedHash  string
		expectedMatch bool
		errExpected   func(err error) bool
	}{
		{
			name: "fsys cannot be nil",
			fsys: nil,
			errExpected: func(err error) bool {
				// we expect to get an error if fsys is nil
				return err != nil
			},
		},
		{
			name: "fsys.Open fails",
			fsys: func(t *testing.T) afero.Fs {
				fs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
					"thing": {
						Content: []byte(`blah`),
					},
				})
				assert.NoError(t, err)
				brokenOpenFs := fsWithBrokenOpen{
					Fs: *fs,
				}
				return &brokenOpenFs

			}(t),
			errExpected: func(err error) bool {
				return assert.Error(t, err)
			},
		},
		{
			name:          "hash matches",
			cfile:         "thing.cfile",
			expectedHash:  "8b7df143d91c716ecfa5fc1730022f6b421b05cedee8fd52b1fc65a96030ad52",
			expectedMatch: true,
			fsys: func(t *testing.T) afero.Fs {
				fs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
					"thing": {
						Content: []byte(`blah`),
					},
				})
				require.NoError(t, err)
				return *fs
			}(t),
		},
		{
			name:          "hash does not match",
			cfile:         "thing.cfile",
			expectedHash:  "8b7df143d91c716ecfa5fc1730022f6b421b05cedee8fd52b1fc65a96030ad52",
			expectedMatch: false,
			fsys: func(t *testing.T) afero.Fs {
				fs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
					"thing": {
						Content: []byte(`stuff`),
					},
				})
				require.NoError(t, err)
				return *fs
			}(t),
		},
	}
	t.Parallel()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger.Init("metadata_test", true, false, io.Discard)
			actual, err := HashFromCfileMatches(test.fsys, test.cfile, test.expectedHash)
			if test.errExpected != nil {
				assert.Equal(t, true, test.errExpected(err))
			}
			assert.Equal(t, test.expectedMatch, actual)

		})
	}
}
