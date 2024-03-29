package cli

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/discentem/cavorite/metadata"
	"github.com/discentem/cavorite/stores"
	"github.com/discentem/cavorite/testutils"
	"github.com/gonuts/go-shellquote"
	"github.com/google/logger"
	multierr "github.com/hashicorp/go-multierror"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// bazel test //internal/cli_test fails due to TestRetrieveAlreadyHaveAllFilesLocally unless we initialize the logger here
var (
	w = bytes.NewBufferString("")
)

func init() {
	logger.Init("test", true, false, w)
}
func TestRetrieveCmd(t *testing.T) {
	expectedRetrieveCmdArgs := "retrieve ./test_file_one ./test_file_two"

	retrieveCmd := retrieveCmd()

	// Split the args and handle bash escape characters
	args, err := shellquote.Split(expectedRetrieveCmdArgs)
	require.NoError(t, err)

	// Traverse splits the beginning command and args into separate parts
	subCmd, subArgs, err := retrieveCmd.Traverse(args)
	require.NoError(t, err)
	assert.NotNil(t, subCmd)
	assert.Equal(t, subCmd.UseLine(), "retrieve")

	// Test the the subArgs equal the expected expectedRetrieveCmdArgs and flags
	assert.NoError(t, subCmd.ParseFlags(subArgs))

	// Test that subArgs expects the same args above
	assert.Equal(t, subArgs, []string{"retrieve", "./test_file_one", "./test_file_two"})
}

type simpleStoreForRetrieve struct {
	// sourceFsys acts as the local file system where objects will be open for retrieving
	sourceFsys afero.Fs
	// bucketFsys acts as remote artifact storage, where objects will be downloaded from
	bucketFsys afero.Fs
	options    stores.Options
}

func (s simpleStoreForRetrieve) Upload(ctx context.Context, objects ...string) error {
	return nil
}
func (s simpleStoreForRetrieve) Retrieve(ctx context.Context, mmap metadata.CfileMetadataMap, cfiles ...string) error {
	var result error
	for _, c := range cfiles {
		m, ok := mmap[c]
		if !ok {
			result = multierr.Append(result, fmt.Errorf("%q not found in mmap", c))
			continue
		}
		object := m.Name
		objectHandle, err := s.bucketFsys.Open(object)
		if err != nil {
			e := fmt.Errorf("could not find %q (a cfile) in bucket: %w", c, err)
			result = multierr.Append(result, e)
			continue
		}
		objInfo, err := objectHandle.Stat()
		if err != nil {
			result = multierr.Append(result, err)
			continue
		}
		b := make([]byte, objInfo.Size())
		_, err = objectHandle.Read(b)
		if err != nil {
			e := fmt.Errorf("failed to read bytes from objectHandle: %w", err)
			result = multierr.Append(result, e)
			continue
		}
		sourceRepoObjectPath := strings.TrimSuffix(c, filepath.Ext(c))
		err = afero.WriteFile(s.sourceFsys, sourceRepoObjectPath, b, fs.ModeTemporary)
		if err != nil {
			result = multierr.Append(result, err)
			continue
		}
		matches, err := metadata.HashFromCfileMatches(s.sourceFsys, c, m.Checksum)
		if err != nil {
			result = multierr.Append(result, err)
			continue
		}
		if !matches {
			if err := s.sourceFsys.Remove(sourceRepoObjectPath); err != nil {
				result = multierr.Append(result, err)
				continue
			}
			result = multierr.Append(result, metadata.ErrRetrieveFailureHashMismatch)
			continue
		}
	}
	return result
}
func (s simpleStoreForRetrieve) GetOptions() (stores.Options, error) {
	return s.options, nil
}

func (s simpleStoreForRetrieve) Close() error {
	return nil
}

