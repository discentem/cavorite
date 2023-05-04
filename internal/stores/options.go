package stores

type Options struct {
	PantriAddress         string `json:"pantri_address" mapstructure:"pantri_address"`
	MetaDataFileExtension string `json:"metadata_file_extension" mapstructure:"metadata_file_extension"`
	Region                string `json:"region" mapstructure:"region"`
}
