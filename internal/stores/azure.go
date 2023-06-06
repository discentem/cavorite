package stores

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/discentem/cavorite/internal/fileutils"
	"github.com/discentem/cavorite/internal/metadata"
	"github.com/google/logger"
	"github.com/spf13/afero"

	progressbar "github.com/schollz/progressbar/v3"
)

// azureBlobishClient is derived from https://github.com/Azure/azure-sdk-for-go/blob/sdk/storage/azblob/v1.0.0/sdk/storage/azblob/client.go#L34
type azureBlobishClient interface {
	// DownloadBuffer(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.DownloadBufferOptions) (int64, error)
	// UploadBuffer(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.UploadBufferOptions) (azblob.UploadBufferResponse, error)
	UploadStream(ctx context.Context, containerName string, blobName string, body io.Reader, o *azblob.UploadStreamOptions) (azblob.UploadStreamResponse, error)
	DownloadStream(ctx context.Context, containerName string, blobName string, o *azblob.DownloadStreamOptions) (azblob.DownloadStreamResponse, error)
	// NewBlobClient(blobName string) *blob.Client
}

type AzureBlobStore struct {
	Options         Options
	containerClient *container.Client
	fsys            afero.Fs
}

func (s *AzureBlobStore) GetOptions() Options { return s.Options }
func (s *AzureBlobStore) GetFsys() afero.Fs   { return s.fsys }

func bytesTransferredFn(w io.Writer, progbar *progressbar.ProgressBar) func(bytesTransferred int64) {
	return func(bytesTransferred int64) {
		progbar.Set64(bytesTransferred)
		w.Write([]byte(progbar.String()))
	}
}

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

		blobClient := s.containerClient.NewBlockBlobClient(o)
		stat, err := f.Stat()
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}
		b, err := fileutils.BytesFromAferoFile(f)
		if err != nil {
			return err
		}
		progbar := progressbar.DefaultBytesSilent(stat.Size(), o)

		blobClient.UploadBuffer(ctx, b, &blockblob.UploadBufferOptions{
			Progress: bytesTransferredFn(os.Stdout, progbar),
		})
		fmt.Println(progbar.String())

		if err != nil {
			if err := cleanupFn(); err != nil {
				return err
			}
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
			logger.Infof("containerName: %s", containerName)

			blobClient := s.containerClient.NewBlobClient(o)
			blobProps, err := blobClient.GetProperties(ctx, nil)
			size := blobProps.ContentLength
			if err != nil {
				return err
			}
			if err := f.Truncate(*size); err != nil {
				return err
			}

			var b []byte
			progbar := progressbar.DefaultBytesSilent(*size, o)

			// Download the file
			_, err = blobClient.DownloadBuffer(
				ctx,
				b,
				&blob.DownloadBufferOptions{
					Progress: bytesTransferredFn(os.Stdout, progbar),
				})
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

func newAZContainerClient(backendAddress string) (*container.Client, error) {
	u, err := url.Parse(backendAddress)
	if err != nil {
		return nil, err
	}
	// We only support Azure CLI authentication.
	// In the future we could support multiple types with
	// azidentity.NewChainedTokenCredential()
	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return nil, err
	}
	containerURL := fmt.Sprintf("https://%s/%s", u.Host, path.Base(backendAddress))

	container, err := container.NewClient(
		// Construct container url
		containerURL,
		cred,
		&container.ClientOptions{},
	)
	if err != nil {
		return nil, err
	}
	return container, nil
}

func NewAzureBlobStore(ctx context.Context, fsys afero.Fs, storeOpts Options, azureBlobOptions azblob.ClientOptions) (*AzureBlobStore, error) {
	containerClient, err := newAZContainerClient(
		storeOpts.BackendAddress,
	)
	if err != nil {
		return nil, err
	}
	return &AzureBlobStore{
		Options:         storeOpts,
		containerClient: containerClient,
		fsys:            fsys,
	}, nil
}
