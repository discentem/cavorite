package stores

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/discentem/cavorite/fileutils"
	"github.com/discentem/cavorite/metadata"
	"github.com/google/logger"
	"github.com/spf13/afero"
)

// azureBlobishClient is derived from https://github.com/Azure/azure-sdk-for-go/blob/sdk/storage/azblob/v1.0.0/sdk/storage/azblob/client.go#L34
type azureBlobishClient interface {
	UploadStream(ctx context.Context, containerName string, blobName string, body io.Reader, o *azblob.UploadStreamOptions) (azblob.UploadStreamResponse, error)
	DownloadStream(ctx context.Context, containerName string, blobName string, o *azblob.DownloadStreamOptions) (azblob.DownloadStreamResponse, error)
}

type AzureBlobStore struct {
	Options         Options
	containerClient azureBlobishClient
	fsys            afero.Fs
}

func (s *AzureBlobStore) GetOptions() (Options, error) { return s.Options, nil }
func (s *AzureBlobStore) GetFsys() (afero.Fs, error)   { return s.fsys, nil }

func (s *AzureBlobStore) Upload(ctx context.Context, objects ...string) error {
	for _, o := range objects {
		f, err := s.fsys.Open(o)
		if err != nil {
			return err
		}
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			return nil
		}
		containerName := path.Base(s.Options.BackendAddress)
		_, err = s.containerClient.UploadStream(
			ctx,
			containerName,
			o,
			f,
			&blockblob.UploadStreamOptions{
				Concurrency: 25,
			},
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
		f, err := fileutils.OpenOrCreateFile(s.fsys, objectPath)
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
			// Download the file
			resp, err := s.containerClient.DownloadStream(
				ctx,
				containerName,
				objectPath,
				&blob.DownloadStreamOptions{})
			if err != nil {
				return err
			}
			if resp.Body == nil {
				return fmt.Errorf("blob %q was nil", o)
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
			return metadata.ErrRetrieveFailureHashMismatch
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
		// &azblob.ClientOptions{
		// 	ClientOptions: policy.ClientOptions{
		// 		Retry: policy.RetryOptions{
		// 			TryTimeout:    time.Second * 5,
		// 			MaxRetryDelay: time.Second * 10,
		// 		},
		// 	},
		// },
		nil,
	)
	if err != nil {
		return nil, err
	}
	return container, nil
}

func NewAzureBlobStore(ctx context.Context, fsys afero.Fs, storeOpts Options, azureBlobOptions azblob.ClientOptions) (*AzureBlobStore, error) {
	u, err := url.Parse(storeOpts.BackendAddress)
	if err != nil {
		return nil, err
	}
	containerClient, err := newAzureContainerClient(
		fmt.Sprintf("https://%s/", u.Host),
		azureBlobOptions,
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

func (s *AzureBlobStore) Close() error {
	// FIXME: implement
	return nil
}
