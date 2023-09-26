package stores

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/discentem/cavorite/metadata"
	"github.com/spf13/afero"
)

type StoreType string

const (
	StoreTypeUndefined StoreType = "undefined"
	StoreTypeS3        StoreType = "s3"
	StoreTypeGCS       StoreType = "gcs"
	StoreTypeAzureBlob StoreType = "azure"
	StoreTypeGoPlugin  StoreType = "plugin"
)

var (
	_ = Store(&S3Store{})
	// _ = Store(&GCSStore{})
	// _ = Store(&AzureBlobStore{})
	_ = Store(&PluggableStore{})
)

type StoreWithGetters interface {
	Store
	GetFsys() (afero.Fs, error)
}

type Store interface {
	Upload(ctx context.Context, keys ...string) error
	Retrieve(ctx context.Context, mmap metadata.CfileMetadataMap, keys ...string) error
	GetOptions() (Options, error)
	Close() error
}

func inferObjPath(cfilePath string) string {
	return strings.TrimSuffix(cfilePath, filepath.Ext(cfilePath))
}
