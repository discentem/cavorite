package stores

import (
	"context"
	"errors"

	"github.com/spf13/afero"
)

type Store interface {
	Upload(ctx context.Context, fsys afero.Fs, sourceRepo string, destination string, objects ...string) error
	Retrieve(ctx context.Context, fsys afero.Fs, sourceRepo string, objects ...string) error
	Init(ctx context.Context)
}
type Options struct {
	MetaDataFileExtension string `json:"metadata_file_extension" mapstructure:"metadata_file_extension"`
	// TODO(discentem) remove this option. See #15
	RemoveFromSourceRepo *bool `json:"remove_from_sourcerepo" mapstructure:"remove_from_sourcerepo"`
}

var (
	ErrRetrieveFailureHashMismatch = errors.New("hashes don't match, Retrieve aborted")
)
