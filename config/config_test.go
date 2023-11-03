package config

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/discentem/cavorite/stores"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateFuncNil(t *testing.T) {
	conf := Config{
		Validate: nil,
	}
	err := conf.Write(afero.NewMemMapFs(), "")
	assert.ErrorIs(t, err, ErrValidateNil)
}
func TestValidateFails(t *testing.T) {
	conf := Config{
		Validate: func() error {
			return errors.New("failed")
		},
	}
	err := conf.Write(afero.NewMemMapFs(), "")
	assert.ErrorIs(t, err, ErrValidate)
}

func TestDirExpanderNil(t *testing.T) {
	conf := Config{
		Validate: func() error { return nil },
		Expander: nil,
	}
	err := conf.Write(afero.NewMemMapFs(), "")
	assert.ErrorIs(t, err, ErrDirExpanderNil)
}

func TestDirExpanderFails(t *testing.T) {
	conf := Config{
		Validate: func() error { return nil },
		Expander: func(path string) (string, error) {
			return "", errors.New("borked")
		},
	}
	err := conf.Write(afero.NewMemMapFs(), "")
	assert.ErrorIs(t, err, ErrDirExpander)
}

type MemFsBrokenStat struct {
	afero.Fs
}

func (m *MemFsBrokenStat) Stat(name string) (os.FileInfo, error) {
	return nil, errors.New("broken")
}

func TestWriteStatFails(t *testing.T) {
	conf := Config{
		Validate: func() error { return nil },
		Expander: func(path string) (string, error) {
			return path, nil
		},
	}
	err := conf.Write(&MemFsBrokenStat{Fs: afero.NewMemMapFs()}, "")
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "does not exist"))
}

type MemFsBrokenMkdirAll struct {
	afero.Fs
}

func (m *MemFsBrokenMkdirAll) MkdirAll(string, fs.FileMode) error {
	return errors.New("broken")
}

func TestWriteMkdirAllFails(t *testing.T) {
	conf := Config{
		Validate: func() error { return nil },
		Expander: func(path string) (string, error) {
			return path, nil
		},
	}
	err := conf.Write(&MemFsBrokenMkdirAll{Fs: afero.NewMemMapFs()}, "")
	require.Error(t, err)
}

func TestSuccessfulWrite(t *testing.T) {
	cfg := Config{
		StoreType: stores.StoreTypeS3,
		Options: stores.Options{
			BackendAddress:        "s3://blahaddress/bucket",
			MetadataFileExtension: "cfile",
			Region:                "us-east-9876",
		},
		Validate: func() error { return nil },
		// this is the emulates normal homedir.Expand behavior
		Expander: func(path string) (string, error) {
			return path, nil
		},
	}

	fsys := afero.NewMemMapFs()
	err := cfg.Write(fsys, ".")
	assert.NoError(t, err)
	f, err := fsys.Open(".cavorite/config")
	assert.NoError(t, err)
	b, err := afero.ReadAll(f)
	assert.NoError(t, err)

	expected := `{
 "store_type": "s3",
 "options": {
  "backend_address": "s3://blahaddress/bucket",
  "metadata_file_extension": "cfile",
  "region": "us-east-9876"
 }
}`
	assert.Equal(t, expected, string(b))

}

func TestWrite(t *testing.T) {
	t.Run("fail if validate() nil", TestValidateFuncNil)
	t.Run("fail if validate() fails", TestValidateFails)
	t.Run("fail if dirExpander() fails",
		TestDirExpanderFails,
	)
	t.Run("successful write", TestSuccessfulWrite)
}

func TestInitializeStoreTypeOf(t *testing.T) {
	ctx := context.Background()
	fsys := afero.NewMemMapFs()

	opts := stores.Options{
		BackendAddress:        "my-test-bucket",
		MetadataFileExtension: "cfile",
		Region:                "us-east-9876",
	}

	cfg := InitalizeStoreTypeOf(
		ctx,
		stores.StoreTypeGCS,
		fsys,
		"~/some_repo_root",
		opts,
	)

	// Assert the S3Store Config matches all of the inputs
	assert.Equal(t, cfg.StoreType, stores.StoreTypeGCS)
	assert.Equal(t, cfg.Options.BackendAddress, opts.BackendAddress)
	assert.Equal(t, cfg.Options.MetadataFileExtension, opts.MetadataFileExtension)
	assert.Equal(t, cfg.Options.Region, opts.Region)

	assert.NoError(t, cfg.Validate())
}
