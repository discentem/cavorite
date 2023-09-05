package stores

import "fmt"

type Options struct {
	BackendAddress        string `json:"backend_address" mapstructure:"backend_address"`
	PluginAddress         string `json:"plugin_address" mapstructure:"plugin_address"`
	MetadataFileExtension string `json:"metadata_file_extension" mapstructure:"metadata_file_extension"`
	Region                string `json:"region" mapstructure:"region"`
}

var ErrMetadataFileExtensionEmpty = fmt.Errorf("options.MetadatafileExtension cannot be %q", "")
