package stores

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/discentem/cavorite/internal/testutils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type aferoAzureBlobServer struct {
	containers map[string]afero.Fs
}

var (
	// ensure test server meets interface
	_ = azureBlobishClient(aferoAzureBlobServer{})
)

func (s aferoAzureBlobServer) UploadStream(
	ctx context.Context,
	containerName string,
	blobName string,
	body io.Reader, o *azblob.UploadStreamOptions) (azblob.UploadStreamResponse, error) {

	// check if the containerName exists in s.containers
	_, ok := s.containers[containerName]
	if !ok {
		return azblob.UploadStreamResponse{},
			fmt.Errorf(
				"%s does not exist in this aferoAzureBlobServer",
				containerName,
			)
	}
	b, err := io.ReadAll(body)
	if err != nil {
		return azblob.UploadStreamResponse{}, err
	}

	// create a filesystem for container referenced in input
	containerfs, _ := testutils.MemMapFsWith(map[string]testutils.MapFile{
		blobName: {
			// write input body to bucketfs
			Content: b,
		},
	})

	// write containerfs to associated "container"
	s.containers[containerName] = *containerfs

	return azblob.UploadStreamResponse{}, nil
}

func (s aferoAzureBlobServer) DownloadStream(
	ctx context.Context,
	containerName string,
	blobName string,
	o *azblob.DownloadStreamOptions) (azblob.DownloadStreamResponse, error) {
	return azblob.DownloadStreamResponse{}, nil
}

func TestAzureBlobStoreUpload(t *testing.T) {
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	memfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"test": {
			Content: []byte("tree"),
			ModTime: &mTime,
		},
	})
	assert.NoError(t, err)

	fakeAzureBlobServer := aferoAzureBlobServer{
		containers: map[string]afero.Fs{
			// create a bucket in our fake azure blob server
			"test": afero.NewMemMapFs(),
		},
	}
	store := AzureBlobStore{
		Options: Options{
			BackendAddress:        "http://whatever/test",
			MetadataFileExtension: "cfile",
		},
		fsys:            *memfs,
		containerClient: fakeAzureBlobServer,
	}
	err = store.Upload(context.Background(), "test")
	require.NoError(t, err)
	b, _ := afero.ReadFile(*memfs, "test.cfile")
	assert.Equal(t, `{
 "name": "test",
 "checksum": "dc9c5edb8b2d479e697b4b0b8ab874f32b325138598ce9e7b759eb8292110622",
 "date_modified": "2014-11-12T11:45:26.371Z"
}`, string(b))
}
