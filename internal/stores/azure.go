package stores

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureBlobStore struct {
	Options         Options
	containerClient *azblob.Client
	storageAccount  string
	// containerName is equivalent to bucketName in S3 and GCS
	containerName string
}

func newAzureContainerClient(
	storageAccount string,
	containerName string,
	cred azcore.TokenCredential,
	options azblob.ClientOptions) (
	*azblob.Client,
	error,
) {
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
