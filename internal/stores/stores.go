package stores

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/discentem/cavorite/internal/metadata"
	"github.com/google/logger"
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
	_ = Store(&s3Store{})
	_ = Store(&GCSStore{})
	_ = Store(&AzureBlobStore{})
	_ = Store(&PluggableStore{})
)

var (
	ErrRetrieveFailureHashMismatch = errors.New("hashes don't match, Retrieve aborted")
)

type StoreWithGetters interface {
	Store
	GetFsys() (afero.Fs, error)
}

type Store interface {
	Upload(ctx context.Context, objects ...string) error
	Retrieve(ctx context.Context, objects ...string) error
	GetOptions() (Options, error)
}

func openOrCreateFile(fsys afero.Fs, filename string) (afero.File, error) {
	file, err := fsys.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func inferObjPath(cfilePath string) string {
	return strings.TrimSuffix(cfilePath, filepath.Ext(cfilePath))
}

// WriteMetadata generates Cavorite metadata for obj and writes it to s.Fsys
func WriteMetadataToFsys(s StoreWithGetters, obj string, f afero.File) (cleanup func() error, err error) {
	opts, err := s.GetOptions()
	if err != nil {
		return nil, err
	}
	if opts.MetadataFileExtension == "" {
		return nil, metadata.ErrFileExtensionEmpty
	}
	fsys, err := s.GetFsys()
	if err != nil {
		return nil, err
	}
	logger.V(2).Infof("object: %s", obj)

	// generate metadata
	m, err := metadata.GenerateFromFile(f)
	if err != nil {
		return nil, err
	}
	logger.V(2).Infof("%s has a checksum of %q", obj, m.Checksum)
	// convert metadata to json
	blob, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return nil, err
	}
	// Write metadata to disk
	metadataPath := fmt.Sprintf("%s.%s", obj, opts.MetadataFileExtension)
	logger.V(2).Infof("writing metadata to %s", metadataPath)
	if err := afero.WriteFile(fsys, metadataPath, blob, 0644); err != nil {
		return nil, err
	}

	cleanup = func() error {
		return fsys.Remove(metadataPath)
	}

	return cleanup, nil
}
