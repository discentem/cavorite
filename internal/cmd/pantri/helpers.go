package pantri

import (
	"fmt"
	"os"
	"strings"
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
