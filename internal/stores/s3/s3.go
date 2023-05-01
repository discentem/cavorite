package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/logger"
	"github.com/spf13/afero"

	"github.com/discentem/pantri_but_go/internal/metadata"
	"github.com/discentem/pantri_but_go/internal/stores"

	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

type S3Store struct {
	Opts      stores.Options `mapstructure:"options"`
	fsys      afero.Fs
	awsRegion string
	// Migrate to internal/s3Client instead of using s3Client directly from AWS_SDK
	s3Client     *s3.Client
	s3Uploader   *s3manager.Uploader
	s3Downloader *s3manager.Downloader
	//
}

func NewS3StoreClient(ctx context.Context, fsys afero.Fs, awsRegion, sourceRepo string, opts stores.Options) (*S3Store, error) {
	if opts.RemoveFromSourceRepo == nil {
		b := false
		opts.RemoveFromSourceRepo = &b
	}
	if opts.MetaDataFileExtension == "" {
		e := ".pfile"
		opts.MetaDataFileExtension = e
	}
	cfg, err := getConfig(
		awsRegion,
		opts.PantriAddress,
	)
	if err != nil {
		return nil, err
	}

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
		Opts:         opts,
		fsys:         fsys,
		awsRegion:    awsRegion,
		s3Client:     s3Client,
		s3Uploader:   s3Uploader,
		s3Downloader: s3Downloader,
	}, nil
}

func getConfig(region string, pantriAddress string) (*aws.Config, error) {
	var cfg aws.Config
	var err error

	if strings.HasPrefix(pantriAddress, "s3://") {
		cfg, err = config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	} else if strings.HasPrefix(pantriAddress, "https://") || strings.HasPrefix(pantriAddress, "http://") {
		// e.g. http://127.0.0.1:9000/test becomes http://127.0.0.1:9000
		server, _ := path.Split(pantriAddress)
		// https://stackoverflow.com/questions/67575681/is-aws-go-sdk-v2-integrated-with-local-minio-server
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               server,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		})

		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(resolver),
		)
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return nil, errors.New("pantriAddress did not contain s3://, http://, or https:// prefix")
}

// func (s *S3Store) init(ctx context.Context, fsys afero.Fs, sourceRepo string) error {
// 	c := pantri.Config{
// 		Type:          "s3",
// 		Opts:          s.Opts,
// 		Validate: func() error {
// 			cfg, err := getConfig(s.awsRegion, s.Opts.PantriAddress)
// 			if err != nil {
// 				return err
// 			}
// 			uploader := s3.NewFromConfig(*cfg)
// 			// s3://test --> test
// 			// http://stuff/test --> test
// 			_, buck := path.Split(s.Opts.PantriAddress)
// 			_, err = uploader.HeadBucket(ctx, &s3.HeadBucketInput{
// 				Bucket: &buck,
// 			})
// 			return err
// 		},
// 	}

// 	return c.Write(fsys, sourceRepo)
// }

// func New(ctx context.Context, fsys afero.Fs, sourceRepo, pantriAddress string, o stores.Options) (*S3Store, error) {
// 	if o.RemoveFromSourceRepo == nil {
// 		b := false
// 		o.RemoveFromSourceRepo = &b
// 	}
// 	if o.MetaDataFileExtension == "" {
// 		e := ".pfile"
// 		o.MetaDataFileExtension = e
// 	}
// 	s := &S3Store{
// 		Opts: o,
// 	}
// 	err := s.init(ctx, fsys, sourceRepo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return s, nil
// }

func (s *S3Store) GetOptions() stores.Options {
	return s.Opts
}

// TODO(discentem): #34 largely copy-pasted from stores/local/local.go. Can be consolidated
func (s *S3Store) Upload(ctx context.Context, sourceRepo string, objects ...string) error {
	for _, o := range objects {
		f, err := os.Open(o)
		if err != nil {
			return err
		}
		defer f.Close()
		// TODO(discentem): probably inefficient, reading entire file into memory
		b, err := os.ReadFile(o)
		if err != nil {
			return err
		}

		// generate pantri metadata
		m, err := metadata.GenerateFromFile(*f)
		if err != nil {
			return err
		}
		logger.V(2).Infof("%s has a checksum of %q", o, m.Checksum)
		// convert to json
		blob, err := json.MarshalIndent(m, "", " ")
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(path.Join(sourceRepo, fmt.Sprintf("%s.%s", o, s.Opts.MetaDataFileExtension))), os.ModePerm); err != nil {
			return err
		}
		if err := os.WriteFile(path.Join(sourceRepo, fmt.Sprintf("%s.%s", o, s.Opts.MetaDataFileExtension)), blob, 0644); err != nil {
			return err
		}

		if s.Opts.RemoveFromSourceRepo != nil {
			if *s.Opts.RemoveFromSourceRepo {
				if err := os.Remove(o); err != nil {
					return err
				}
			}
		}

		_, buck := path.Split(s.Opts.PantriAddress)
		obj := s3.PutObjectInput{
			Bucket: aws.String(buck),
			Key:    &o,
			Body:   bytes.NewReader(b),
		}
		out, err := s.s3Uploader.Upload(ctx, &obj)
		if err != nil {
			logger.Error(out)
			return err
		}
	}
	return nil
}

func (s *S3Store) Retrieve(ctx context.Context, sourceRepo string, objects ...string) error {
	for _, o := range objects {
		retrievePath := filepath.Join(sourceRepo, o)
		f, err := os.Create(retrievePath)
		if err != nil {
			return err
		}
		defer f.Close()
		_, buck := path.Split(s.Opts.PantriAddress)
		obj := &s3.GetObjectInput{
			Bucket: aws.String(buck),
			Key:    aws.String(o),
		}
		_, err = s.s3Downloader.Download(ctx, f, obj)
		if err != nil {
			return err
		}
		hash, err := metadata.SHA256FromReader(f)
		if err != nil {
			return err
		}
		var ext string
		if s.Opts.MetaDataFileExtension == "" {
			ext = ".pfile"
		} else {
			ext = s.Opts.MetaDataFileExtension
		}
		pfilePath := filepath.Join(sourceRepo, o)

		m, err := metadata.ParsePfile(s.fsys, pfilePath, ext)
		if err != nil {
			return err
		}
		if hash != m.Checksum {
			fmt.Println(hash, m.Checksum)
			return stores.ErrRetrieveFailureHashMismatch
		}
		op := path.Join(sourceRepo, o)
		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		if err := os.WriteFile(op, b, 0644); err != nil {
			return err
		}
	}
	return nil
}

// func Load(m map[string]interface{}, ctx context.Context, fsys afero.Fs, region string) (stores.Store, error) {
// 	logger.Infof("type %q detected in pantri %q", m["type"], m["pantri_address"])
// 	if err := mapstructure.Decode(m, &s); err != nil {
// 		return nil, err
// 	}
// 	s3, err := s3.NewS3StoreClient(ctx, fsys, "us-east-1", sourceRepo, opts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return stores.Store(s3), nil
// }

// // from loader //////////////////////////////////
// func Initialize(ctx context.Context, fsys afero.Fs, sourceRepo, backend, address string, opts stores.Options) (stores.Store, error) {
// 	var s stores.Store
// 	switch b := (backend); b {
// 	case "s3":
// 		s3, err := s3.NewS3StoreClient(ctx, fsys, "us-east-1", sourceRepo, opts)
// 		if err != nil {
// 			return nil, err
// 		}
// 		s = stores.Store(s3)
// 	default:
// 		return nil, fmt.Errorf("%s: %w", b, ErrUnsupportedStore)
// 	}
// 	return s, nil

// }

// //
