package stores

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/discentem/cavorite/metadata"
	"github.com/discentem/cavorite/testutils"
	"github.com/google/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	// "github.com/discentem/cavorite/internal/cli"
)

type aferoS3Server struct {
	buckets map[string]afero.Fs
}

func (s aferoS3Server) Upload(ctx context.Context,
	input *s3.PutObjectInput,
	opts ...func(*s3manager.Uploader)) (
	*s3manager.UploadOutput, error,
) {
	bucket := *input.Bucket
	// check if the bucket referenced in input exists
	_, ok := s.buckets[bucket]
	if !ok {
		return nil, fmt.Errorf("%s does not exist in this aferoS3Server", bucket)
	}
	inputBytes, err := io.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}
	// create a filesystem for bucket referenced in input
	bucketfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		*input.Key: {
			// write input body to bucketfs
			Content: inputBytes,
		},
	})
	if err != nil {
		return nil, err
	}
	// write bucketfs to associated bucket
	s.buckets[bucket] = *bucketfs
	// S3Store doesn't use UploadOutput, so in the test we don't either
	return nil, nil
}

func (s aferoS3Server) Download(
	ctx context.Context,
	w io.WriterAt,
	input *s3.GetObjectInput,
	options ...func(*s3manager.Downloader)) (n int64, err error) {

	bucket := *input.Bucket
	// check if the bucket referenced in input exists
	_, ok := s.buckets[bucket]
	if !ok {
		return 0, fmt.Errorf("%s does not exist in this aferoS3Server", bucket)
	}

	objectHandle, err := s.buckets[bucket].Open(*input.Key)
	if err != nil {
		return 0, fmt.Errorf("could not find %s in bucket %s: %w", *input.Key, bucket, err)
	}
	objInfo, err := objectHandle.Stat()
	if err != nil {
		return 0, err
	}
	b := make([]byte, objInfo.Size())
	_, err = objectHandle.Read(b)
	if err != nil {
		return 0, fmt.Errorf("failed to read bytes from objectHandle: %w", err)
	}
	nbw, err := w.WriteAt(b, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to write objectHandle bytes to w: %w", err)
	}
	return int64(nbw), nil

}

func TestS3StoreUpload(t *testing.T) {
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	memfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"test": {
			Content: []byte("bla"),
			ModTime: &mTime,
		},
	})
	assert.NoError(t, err)

	fakeS3Server := aferoS3Server{
		buckets: map[string]afero.Fs{
			// create a bucket in our fake s3 server
			"test": afero.NewMemMapFs(),
		},
	}

	store := S3Store{
		Options: Options{
			BackendAddress:        "s3://test",
			MetadataFileExtension: "cfile",
		},
		fsys:         *memfs,
		awsRegion:    "us-east-1",
		s3Uploader:   fakeS3Server,
		s3Downloader: fakeS3Server,
	}

	err = store.Upload(context.Background(), "test")
	assert.NoError(t, err)

	buck := fakeS3Server.buckets["test"]
	f, err := buck.Open("test")
	assert.NoError(t, err)
	fstat, err := f.Stat()
	assert.NoError(t, err)
	b := make([]byte, fstat.Size())
	_, err = f.Read(b)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`bla`), b)
}

func TestS3StoreRetrieve(t *testing.T) {
	logger.Init("TestS3StoreRetrieve", true, true, os.Stdout)
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	// create bucket content
	bucketfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"someObject": {
			Content: []byte("tla"),
			ModTime: &mTime,
		},
	})
	require.NoError(t, err)

	fakeS3Server := aferoS3Server{
		buckets: map[string]afero.Fs{
			// create a bucket in our fake s3 server with the content
			"aFakeBucket": *bucketfs,
		},
	}
	som := `{
		"name": "someObject",
		"checksum": "59e5ad2a03d2499749f7943c9dded0f303ad7542befef6d0aead8a7888587f66",
		"date_modified": "2014-11-12T11:45:26.371Z"
	   }`

	localFs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"someObject.cfile": {
			Content: []byte(som),
			ModTime: &mTime,
		},
	})
	require.NoError(t, err)

	store := S3Store{
		Options: Options{
			BackendAddress:        "s3://aFakeBucket",
			MetadataFileExtension: "cfile",
		},
		fsys:         *localFs,
		awsRegion:    "us-east-1",
		s3Uploader:   fakeS3Server,
		s3Downloader: fakeS3Server,
	}
	var meta metadata.ObjectMetaData
	merr := json.Unmarshal([]byte(som), &meta)
	require.NoError(t, merr)

	err = store.Retrieve(context.Background(), metadata.CfileMetadataMap{
		"someObject.cfile": meta,
	}, "someObject.cfile")
	require.NoError(t, err)

	assert.NoError(t, err)

	// ensure the content of the file is correct
	b, err := afero.ReadFile(*localFs, "someObject")
	require.NoError(t, err)
	require.Equal(t, `tla`, string(b))

}

func TestS3GetBucketNameWithS3Prefix(t *testing.T) {
	expectedBackendAddress := "s3://aFakeBucket"

	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	memfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"test": {
			Content: []byte("bla"),
			ModTime: &mTime,
		},
	})
	assert.NoError(t, err)

	fakeS3Server := aferoS3Server{
		buckets: map[string]afero.Fs{
			// create a bucket in our fake s3 server
			"test": afero.NewMemMapFs(),
		},
	}

	store := S3Store{
		Options: Options{
			BackendAddress:        expectedBackendAddress,
			MetadataFileExtension: "cfile",
		},
		fsys:         *memfs,
		awsRegion:    "us-east-1",
		s3Uploader:   fakeS3Server,
		s3Downloader: fakeS3Server,
	}

	bucketName, err := store.getBucketName()

	assert.NoError(t, err)
	assert.Equal(t, "aFakeBucket", bucketName)

}

func TestS3GetBucketNameWithHTTPPrefix(t *testing.T) {
	type test struct {
		name        string
		addressName string
		want        string
	}
	tests := []test{
		{
			name:        "simple address",
			addressName: "http://127.0.0.1:9000/aFakeBucket",
			want:        "aFakeBucket",
		},
		{
			name:        "complex address",
			addressName: "http://127.0.0.1:9000/aFakeBucket/a/b/c",
			want:        "aFakeBucket",
		},
	}

	for _, test := range tests {
		expectedBackendAddress := test.addressName

		mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
		memfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
			"test": {
				Content: []byte("bla"),
				ModTime: &mTime,
			},
		})
		assert.NoError(t, err)

		fakeS3Server := aferoS3Server{
			buckets: map[string]afero.Fs{
				// create a bucket in our fake s3 server
				"test": afero.NewMemMapFs(),
			},
		}

		store := S3Store{
			Options: Options{
				BackendAddress:        expectedBackendAddress,
				MetadataFileExtension: "cfile",
			},
			fsys:         *memfs,
			awsRegion:    "us-east-1",
			s3Uploader:   fakeS3Server,
			s3Downloader: fakeS3Server,
		}

		bucketName, err := store.getBucketName()

		assert.NoError(t, err)
		assert.Equal(t, test.want, bucketName)
	}

}

func TestGetBucketNameUnsupportedPrefix(t *testing.T) {
	store := S3Store{
		Options: Options{
			BackendAddress:        "blah://thing",
			MetadataFileExtension: "cfile",
		},
	}

	_, err := store.getBucketName()
	assert.Error(t, err)
}

func TestRetrieveZeroCfiles(t *testing.T) {
	s := S3Store{}
	err := s.Retrieve(context.Background(), metadata.CfileMetadataMap{})
	require.ErrorIs(t, err, ErrCfilesLengthZero)
}
