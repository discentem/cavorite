package cli

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/logger"
	"github.com/spf13/afero"

	"github.com/discentem/cavorite/internal/config"
	"github.com/discentem/cavorite/internal/stores"
)

func rootOfSourceRepo() (*string, error) {
	absPathOfConfig, err := filepath.Abs(".cavorite/config")
	if err != nil {
		return nil, errors.New(".cavorite/config not detected, not in sourceRepo root")
	}
	logger.V(2).Infof("absPathOfconfig: %q", absPathOfConfig)
	root := filepath.Dir(filepath.Dir(absPathOfConfig))
	return &root, nil
}

func removePathPrefix(objects []string, prefix string) ([]string, error) {
	for i, object := range objects {
		absObject, err := filepath.Abs(object)
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(absObject, prefix) {
			return nil, fmt.Errorf("%q does not exist relative to source_repo: %q", object, prefix)
		}
		objects[i] = strings.TrimPrefix(absObject, fmt.Sprintf("%s/", prefix))
	}

	return objects, nil
}

func initStoreFromConfig(ctx context.Context, cfg config.Config, fsys afero.Fs, opts stores.Options) (stores.Store, error) {
	var s stores.Store
	switch cfg.StoreType {
	case stores.StoreTypeS3:
		s3, err := stores.NewS3StoreClient(ctx, fsys, opts)
		if err != nil {
			return nil, fmt.Errorf("improper stores.S3Client init: %v", err)
		}
		s = stores.Store(s3)
	case stores.StoreTypeGCS:
		gcs, err := stores.NewGCSStoreClient(ctx, fsys, opts)
		if err != nil {
			return nil, fmt.Errorf("improper stores.GCSClient init: %v", err)
		}
		s = stores.Store(gcs)
	case stores.StoreTypeAzureBlob:
		az, err := stores.NewAzureBlobStore(
			ctx,
			fsys,
			opts,
			azblob.ClientOptions{},
		)
		if err != nil {
			return nil, fmt.Errorf("improper stores.AzureBlobStore init: %v", err)
		}
		s = stores.Store(az)
	case stores.StoreTypeGoPlugin:
		// FIXME: allow specifying command arguments for plugin
		ps, err := stores.NewPluggableStore(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("improper plugin init: %v", err)
		}
		return ps, nil
	default:
		return nil, fmt.Errorf("type %s is not supported", cfg.StoreType)
	}

	return s, nil
}

func setLoggerOpts() {
	if vv {
		logger.SetLevel(2)
	}
	logger.SetFlags(log.LUTC)
}
