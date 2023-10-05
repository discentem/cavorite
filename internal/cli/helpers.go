package cli

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/google/logger"
	"github.com/spf13/afero"

	"github.com/discentem/cavorite/config"
	"github.com/discentem/cavorite/stores"
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

func initStoreFromConfig(ctx context.Context, cfg config.Config, fsys afero.Fs) (stores.Store, error) {
	switch cfg.StoreType {
	case stores.StoreTypeS3:
		s3, err := stores.NewS3StoreClient(ctx, fsys, cfg.Options)
		if err != nil {
			return nil, fmt.Errorf("improper stores.S3Client init: %v", err)
		}
		return stores.Store(s3), nil
	case stores.StoreTypeGCS:
		return nil, fmt.Errorf("type %s is not currently supported; support for %q will be re-enabled in a future release", cfg.StoreType, cfg.StoreType)
	case stores.StoreTypeAzureBlob:
		return nil, fmt.Errorf("type %s is not currently supported; support for %q will be re-enabled in a future release", cfg.StoreType, cfg.StoreType)
	case stores.StoreTypeGoPlugin:
		// TODO(discentem): allow specifying command arguments for plugin
		ps, err := stores.NewPluggableStore(ctx, cfg.Options)
		if err != nil {
			return nil, fmt.Errorf("improper plugin init: %v", err)
		}
		return ps, nil
	default:
		return nil, fmt.Errorf("type %s is not supported", cfg.StoreType)
	}
}

func setLoggerOpts() {
	if VV {
		logger.SetLevel(2)
	}
	logger.SetFlags(log.LUTC)
}
