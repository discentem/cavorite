package config

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestValidateFuncNil(t *testing.T) {
	conf := Config{
		Validate: nil,
	}
	err := conf.WriteToDisk(afero.NewMemMapFs(), "")
	assert.ErrorIs(t, err, ErrValidateNil)
}
func TestValidateFails(t *testing.T) {
	conf := Config{
		Validate: func() error {
			return errors.New("failed")
		},
	}
	err := conf.WriteToDisk(afero.NewMemMapFs(), "")
	assert.ErrorIs(t, err, ErrValidate)
}

func TestJsonMarshalIndenterFails(t *testing.T) {
	conf := Config{
		Validate: func() error { return nil },
	}
	// override marshaler with broken one
	jsonMarshalIndenter = func(v any, prefix string, indent string) ([]byte, error) {
		return nil, errors.New("borked")
	}
	err := conf.WriteToDisk(
		afero.NewMemMapFs(),
		"",
	)
	assert.ErrorIs(t, err, ErrJsonMarshal)
}

func TestDirExpanderFails(t *testing.T) {
	conf := Config{
		Validate: func() error { return nil },
	}
	dirExpander = func(path string) (string, error) {
		return "", errors.New("borked")
	}
	// TODO(discentem): this is shitty, one test should not break another
	// reset jsonMarshalIndenter back after previous test
	jsonMarshalIndenter = json.MarshalIndent
	err := conf.WriteToDisk(afero.NewMemMapFs(), "")
	assert.ErrorIs(t, err, ErrDirExpander)
}

func TestWriteToDisk(t *testing.T) {
	t.Run("fail if validate() nil", TestValidateFuncNil)
	t.Run("fail if validate() fails", TestValidateFails)
	t.Run("fail if jsonMarshalIndenter errors",
		TestJsonMarshalIndenterFails,
	)
	t.Run("fail if dirExpander errors",
		TestDirExpanderFails,
	)
}
