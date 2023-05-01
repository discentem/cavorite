package config

import (
	"errors"
	"testing"

	"github.com/discentem/pantri_but_go/internal/stores"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
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

func TestDirExpanderFails(t *testing.T) {
	conf := Config{
		Validate: func() error { return nil },
	}
	expander = func(path string) (string, error) {
		return "", errors.New("borked")
	}
	err := conf.Write(afero.NewMemMapFs(), "")
	assert.ErrorIs(t, err, ErrDirExpander)
}

func TestSuccessfulWrite(t *testing.T) {
	conf := Config{
		StoreType: stores.StoreTypeUndefined,
		Options: stores.Options{
			PantriAddress: "s3://blahaddress/bucket",
		},
		Validate: func() error { return nil },
	}
	// override back to a dirExpander that will succeed, as opposed to previous test
	expander = func(path string) (string, error) {
		return path, nil
	}
	fsys := afero.NewMemMapFs()
	err := conf.Write(fsys, ".")
	assert.NoError(t, err)
	f, err := fsys.Open(".pantri/config")
	assert.NoError(t, err)
	b, err := afero.ReadAll(f)
	assert.NoError(t, err)

	expected := `{
 "store_type": "undefined",
 "options": {
  "pantri_address": "s3://blahaddress/bucket"
  "metadata_file_extension": "",
  "remove_from_sourcerepo": null
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
