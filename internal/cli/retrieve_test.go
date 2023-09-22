package cli

import (
	"context"
	"fmt"
	"io/fs"
	"testing"
	"time"

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
func (s simpleStoreForRetrieve) Retrieve(ctx context.Context, objects ...string) error {
	var result error
	for _, o := range objects {
		objectHandle, err := s.bucketFsys.Open(o)
		if err != nil {
			e := fmt.Errorf("could not find %s in bucket: %w", o, err)
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
		err = afero.WriteFile(s.sourceFsys, o, b, fs.ModeTemporary)
		if err != nil {
			result = multierr.Append(result, err)
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
  "checksum": "35bafb1ce99aef3ab068afbaabae8f21fd9b9f02d3a9442e364fa92c0b3eeef0",
  "date_modified": "2014-11-12T11:45:26.371Z"
}`),
		},
		"someOtherFile.cfile": {
			ModTime: &mTime,
			Content: []byte(`{
  "name": "repo/someOtherFile",
  "checksum": "35bafb1ce99aef3ab068afbaabae8f21fd9b9f02d3a9442e364fa92c0b3eeef0",
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
	err = retrieve(
		context.Background(),
		*sourceFsys,
		simpleStoreForRetrieve{
			sourceFsys: *sourceFsys,
			bucketFsys: *bucket,
			options: stores.Options{
				MetadataFileExtension: "cfile",
				BackendAddress:        "simpleStore/Test",
			},
		},
		"someFile.cfile",
		"someOtherFile.cfile",
	)
	assert.NoError(t, err)
}
