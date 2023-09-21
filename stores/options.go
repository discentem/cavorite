package stores

import "fmt"

type Options struct {
	BackendAddress        string `json:"backend_address" mapstructure:"backend_address"`
	PluginAddress         string `json:"plugin_address,omitempty" mapstructure:"plugin_address"`
	MetadataFileExtension string `json:"metadata_file_extension" mapstructure:"metadata_file_extension"`
	Region                string `json:"region" mapstructure:"region"`
	/*
		If ObjectKeyPrefix is set to "team-bucket", and the initialized backend supports it,
			- `cavorite upload whatever/thing` will be written to `team-bucket/whatever/thing`
			- `cavorite retrieve whatever/thing` will request `team-bucket/whatever/thing`
	*/
	ObjectKeyPrefix string `json:"object_key_prefix" mapstructure:"object_key_prefix"`
}

var ErrMetadataFileExtensionEmpty = fmt.Errorf("options.MetadatafileExtension cannot be %q", "")
