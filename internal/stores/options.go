package stores

type Options struct {
	BackendAddress        string `json:"backend_address" mapstructure:"backend_address"`
	MetaDataFileExtension string `json:"metadata_file_extension" mapstructure:"metadata_file_extension"`
	Region                string `json:"region" mapstructure:"region"`
}
