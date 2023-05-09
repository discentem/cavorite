package cli

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/discentem/pantri_but_go/internal/config"
	"github.com/discentem/pantri_but_go/internal/metadata"
	"github.com/discentem/pantri_but_go/internal/stores"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_removePathPrefix(t *testing.T) {
	pathPrefix, err := os.Getwd()
	assert.NoError(t, err)

	expectedRemovePathPrefixes := []string{"foo/foo.dmg", "bar/bar.pkg"}
	testRemovePathPrefixes, err := removePathPrefix([]string{
		fmt.Sprintf("%s/%s", pathPrefix, "foo/foo.dmg"),
		fmt.Sprintf("%s/%s", pathPrefix, "bar/bar.pkg"),
	})
	assert.NoError(t, err)
	assert.Equal(t, expectedRemovePathPrefixes, testRemovePathPrefixes)
}

func Test_buildStoresFromConfig(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		StoreType: stores.StoreTypeS3,
		Options: stores.Options{
			PantriAddress:         "s3://test-bucket",
			MetaDataFileExtension: metadata.MetaDataFileExtension,
			Region:                "us-east-9876",
		},
	}

	fsys := afero.NewMemMapFs()

	s, err := buildStoresFromConfig(
		ctx,
		cfg,
		fsys,
		cfg.Options,
	)
	assert.NoError(t, err)

	// Test GetOptions from newly created S3Store
	assert.Equal(t, s.GetOptions().PantriAddress, "s3://test-bucket")

	// Test if type is equal to stores.S3Store
	// to ensure sure the buildStores returns
	// the proper stores based on type
	assert.Equal(
		t,
		reflect.TypeOf(stores.Store(&stores.S3Store{})).Elem(),
		reflect.TypeOf(s).Elem(),
	)
}

func Test_buildStoresFromConfig_Improper_S3Client(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		StoreType: stores.StoreTypeS3,
		Options: stores.Options{
			PantriAddress:         "s4://test-bucket",
			MetaDataFileExtension: metadata.MetaDataFileExtension,
			Region:                "us-east-9876",
		},
	}

	fsys := afero.NewMemMapFs()

	_, err := buildStoresFromConfig(
		ctx,
		cfg,
		fsys,
		cfg.Options,
	)

	t.Log(err.Error())

	assert.ErrorContains(t,
		err,
		`improper stores.S3Client init: pantriAddress did not contain s3://, http://, or https:// prefix`,
	)
}

func Test_buildStoresFromConfig_ImproperStoreType(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		StoreType: "s4",
		Options: stores.Options{
			PantriAddress:         "s4://test-bucket",
			MetaDataFileExtension: metadata.MetaDataFileExtension,
			Region:                "us-east-9876",
		},
	}

	fsys := afero.NewMemMapFs()

	_, err := buildStoresFromConfig(
		ctx,
		cfg,
		fsys,
		cfg.Options,
	)

	t.Log(err.Error())

	assert.ErrorContains(t,
		err,
		fmt.Sprintf("type %s is not supported", "s4"),
	)
}
