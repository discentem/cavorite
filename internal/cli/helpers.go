package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/discentem/pantri_but_go/internal/config"
	"github.com/discentem/pantri_but_go/internal/stores"
	"github.com/google/logger"
	"github.com/spf13/afero"
)

func removePathPrefix(objects []string) ([]string, error) {
	// Our current path is the prefix to remove as pantri can only be run from the root
	// of the repo
	pathPrefix, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for i, object := range objects {
		objects[i] = strings.TrimPrefix(object, fmt.Sprintf("%s/", pathPrefix))
	}

	return objects, nil
}

func buildStoresFromConfig(ctx context.Context, cfg config.Config, fsys afero.Fs, opts stores.Options) (stores.Store, error) {
	var s stores.Store
	switch cfg.StoreType {
	case stores.StoreTypeS3:
		s3, err := stores.NewS3StoreClient(ctx, fsys, opts)
		if err != nil {
			return nil, fmt.Errorf("improper stores.S3Client init: %v", err)
		}
		s = stores.Store(s3)
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
