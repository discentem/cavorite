package cli

import (
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
	multierr "github.com/hashicorp/go-multierror"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

			// 44e77fe9966d66d9dbc24fb67c2be59f0d7681c541b8221bd6417def2ac561c4
			Content: []byte(`{
  "name": "repo/someFile",
  "checksum": "4df3c3f68fcc83b27e9d42c90431a72499f17875c81a599b566c9889b9696703",
  "date_modified": "2014-11-12T11:45:26.371Z"
}`),
		},
		"someOtherFile.cfile": {
			ModTime: &mTime,
			//4937c76e9ef850fbedd79cad64b8bed25d261a334b41303a3743f8ab7e4af2db
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
	assert.NoError(t, err)
}
