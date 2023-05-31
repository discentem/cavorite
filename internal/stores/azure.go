package stores

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/spf13/afero"
)

type AzureBlobStore struct {
	Options         Options
	containerClient *azblob.Client
	storageAccount  string
	// containerName is equivalent to bucketName in S3 and GCS
	containerName string
	fsys          afero.Fs
}

func (s *AzureBlobStore) Upload(ctx context.Context, objects ...string) error   { return nil }
func (s *AzureBlobStore) Retrieve(ctx context.Context, objects ...string) error { return nil }
func (s *AzureBlobStore) GetOptions() Options                                   { return s.Options }
func (s *AzureBlobStore) GetFsys() afero.Fs                                     { return s.fsys }

func newAzureContainerClient(storageAccount string, containerName string, options azblob.ClientOptions) (*azblob.Client, error) {
	// We only support Azure CLI authentication.
	// In the future we could support multiple types with
	// azidentity.NewChainedTokenCredential()
	cred, err := azidentity.NewAzureCLICredential(nil)
	if err != nil {
		return nil, err
	}
	container, err := azblob.NewClient(
		fmt.Sprintf(
			"https://%s.blob.core.windows.net/%s",
			storageAccount,
			containerName,
		),
		cred,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return container, nil
}

func NewAzureBlobStore(storeOptions Options, storageAccount string, containerName string, azureBlobOptions azblob.ClientOptions) (*AzureBlobStore, error) {
	containerClient, err := newAzureContainerClient(
		storageAccount,
		containerName,
		azureBlobOptions,
	)
	if err != nil {
		return nil, err
	}
	return &AzureBlobStore{
		Options:         storeOptions,
		containerClient: containerClient,
		storageAccount:  storageAccount,
		containerName:   containerName,
	}, nil
}
