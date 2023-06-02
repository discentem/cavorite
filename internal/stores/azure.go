package stores

import (
	"context"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/discentem/cavorite/internal/metadata"
	"github.com/google/logger"
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
func (s *AzureBlobStore) Retrieve(ctx context.Context, objects ...string) error {

	for _, o := range objects {
		// For Retrieve, the object is the cfile itself, which we derive the actual filename from
		objectPath := strings.TrimSuffix(o, filepath.Ext(o))
		// We will either read the file that already exists or download it because it
		// is missing
		f, err := openOrCreateFile(s.fsys, objectPath)
		if err != nil {
			return err
		}
		fileInfo, err := f.Stat()
		if err != nil {
			return err
		}
		if fileInfo.Size() > 0 {
			logger.Infof("%s already exists", objectPath)
		} else {
			containerName := path.Base(s.Options.BackendAddress)
			// Download the file
			resp, err := s.containerClient.DownloadStream(ctx, containerName, objectPath, nil)
			if err != nil {
				return err
			}
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			_, err = f.Write(b)
			if err != nil {
				return err
			}

		}
		// Get the hash for the downloaded file
		hash, err := metadata.SHA256FromReader(f)
		if err != nil {
			return err
		}
		// Get the metadata from the metadata file
		m, err := metadata.ParseCfile(s.fsys, o)
		if err != nil {
			return err
		}
		// If the hash of the downloaded file does not match the retrieved file, return an error
		if hash != m.Checksum {
			logger.Infof("Hash mismatch, got %s but expected %s", hash, m.Checksum)
			if err := s.fsys.Remove(objectPath); err != nil {
				return err
			}
			return ErrRetrieveFailureHashMismatch
		}
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil

}

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
