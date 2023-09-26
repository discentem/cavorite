package main

import (
	"context"
	"os"
	"testing"

	"github.com/carolynvs/aferox"
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
	// TODO(discentem): check contains of cfiles
	err = s.Upload(context.Background(), "thing2")
	require.NoError(t, err)

	b, err = afero.ReadFile(s.fsys, "../artifactStorage/thing2")
	require.NoError(t, err)
	require.Equal(t, []byte(`stuff`), b)
}
