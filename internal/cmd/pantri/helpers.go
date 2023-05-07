package pantri

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/logger"
)

func rootOfSourceRepo() (*string, error) {
	absPathOfConfig, err := filepath.Abs(".pantri/config")
	if err != nil {
		return nil, errors.New(".pantri/config not detected, not in sourceRepo root")
	}
	logger.V(2).Infof("absPathOfconfig: %q", absPathOfConfig)
	root := filepath.Dir(filepath.Dir(absPathOfConfig))
	return &root, nil
}

func removePathPrefix(objects []string, prefix string) ([]string, error) {
	for i, object := range objects {
		if !strings.HasPrefix(object, prefix) {
			return nil, fmt.Errorf("%q does not exist relative to source_repo: %q", object, prefix)
		}
		objects[i] = strings.TrimPrefix(object, fmt.Sprintf("%s/", prefix))
	}

	return objects, nil
}