func TestRetrieve(t *testing.T) {
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	sourceFsys, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"someFile.cfile": {
			ModTime: &mTime,
			Content: []byte(`{
  "name": "repo/someFile",
  "checksum": "4df3c3f68fcc83b27e9d42c90431a72499f17875c81a599b566c9889b9696703",
  "date_modified": "2014-11-12T11:45:26.371Z"
}`),
		},
		"someOtherFile.cfile": {
			ModTime: &mTime,
			Content: []byte(`{
  "name": "repo/someOtherFile",
  "checksum": "4df3c3f68fcc83b27e9d42c90431a72499f17875c81a599b566c9889b9696703",
  "date_modified": "2014-11-12T11:45:26.371Z"
}`),
		},
	})
	assert.NoError(t, err)

	bucket, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"repo/someFile": {
			Content: []byte(`bla`),
			ModTime: &mTime,
		},
		"repo/someOtherFile": {
			Content: []byte(`bla`),
			ModTime: &mTime,
		},
	})
	assert.NoError(t, err)
	err = Retrieve(
		context.Background(),
		*sourceFsys,
		simpleStoreForRetrieve{
			sourceFsys: *sourceFsys,
			bucketFsys: *bucket,
			options: stores.Options{
				MetadataFileExtension: "cfile",
				BackendAddress:        "simpleStore/Test",
				ObjectKeyPrefix:       "repo",
			},
		},
		"someFile.cfile",
		"someOtherFile.cfile",
	)
	// TODO(discentem):
	assert.NoError(t, err)
}

func TestRetrieveAlreadyHaveAllFilesLocally(t *testing.T) {
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	sourceFsys, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"someFile.cfile": {
			ModTime: &mTime,
			Content: []byte(`{
  "name": "repo/someFile",
  "checksum": "4df3c3f68fcc83b27e9d42c90431a72499f17875c81a599b566c9889b9696703",
  "date_modified": "2014-11-12T11:45:26.371Z"
}`),
		},
		"someFile": {
			Content: []byte(`bla`),
			ModTime: &mTime,
		},
	})
	require.NoError(t, err)

	bucket, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"repo/someFile": {
			Content: []byte(`bla`),
			ModTime: &mTime,
		},
	})
	require.NoError(t, err)
	err = Retrieve(
		context.Background(),
		*sourceFsys,
		simpleStoreForRetrieve{
			sourceFsys: *sourceFsys,
			bucketFsys: *bucket,
			options: stores.Options{
				MetadataFileExtension: "cfile",
			},
		},
		"someFile.cfile",
	)
	require.NoError(t, err)
	require.True(t, strings.Contains(w.String(), "retrieval not needed, all requested files are present from"))
}
func TestShouldRetrieve(t *testing.T) {
	tests := []struct {
		name           string
		fsys           afero.Fs
		m              *metadata.ObjectMetaData
		cfile          string
		shouldRetrieve bool
		expectedError  func(e error) bool
	}{
		{
			name:           "metadata.ObjectMetadata is nil",
			fsys:           nil,
			m:              nil,
			cfile:          "",
			shouldRetrieve: true,
			expectedError: func(err error) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "object is not found locally, should retrieve",
			fsys: afero.NewMemMapFs(),
			m: &metadata.ObjectMetaData{
				Checksum: "blah",
			},
			cfile:          "blah",
			shouldRetrieve: true,
			expectedError: func(err error) bool {
				return strings.Contains(err.Error(), "file does not exist")

			},
		},
		{
			name: "object found locally, wrong hash, should retrieve",
			fsys: func(t *testing.T) afero.Fs {
				memfs := afero.NewMemMapFs()
				_, err := memfs.Create("blah")
				assert.NoError(t, err)
				return memfs
			}(t),
			m: &metadata.ObjectMetaData{
				Checksum: "blah",
			},
			cfile:          "blah.cfile",
			shouldRetrieve: true,
			expectedError: func(err error) bool {
				return err == nil
			},
		},
		{
			name: "object found locally, correct hash, no need to retrieve",
			fsys: func(t *testing.T) afero.Fs {
				memfs := afero.NewMemMapFs()
				_, err := memfs.Create("blah")
				assert.NoError(t, err)
				return memfs
			}(t),
			m: &metadata.ObjectMetaData{
				Checksum: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
			cfile:          "blah.cfile",
			shouldRetrieve: false,
			expectedError: func(err error) bool {
				return err == nil
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := shouldRetrieve(test.fsys, test.m, test.cfile)
			expectedErr := test.expectedError(err)
			assert.Equal(t, true, expectedErr)
			assert.Equal(t, test.shouldRetrieve, actual)
		})

	}
}
