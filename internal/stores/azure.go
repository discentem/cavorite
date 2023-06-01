package stores

import (
	"context"
	"io"
	"path"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/spf13/afero"
)

// azureBlobishClient is derived from https://github.com/Azure/azure-sdk-for-go/blob/sdk/storage/azblob/v1.0.0/sdk/storage/azblob/client.go#L34
type azureBlobishClient interface {
	// UploadStream copies the file held in io.Reader to the Blob at blockBlobClient.
	// A Context deadline or cancellation will cause this to error.
	UploadStream(ctx context.Context, containerName string, blobName string, body io.Reader, o *azblob.UploadStreamOptions) (azblob.UploadStreamResponse, error)
	// DownloadStream reads a range of bytes from a blob. The response also includes the blob's properties and metadata.
	// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob.
	DownloadStream(ctx context.Context, containerName string, blobName string, o *azblob.DownloadStreamOptions) (azblob.DownloadStreamResponse, error)
}

type AzureBlobStore struct {
	Options         Options
	containerClient azureBlobishClient
	fsys            afero.Fs
}

func (s *AzureBlobStore) GetOptions() Options { return s.Options }
func (s *AzureBlobStore) GetFsys() afero.Fs   { return s.fsys }

func (s *AzureBlobStore) Upload(ctx context.Context, objects ...string) error {
	for _, o := range objects {
		f, err := s.fsys.Open(o)
		if err != nil {
			return err
		}

		// cleanupFn is function that can be called if
		// uploading to blob storage fails. cleanupFn deletes the cfile
		// so that we don't retain a cfile without a corresponding binary
		cleanupFn, err := WriteMetadataToFsys(s, o, f)
		if err != nil {
			return err
		}

		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			if err := cleanupFn(); err != nil {
				return err
			}
			return err
		}

		containerName := path.Base(s.Options.BackendAddress)
		_, err = s.containerClient.UploadStream(
			ctx,
			containerName,
			o,
			f,
			nil,
		)
		if err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}

	return nil
}
func (s *AzureBlobStore) Retrieve(ctx context.Context, objects ...string) error { return nil }

func newAzureContainerClient(serviceURL string, options azblob.ClientOptions) (*azblob.Client, error) {
	// We only support Azure CLI authentication.
	// In the future we could support multiple types with
	// azidentity.NewChainedTokenCredential()
	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return nil, err
	}
	container, err := azblob.NewClient(
		serviceURL,
		cred,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return container, nil
}

func NewAzureBlobStore(backendAddress string, azureBlobOptions azblob.ClientOptions) (*AzureBlobStore, error) {
	containerClient, err := newAzureContainerClient(
		backendAddress,
		azureBlobOptions,
	)
	if err != nil {
		return nil, err
	}
	return &AzureBlobStore{
		Options: Options{
			BackendAddress:        backendAddress,
			MetadataFileExtension: "cfile",
		},
		containerClient: containerClient,
	}, nil
}
