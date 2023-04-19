package stores

import (
	"context"
	"errors"
)

type Store interface {
	WriteConfig(ctx context.Context, sourceRepo string) error
	Upload(ctx context.Context, sourceRepo string, destination string, objects ...string) error
	Retrieve(ctx context.Context, sourceRepo string, objects ...string) error
}
type Options struct {
	PantriAddress         string `json:"pantri_address" mapstructure:"pantri_address"`
	MetaDataFileExtension string `json:"metadata_file_extension" mapstructure:"metadata_file_extension"`
	// TODO(discentem) remove this option. See #15
	RemoveFromSourceRepo *bool `json:"remove_from_sourcerepo" mapstructure:"remove_from_sourcerepo"`
}

var (
	ErrRetrieveFailureHashMismatch = errors.New("hashes don't match, Retrieve aborted")
)
