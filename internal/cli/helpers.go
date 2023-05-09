package cli

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/discentem/pantri_but_go/internal/config"
	"github.com/discentem/pantri_but_go/internal/stores"
	"github.com/google/logger"
	"github.com/spf13/afero"
)

func removePathPrefix(pathPrefix string, objects []string) ([]string, error) {
	for i, object := range objects {
		objects[i] = strings.TrimPrefix(object, fmt.Sprintf("%s/", pathPrefix))
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
