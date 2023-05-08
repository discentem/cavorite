package root

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
