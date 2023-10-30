package config

import (
	"testing"

	"github.com/discentem/cavorite/stores"
	"github.com/discentem/cavorite/testutils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig creates a pantri config file in memory
// to be read and parsed by viper
// Test inspired from https://github.com/spf13/viper/blob/master/viper_test.go
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name              string
		parseInto         *Config
		fsys              afero.Fs
		expected          *Config
		expectedLoadError func(e error) bool
	}{
		{
			name:      "valid config is parsed correctly",
			parseInto: &Config{},
			// creating an fsys with a valid cavorite config
			fsys: func(t *testing.T) afero.Fs {
				fs := afero.NewMemMapFs()

				err := fs.Mkdir(".cavorite", 0o777)
				require.NoError(t, err)

				f, err := fs.Create(testutils.AbsFilePath(t, ".cavorite/config"))
				require.NoError(t, err)

				_, err = f.Write([]byte(`{
					"store_type": "s3",
					"options": {
					 "backend_address": "s3://blahaddress/bucket",
					 "metadata_file_extension": "",
					 "region": "us-east-9876"
					}
				   }`),
				)
				require.NoError(t, err)
				f.Close()
				return fs
			}(t),
			expected: &Config{
				StoreType: stores.StoreType("s3"),
				Options: stores.Options{
					BackendAddress:        "s3://blahaddress/bucket",
					MetadataFileExtension: "",
					Region:                "us-east-9876",
				},
			},
			expectedLoadError: func(e error) bool {
				return e == nil
			},
		},
		{
			name:      "invalid config parsing returns error",
			parseInto: nil,
			fsys: func(t *testing.T) afero.Fs {
				fs := afero.NewMemMapFs()

				err := fs.Mkdir(".cavorite", 0o777)
				require.NoError(t, err)

				f, err := fs.Create(testutils.AbsFilePath(t, ".cavorite/config"))
				require.NoError(t, err)

				// "options" is missing the closing curly bracket below
				// so this is expected to cause ErrViperReadConfig
				_, err = f.Write([]byte(`{
					"store_type": "s3",
					"options": {
					 "backend_address": "s3://blahaddress/bucket",
					 "metadata_file_extension": "",
					 "region": "us-east-9876"
				   }`),
				)
				require.NoError(t, err)
				f.Close()
				return fs
			}(t),
			expectedLoadError: func(e error) bool {
				return assert.ErrorIs(t, e, ErrViperReadConfig)
			},
		},
		{
			name: "no cfgDest passed, parse into global cfg variable",
			parseInto: func() *Config {
				return &Cfg
			}(),
			fsys: func(t *testing.T) afero.Fs {
				fs := afero.NewMemMapFs()

				err := fs.Mkdir(".cavorite", 0o777)
				require.NoError(t, err)

				f, err := fs.Create(testutils.AbsFilePath(t, ".cavorite/config"))
				require.NoError(t, err)
				_, err = f.Write([]byte(`{
					"store_type": "s3",
					"options": {
					 "backend_address": "s3://blahaddress/bucket",
					 "metadata_file_extension": "",
					 "region": "us-east-9876"
					}
				   }`),
				)
				require.NoError(t, err)
				f.Close()
				return fs
			}(t),
			expectedLoadError: func(e error) bool {
				return e == nil
			},
			expected: &Config{
				StoreType: stores.StoreType("s3"),
				Options: stores.Options{
					BackendAddress:        "s3://blahaddress/bucket",
					MetadataFileExtension: "",
					Region:                "us-east-9876",
				},
			},
		},
	}

	t.Parallel()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Load(test.fsys, test.parseInto)
			assert.Equal(t, test.expectedLoadError(err), true)
			assert.Equal(t, test.expected, test.parseInto)
		})

	}
}
