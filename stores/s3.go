package stores

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/logger"
	"github.com/spf13/afero"

	"github.com/discentem/cavorite/fileutils"
	"github.com/discentem/cavorite/metadata"
	multierr "github.com/hashicorp/go-multierror"

	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

type S3Downloader interface {
	Download(
		ctx context.Context,
		w io.WriterAt,
		input *s3.GetObjectInput,
		options ...func(*s3manager.Downloader)) (n int64, err error)
}

type S3Uploader interface {
	Upload(ctx context.Context,
		input *s3.PutObjectInput,
		opts ...func(*s3manager.Uploader)) (
		*s3manager.UploadOutput, error,
	)
}

type S3Store struct {
	Options      Options `json:"options" mapstructure:"options"`
	fsys         afero.Fs
	awsRegion    string
	s3Uploader   S3Uploader
	s3Downloader S3Downloader
}

func getConfig(ctx context.Context, region string, address string) (*aws.Config, error) {
	var cfg aws.Config
	var err error

	switch {
	case strings.HasPrefix(address, "s3://"):
		cfg, err = awsConfig.LoadDefaultConfig(
			ctx,
			awsConfig.WithRegion(region),
		)
		if err != nil {
			return nil, err
		}
	case strings.HasPrefix(address, "http://"):
		fallthrough
	case strings.HasPrefix(address, "https://"):
		server, _ := path.Split(address)
		// https://stackoverflow.com/questions/67575681/is-aws-go-sdk-v2-integrated-with-local-minio-server
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               server,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		})

		cfg, err = awsConfig.LoadDefaultConfig(
			ctx,
			awsConfig.WithRegion(region),
			awsConfig.WithEndpointResolverWithOptions(resolver),
			awsConfig.WithClientLogMode(aws.LogRequest|aws.LogRequestWithBody),
		)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("address did not contain s3://, http://, or https:// prefix")
	}

	return &cfg, nil
}

func NewS3StoreClient(ctx context.Context, fsys afero.Fs, opts Options) (*S3Store, error) {
	cfg, err := getConfig(
		ctx,
		opts.Region,
		opts.BackendAddress,
	)
	if err != nil {
		return nil, err
	}

	// TODO (@radsec) - extract this S3 logic to a separate internal client instead of directly from AWS_SDK
	s3Client := s3.NewFromConfig(*cfg)
	s3Uploader := s3manager.NewUploader(
		s3Client,
		func(u *s3manager.Uploader) {
			u.PartSize = 64 * 1024 * 1024 // 64MB per part
		},
	)
	s3Downloader := s3manager.NewDownloader(
		s3Client,
		func(d *s3manager.Downloader) {
			d.Concurrency = 3
		},
	)

	return &S3Store{
		Options:   opts,
		fsys:      fsys,
		awsRegion: opts.Region,
		// s3Uploader meets our interface for S3Uploader
		s3Uploader: s3Uploader,
		// s3Downloader meets our interface for S3Downloader
		s3Downloader: s3Downloader,
	}, nil
}

func (s *S3Store) GetOptions() (Options, error) {
	return s.Options, nil
}

func (s *S3Store) GetFsys() (afero.Fs, error) {
	return s.fsys, nil
}

// Upload generates the metadata, writes it to disk and uploads the file to the S3 bucket
func (s *S3Store) Upload(ctx context.Context, objects ...string) error {
	for _, o := range objects {
		f, err := s.fsys.Open(o)
		if err != nil {
			return err
		}

		// Generate S3 struct for object and upload to S3 bucket
		s3BucketName, err := s.getBucketName()
		if err != nil {
			logger.Errorf("error encountered parsing backend address: %v", err)
			return err
		}
		obj := s3.PutObjectInput{
			Bucket: aws.String(s3BucketName),
			Key:    aws.String(o),
			Body:   f,
		}
		out, err := s.s3Uploader.Upload(ctx, &obj)
		if err != nil {
			logger.Error(out)
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}

var (
	ErrCfilesLengthZero = errors.New("at least one cfile must be specified")
)

// Retrieve gets the file from the S3 bucket, validates the hash is correct and writes it to disk
func (s *S3Store) Retrieve(ctx context.Context, mmap metadata.CfileMetadataMap, cfiles ...string) error {
	var result *multierr.Error
	if len(cfiles) == 0 {
		return ErrCfilesLengthZero
	}

	for _, cfile := range cfiles {
		// We will either read the file that already exists or download it because it
		// is missing
		f, err := fileutils.OpenOrCreateFile(s.fsys, mmap[cfile].Name)
		if err != nil {
			result = multierr.Append(result, err)
			continue
		}
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			result = multierr.Append(result, err)
			continue
		}
		s3BucketName, err := s.getBucketName()
		if err != nil {
			e := fmt.Errorf("error encountered parsing backend address: %v", err)
			result = multierr.Append(result, e)
			continue
		}
		logger.V(2).Infof("mmap: %s", mmap)
		m, ok := mmap[cfile]
		if !ok {
			result = multierr.Append(result, fmt.Errorf("%q not found in mmap", cfile))
			continue
		}

		obj := &s3.GetObjectInput{
			Bucket: aws.String(s3BucketName),
			Key:    aws.String(m.Name),
		}
		// Download the file
		_, err = s.s3Downloader.Download(ctx, f, obj)
		if err != nil {
			result = multierr.Append(result, err)
			continue
		}

		matches, err := metadata.HashFromCfileMatches(s.fsys, cfile, m.Checksum)
		if err != nil {
			result = multierr.Append(result, err)
			continue
		}
		if !matches {
			logger.V(2).Infof("hash for %s did not match expected hash (%q) in %q", m.Name, m.Checksum, cfile)
			if err := s.fsys.Remove(m.Name); err != nil {
				result = multierr.Append(result, err)
				continue
			}
			result = multierr.Append(result, metadata.ErrRetrieveFailureHashMismatch)
			continue
		}
	}
	return result.ErrorOrNil()
}

func (s *S3Store) getBucketName() (string, error) {
	var bucketName string
	logger.Infof("s3store getBucketName: backend address: %s", s.Options.BackendAddress)
	switch {
	case strings.HasPrefix(s.Options.BackendAddress, "s3://"):
		s3BucketUrl, err := url.Parse(s.Options.BackendAddress)
		if err != nil {
			return "", err
		}
		bucketName = s3BucketUrl.Host
	case strings.HasPrefix(s.Options.BackendAddress, "http://"):
		fallthrough
	case strings.HasPrefix(s.Options.BackendAddress, "https://"):
		up, err := url.Parse(s.Options.BackendAddress)
		if err != nil {
			return "", err
		}
		p := strings.SplitN(up.Path, "/", 2)[1]
		fragments := strings.Split(p, "/")
		bucket := fragments[0]
		return bucket, nil
	default:
		return "", fmt.Errorf("unsupported s3 backend address")
	}

	logger.V(2).Infof("s3store getBucketName: bucket: %s", bucketName)
	return bucketName, nil
}

func (s *S3Store) Close() error {
	// FIXME: implement
	return nil
}
