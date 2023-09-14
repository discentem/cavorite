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
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(".cavorite", 0o777)
	require.NoError(t, err)

	file, err := fs.Create(testutils.AbsFilePath(t, ".cavorite/config"))
	require.NoError(t, err)

	cfg := Config{}

	_, err = file.Write([]byte(`{
		"store_type": "s3",
		"options": {
		 "backend_address": "s3://blahaddress/bucket",
		 "metadata_file_extension": "",
		 "region": "us-east-9876"
		}
	   }`),
	)
	require.NoError(t, err)
	file.Close()

	err = Load(fs, &cfg)
	assert.NoError(t, err)

	assert.Equal(t, stores.StoreType("s3"), cfg.StoreType)
	assert.Equal(t, "s3://blahaddress/bucket", cfg.Options.BackendAddress)
	assert.Equal(t, "", cfg.Options.MetadataFileExtension)
	assert.Equal(t, "us-east-9876", cfg.Options.Region)

}
