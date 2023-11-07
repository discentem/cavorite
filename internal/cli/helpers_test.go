package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/carolynvs/aferox"
	"github.com/discentem/cavorite/config"
	"github.com/discentem/cavorite/metadata"
	"github.com/discentem/cavorite/stores"
	"github.com/google/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemovePathPrefix(t *testing.T) {
	pathPrefix, err := os.Getwd()
	assert.NoError(t, err)

	expectedRemovePathPrefixes := []string{"foo/foo.dmg", "bar/bar.pkg"}
	testRemovePathPrefixes, err := removePathPrefix(
		[]string{
			fmt.Sprintf("%s/%s", pathPrefix, "foo/foo.dmg"),
			fmt.Sprintf("%s/%s", pathPrefix, "bar/bar.pkg"),
		},
		pathPrefix,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedRemovePathPrefixes, testRemovePathPrefixes)
}

func TestInitStoreFromConfig(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		StoreType: stores.StoreTypeS3,
		Options: stores.Options{
			BackendAddress:        "s3://test-bucket",
			MetadataFileExtension: metadata.MetadataFileExtension,
			Region:                "us-east-9876",
		},
	}

	fsys := afero.NewMemMapFs()

	s, err := initStoreFromConfig(
		ctx,
		cfg,
		fsys,
	)
	assert.NoError(t, err)

	opts, err := s.GetOptions()
	assert.NoError(t, err)

	// Test GetOptions from newly created S3Store
	assert.Equal(t, opts.BackendAddress, "s3://test-bucket")

	// Test if type is equal to stores.S3Store
	// to ensure sure the buildStores returns
	// the proper stores based on type
	// assert.Equal(
	// 	t,
	// 	reflect.TypeOf(stores.Store(&stores.S3Store{})).Elem(),
	// 	reflect.TypeOf(s).Elem(),
	// )
}

func TestInitStoreFromConfig_InvalidOptions(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		StoreType: stores.StoreTypeS3,
		Options: stores.Options{
			BackendAddress:        "s4://test-bucket",
			MetadataFileExtension: metadata.MetadataFileExtension,
			Region:                "us-east-9876",
		},
	}

	fsys := afero.NewMemMapFs()

	_, err := initStoreFromConfig(
		ctx,
		cfg,
		fsys,
	)

	assert.Errorf(t, err, `improper stores.S3Client init: cavoriteAddress did not contain s3://, http://, or https:// prefix`)
}

func TestInitStoreFromConfig_InvalidateStoreType(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		StoreType: "s4",
		Options: stores.Options{
			BackendAddress:        "s4://test-bucket",
			MetadataFileExtension: metadata.MetadataFileExtension,
			Region:                "us-east-9876",
		},
	}

	fsys := afero.NewMemMapFs()

	_, err := initStoreFromConfig(
		ctx,
		cfg,
		fsys,
	)

	assert.Errorf(t, err, "type %s is not supported", "s4")
}

type aferoxWithAbsErr struct {
	aferox.Aferox
}

func newAferoxWithAbsErr(root string, fs *afero.Fs) *aferoxWithAbsErr {
	return &aferoxWithAbsErr{
		Aferox: aferox.NewAferox(root, *fs),
	}
}

func (a *aferoxWithAbsErr) Abs(path string) (string, error) {
	_, err := a.Afero.Open(path)
	if err != nil {
		return "", err
	}
	return a.Aferox.Abs(path), nil
}

func TestRootOfSourceRepo(t *testing.T) {
	type test struct {
		name        string
		fsyses      []fsWithAbs
		expected    func(*string) bool
		expectedErr error
	}
	tests := []test{
		{
			name: "too many fsyses",
			fsyses: func() []fsWithAbs {
				memfs := afero.NewMemMapFs()
				return []fsWithAbs{
					newAferoxWithAbsErr("", &memfs),
					newAferoxWithAbsErr("", &memfs),
				}
			}(),
			expected: func(s *string) bool {
				return s == nil
			},
			expectedErr: ErrTooManyFsyses,
		},
		{
			name: "no .cavorite/config",
			fsyses: []fsWithAbs{newAferoxWithAbsErr("", func() *afero.Fs {
				memfs := afero.NewMemMapFs()
				return &memfs
			}())},
			expectedErr: errors.New(
				".cavorite/config not detected, not in sourceRepo root",
			),
			expected: func(path *string) bool {
				return path == nil
			},
		},
		{
			name: "with .cavorite/config",
			fsyses: func(t *testing.T) []fsWithAbs {
				memfs := afero.NewMemMapFs()
				err := memfs.MkdirAll("code_repo", 0755)
				require.NoError(t, err)
				f, err := memfs.Create("/code_repo/.cavorite/config")
				require.NoError(t, err)
				_, err = f.Write([]byte(`{
					"store_type": "s3",
					"options": {
						"backend_address": "s3://blahaddress/bucket",
						"metadata_file_extension": "",
						"region": "us-east-9876"
					}	
				}`))
				require.NoError(t, err)
				return []fsWithAbs{newAferoxWithAbsErr("/code_repo", &memfs)}
			}(t),
			expected: func(s *string) bool {
				return assert.Equal(t, "/code_repo", *s)
			},
			expectedErr: nil,
		},
	}
	logger.Init("helpers_test", true, false, os.Stdout)
	t.Parallel()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var root *string
			var err error
			if test.fsyses == nil {
				root, err = rootOfSourceRepo()
			} else {
				root, err = rootOfSourceRepo(test.fsyses...)
			}
			require.Equal(t, test.expectedErr, err)
			if test.expected == nil {
				t.Error("test.expected function must be provided")
				t.Fail()
			} else {
				require.Equal(t, test.expected(root), true)
			}

		})
	}
}
