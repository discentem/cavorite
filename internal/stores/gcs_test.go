package stores

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

type mapFile struct {
	content []byte
	modTime *time.Time
}

// memMapFsWith creates a afero.MemMapFs from a map[string]mapFile
func memMapFsWith(files map[string]mapFile) (*afero.Fs, error) {
	memfsys := afero.NewMemMapFs()
	for fname, mfile := range files {
		afile, err := memfsys.Create(fname)
		if err != nil {
			return nil, err
		}
		_, err = afile.Write(mfile.content)
		if err != nil {
			return nil, err
		}
		if mfile.modTime != nil {
			err := memfsys.Chtimes(fname, time.Time{}, *mfile.modTime)
			if err != nil {
				return nil, err
			}
		}
	}
	return &memfsys, nil
}

// fakeGCSStore creates a GCSStore using fake-gcs-store
func fakeGCSStore(t *testing.T, fsys afero.Fs) GCSStore {
	// create a gcs fake server
	server, err := fakestorage.NewServerWithOptions(fakestorage.Options{})
	// create a gcs bucket on the server
	server.CreateBucketWithOpts(fakestorage.CreateBucketOpts{
		Name: "test",
	})
	assert.NoError(t, err)
	// get a gcs storage client talking to our fake gcs server
	client, err := storage.NewClient(
		context.Background(),
		option.WithoutAuthentication(),
		option.WithHTTPClient(server.HTTPClient()),
	)
	assert.NoError(t, err)
	assert.NoError(t, err)

	gcs := &GCSStore{
		fsys:      fsys,
		gcsClient: client,
		Options: Options{
			BackendAddress: "test",
		},
	}
	return *gcs
}

func TestGCSUpload(t *testing.T) {
	// create an in-memory file system with a small file
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	memfs, err := memMapFsWith(map[string]mapFile{
		"thing": {
			content: []byte(`blah`),
			modTime: &mTime,
		},
	})

	// get GCSStore
	store := fakeGCSStore(t, *memfs)
	// upload object
	err = store.Upload(context.Background(), "thing")
	assert.NoError(t, err)

	// get object handle
	obj := store.gcsClient.Bucket("test").Object("thing")
	r, err := obj.NewReader(context.Background())
	assert.NoError(t, err)

	buf := new(strings.Builder)
	_, err = io.Copy(buf, r)
	assert.NoError(t, err)
	// ensure object's content is correct in bucket
	assert.Equal(t, `blah`, buf.String())
}
