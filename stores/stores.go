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
	PantriAddress string `json:"pantri_address" mapstructure:"pantri_address"`
	MetaDataFileExtension string `json:"metadata_file_extension" mapstructure:"metadata_file_extension"`
	// TODO(discentem) remove this option. See #15
	RemoveFromSourceRepo *bool `json:"remove_from_sourcerepo" mapstructure:"remove_from_sourcerepo"`
}

var (
	ErrRetrieveFailureHashMismatch = errors.New("hashes don't match, Retrieve aborted")
)

func NewS3Store(sourceRepo, backend, address string, opts Options) *S3Store {
	return &S3Store{}
}

func InitLocalStore(fsys afero.Fs, sourceRepo string) *LocalStore {
	epa, err := homedir.Expand(s.PantriAddress)
	if err != nil {
		return err
	}

	c := pantriconfig.Config{
		Type:          "local",
		PantriAddress: epa,
		Opts:          s.Opts,
		Validate: func() error {
			// Ensure s.PantriAddress exists before writing config to disk
			if _, err := os.Stat(s.PantriAddress); err != nil {
				fmt.Println(err)
				return fmt.Errorf("specified pantri_address %q does not exist, so we can't make it a pantri repo", epa)
			}
			return nil
		},
	}

	return c.Write(fsys, sourceRepo)
	// return &LocalStore{
	// 	PantriAddress: 
	// 	Opts: opts,
	// }
}

func LocalStoreClient(fsys afero.Fs, sourceRepo string, opts Options) *LocalStore {

}

