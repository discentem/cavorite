package stores

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

type StoreType string

const (
	StoreTypeUndefined StoreType = "undefined"
	StoreTypeS3        StoreType = "s3"
	StoreTypeGCS       StoreType = "gcs"
)

var (
	ErrRetrieveFailureHashMismatch = errors.New("hashes don't match, Retrieve aborted")
)

type Store interface {
	Upload(ctx context.Context, objects ...string) error
	Retrieve(ctx context.Context, objects ...string) error
	GetOptions() Options
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
