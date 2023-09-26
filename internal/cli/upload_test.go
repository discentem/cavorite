package cli

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/discentem/cavorite/metadata"
	"github.com/discentem/cavorite/stores"
	"github.com/discentem/cavorite/testutils"
	"github.com/gonuts/go-shellquote"
	"github.com/google/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadCmd(t *testing.T) {
	expectedUploadCmdArgs := "upload ./test_file_one ./test_file_two"

	uploadCmd := uploadCmd()

	// Split the args and handle bash escape characters
	args, err := shellquote.Split(expectedUploadCmdArgs)
	require.NoError(t, err)

	// Traverse splits the beginning command and args into separate parts
	subCmd, subArgs, err := uploadCmd.Traverse(args)
	require.NoError(t, err)
	assert.NotNil(t, subCmd)
	assert.Equal(t, subCmd.UseLine(), "upload")

	// Test the the subArgs equal the expected expectedUploadCmdArgs and flags
	assert.NoError(t, subCmd.ParseFlags(subArgs))

	// Test that subArgs expects the same args above
	assert.Equal(t, subArgs, []string{"upload", "./test_file_one", "./test_file_two"})
}

var (
	_ = stores.Store(simpleStore{})
)

type simpleStore struct {
	// sourceFsys acts as the local file system where objects will be open for uploading
	sourceFsys afero.Fs
	// bucketFsys acts as remote artifact storage, where objects will be uploaded
	bucketFsys afero.Fs
	options    stores.Options
}

func (s simpleStore) Upload(ctx context.Context, objects ...string) error {
	return nil
}
func (s simpleStore) Retrieve(ctx context.Context, mmap metadata.CfileMetadataMap, objects ...string) error {
	return nil
}
func (s simpleStore) GetOptions() (stores.Options, error) {
	return s.options, nil
}

func (s simpleStore) Close() error {
	return nil
}

// TestUpload tests whether metadata gets generated correctly
func TestUpload(t *testing.T) {
	logger.Init("TestUpload", false, false, io.Discard)
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	objs := []string{"someFile", "someOtherFile"}
	sourceFsys, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		objs[0]: {
			ModTime: &mTime,
			Content: []byte(`stuff`),
		},
		objs[1]: {
			ModTime: &mTime,
			Content: []byte(`stuff`),
		},
	})
	assert.NoError(t, err)
	bucket, err := testutils.MemMapFsWith(map[string]testutils.MapFile{})
	assert.NoError(t, err)
	sStore := simpleStore{
		sourceFsys: *sourceFsys,
		bucketFsys: *bucket,
		options: stores.Options{
			MetadataFileExtension: "cfile",
			BackendAddress:        "simpleStore/Test",
		},
	}
	err = upload(context.Background(), *sourceFsys, sStore, objs...)
	assert.NoError(t, err)

	require.NoError(t, err)

	sopts, err := sStore.GetOptions()
	assert.NoError(t, err)
	for _, f := range objs {
		b, err := afero.ReadFile(sStore.sourceFsys, fmt.Sprintf("%s.%s", f, sopts.MetadataFileExtension))
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(`{
 "name": "%s",
 "checksum": "35bafb1ce99aef3ab068afbaabae8f21fd9b9f02d3a9442e364fa92c0b3eeef0",
 "date_modified": "2014-11-12T11:45:26.371Z"
}`, f), string(b))
	}

	assert.NoError(t, err)
}

// TestUploadWithPrefix tests whether metadata gets generated correctly
func TestUploadWithPrefix(t *testing.T) {
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	sourceFsys, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"someFile": {
			ModTime: &mTime,
			Content: []byte(`stuff`),
		},
		"someOtherFile": {
			ModTime: &mTime,
			Content: []byte(`stuff`),
		},
	})
	assert.NoError(t, err)
	bucket, err := testutils.MemMapFsWith(map[string]testutils.MapFile{})
	assert.NoError(t, err)
	sStore := simpleStore{
		sourceFsys: *sourceFsys,
		bucketFsys: *bucket,
		options: stores.Options{
			MetadataFileExtension: "cfile",
			BackendAddress:        "simpleStore/Test",
			ObjectKeyPrefix:       "aCoolPrefix",
		},
	}
	err = upload(context.Background(), *sourceFsys, sStore, "someFile", "someOtherFile")
	require.NoError(t, err)

	require.NoError(t, err)

	sopts, err := sStore.GetOptions()
	assert.NoError(t, err)
	for _, f := range []string{"someFile", "someOtherFile"} {
		b, err := afero.ReadFile(*sourceFsys, fmt.Sprintf("%s.%s", f, sopts.MetadataFileExtension))
		require.NoError(t, err)
		var objkey string
		if sopts.ObjectKeyPrefix != "" {
			objkey = fmt.Sprintf("%s/%s", sopts.ObjectKeyPrefix, f)
		} else {
			objkey = f
		}
		expected := fmt.Sprintf(`{
 "name": "%s",
 "checksum": "35bafb1ce99aef3ab068afbaabae8f21fd9b9f02d3a9442e364fa92c0b3eeef0",
 "date_modified": "2014-11-12T11:45:26.371Z"
}`, objkey)
		assert.Equal(t, expected, string(b))
	}
}

// TestUploadPartialFail tests whether metadata generation will succeed for n+1 even if n fails
func TestUploadPartialFail(t *testing.T) {
	logger.Init("TestUpload", false, false, io.Discard)
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	sourceFsys, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"someFile": {
			ModTime: &mTime,
			Content: []byte(`stuff`),
		},
	})
	assert.NoError(t, err)
	bucket, err := testutils.MemMapFsWith(map[string]testutils.MapFile{})
	assert.NoError(t, err)
	sStore := simpleStore{
		sourceFsys: *sourceFsys,
		bucketFsys: *bucket,
	}
	err = upload(context.Background(), *sourceFsys, sStore, "someOtherFileThatDoesntExist", "someFile")

	// upload is expected for fail for someOtherFileThatDoesntExist as it does not exist in sourceFsys
	require.ErrorIs(t, err, ErrOpen)

	sopts, err := sStore.GetOptions()
	assert.NoError(t, err)
	for _, f := range []string{"someFile"} {
		b, _ := afero.ReadFile(sStore.sourceFsys, fmt.Sprintf("%s.%s", f, sopts.MetadataFileExtension))
		assert.Equal(t, fmt.Sprintf(`{
 "name": "%s",
 "checksum": "35bafb1ce99aef3ab068afbaabae8f21fd9b9f02d3a9442e364fa92c0b3eeef0",
 "date_modified": "2014-11-12T11:45:26.371Z"
}`, f), string(b))
	}

	assert.NoError(t, err)
}
