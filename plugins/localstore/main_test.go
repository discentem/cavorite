package main

import (
	"context"
	"os"
	"testing"

	"github.com/carolynvs/aferox"
	"github.com/discentem/cavorite/stores"
	"github.com/discentem/cavorite/testutils"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func localStore() *LocalStore {
	hlog := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	testfs, _ := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"/artifactStorage/blah": {
			Content: []byte(`someContent`),
		},
		"/git_repo/thing": {
			Content: []byte(`blah`),
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

func TestUpload(t *testing.T) {
	s := localStore()
	err := s.Upload(context.Background(), "thing")
	assert.NoError(t, err)
	t.Log(s.fsys.Stat("/artifactStorage/thing"))

}
