package cli

import (
	"testing"

	"github.com/discentem/pantri_but_go/internal/testutil"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test inspired from https://github.com/spf13/viper/blob/master/viper_test.go
func Test_loadConfig(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(".pantri", 0o777)
	require.NoError(t, err)

	file, err := fs.Create(testutil.AbsFilePath(t, ".pantri/config"))
	require.NoError(t, err)

	_, err = file.Write([]byte(`{
		"store_type": "s3",
		"options": {
		 "pantri_address": "s3://blahaddress/bucket",
		 "metadata_file_extension": "",
		 "region": "us-east-9876"
		}
	   }`),
	)
	require.NoError(t, err)
	file.Close()

	err = loadConfig(fs)
	assert.NoError(t, err)
}
