package stores

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"

	"github.com/discentem/cavorite/internal/testutils"
)

func fakeBucketClient(
	t *testing.T,
	storageOpts fakestorage.Options,
	bucketOpts fakestorage.CreateBucketOpts) *storage.Client {

	server, err := fakestorage.NewServerWithOptions(storageOpts)
	// create a gcs bucket on the server
	server.CreateBucketWithOpts(bucketOpts)
	assert.NoError(t, err)
	// get a gcs storage client talking to our fake gcs server
	client, err := storage.NewClient(
		context.Background(),
		option.WithoutAuthentication(),
		option.WithHTTPClient(server.HTTPClient()),
	)
	assert.NoError(t, err)
	return client

}

func TestGCSUpload(t *testing.T) {
	// create an in-memory file system with a small file
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	memfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"thing": {
			Content: []byte(`blah`),
			ModTime: &mTime,
		},
	})
	assert.NoError(t, err)

	// gcsstore with faked client
	store := GCSStore{
		fsys: *memfs,
		gcsClient: fakeBucketClient(
			t,
			fakestorage.Options{},
			fakestorage.CreateBucketOpts{
				Name: "test",
			},
		),
		Options: Options{
			BackendAddress: "test",
			// GCSStore.Upload will error if this isn't set
			MetadataFileExtension: "cfile",
		},
	}
	// upload object
	err = store.Upload(context.Background(), "thing")
	assert.NoError(t, err)

	f, err := store.fsys.Open("thing")
	assert.NoError(t, err)
	fstat, err := f.Stat()
	assert.NoError(t, err)
	b := make([]byte, fstat.Size())
	_, err = f.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`blah`), b)

}

func TestGCSRetrieve(t *testing.T) {
	// create an in-memory file system with a small file
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	memfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"thing/a.cfile": {
			ModTime: &mTime,
			Content: []byte(`{
				"name":"a",
				"checksum":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				"date_modified":"2014-11-12T11:45:26.371Z"
			}`),
		},
	})
	assert.NoError(t, err)

	// get GCSStore
	store := GCSStore{
		fsys: *memfs,
		gcsClient: fakeBucketClient(
			t,
			fakestorage.Options{
				InitialObjects: []fakestorage.Object{
					{
						ObjectAttrs: fakestorage.ObjectAttrs{
							BucketName: "test",
							Name:       "thing/a",
						},
						Content: []byte(
							`whatever`,
						),
					},
				},
			},
			fakestorage.CreateBucketOpts{
				Name: "test",
			},
		),
		Options: Options{
			BackendAddress: "test",
		},
	}
	// retrieve ensures the hash of the file matches a.cfile
	err = store.Retrieve(context.Background(), "thing/a.cfile")
	assert.NoError(t, err)

	// ensure the content of the file is correct
	b, _ := afero.ReadFile(*memfs, "thing/a")
	assert.Equal(t, `whatever`, string(b))
}
