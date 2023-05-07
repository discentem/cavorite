package pantri

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/logger"
	"github.com/spf13/afero"
)

func removePathPrefix(fsys afero.Fs, objects []string) ([]string, error) {
	logger.V(2).Infof("objects before removePathPrefix: %v", objects)
	// Our current path is the prefix to remove as pantri can only be run from the root
	// of the repo
	pathPrefix, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	pathPrefix = filepath.Dir(pathPrefix)
	logger.V(2).Infof("pathPrefix: %s", pathPrefix)
	newObjects := []string{}

	for _, object := range objects {
		newObjects = append(newObjects, strings.TrimPrefix(object, fmt.Sprintf("%s/", pathPrefix)))
	}
	logger.V(2).Infof("objects after removePathPrefix: %v", newObjects)

	return newObjects, nil
}
