package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/carolynvs/aferox"
	"github.com/discentem/cavorite/metadata"
	"github.com/discentem/cavorite/stores"
	"github.com/discentem/cavorite/testutils"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func localStoreUpload() *LocalStore {
	hlog := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	testfs, _ := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"/artifactStorage/blah": {
			Content: []byte(`someContent`),
		},
		"/git_repo/thing1": {Content: []byte(`bla`)},
		"/git_repo/thing2": {Content: []byte(`stuff`)},
	})

	return &LocalStore{
		logger: hlog,
		fsys:   aferox.NewAferox("/git_repo", *testfs),
		opts: &stores.Options{
			BackendAddress: "/artifactStorage",
		},
	}
}

func TestUpload(t *testing.T) {
	s := localStoreUpload()

	err := s.Upload(context.Background(), "thing1")
	require.NoError(t, err)
	b, err := afero.ReadFile(s.fsys, "/artifactStorage/thing1")
	require.NoError(t, err)
	require.Equal(t, []byte(`bla`), b)

	err = s.Upload(context.Background(), "thing2")
	require.NoError(t, err)
	b, err = afero.ReadFile(s.fsys, "../artifactStorage/thing2")
	require.NoError(t, err)
	require.Equal(t, []byte(`stuff`), b)
}

func localStoreRetrieve() *LocalStore {
	hlog := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	testfs, _ := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"/artifactStorage/someObject": {
			Content: []byte(`tla`),
		},
		"/git_repo/someObject.cfile": {Content: []byte(`{
			"name": "someObject",
			"checksum": "59e5ad2a03d2499749f7943c9dded0f303ad7542befef6d0aead8a7888587f66",
			"date_modified": "2014-11-12T11:45:26.371Z"
		   }`),
		},
	})

	return &LocalStore{
		logger: hlog,
		fsys:   aferox.NewAferox("/git_repo", *testfs),
		opts: &stores.Options{
			BackendAddress: "/artifactStorage",
		},
	}
}

func TestRetrieve(t *testing.T) {
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	s := localStoreRetrieve()
	err := s.Retrieve(context.Background(), metadata.CfileMetadataMap{
		"someObject.cfile": metadata.ObjectMetaData{
			Name:         "someObject",
			Checksum:     "59e5ad2a03d2499749f7943c9dded0f303ad7542befef6d0aead8a7888587f66",
			DateModified: mTime,
		},
	}, "someObject.cfile")
	require.NoError(t, err)
	_, err = s.fsys.Stat("/git_repo/someObject")
	require.NoError(t, err)
}

func TestRetrieveZeroCfiles(t *testing.T) {
	s := LocalStore{}
	err := s.Retrieve(context.Background(), metadata.CfileMetadataMap{})
	require.ErrorIs(t, err, ErrCfilesLengthZero)
}
