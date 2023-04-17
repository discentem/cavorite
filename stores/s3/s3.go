package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/logger"
	"github.com/spf13/afero"

	"github.com/discentem/pantri_but_go/internal/metadata"
	"github.com/discentem/pantri_but_go/stores"

	pantriconfig "github.com/discentem/pantri_but_go/pantri"

	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/mitchellh/mapstructure"
)

type Store struct {
	PantriAddress string         `mapstructure:"pantri_address"`
	Opts          stores.Options `mapstructure:"options"`
	awsRegion     string
	s3Client      *s3.Client
	s3Uploader    *s3manager.Uploader
	s3Downloader  *s3manager.Downloader
}

func (s *Store) bucketName() *string {
	_, buck := path.Split(s.PantriAddress)
	return aws.String(buck)
}

func getConfig(ctx context.Context, awsRegion string, pantriAddress string) (*aws.Config, error) {
	var cfg aws.Config

	logger.V(2).Infof("pantriAddress: %s", pantriAddress)
	if !strings.HasPrefix(pantriAddress, "s3://") && !strings.HasPrefix(pantriAddress, "https://") && !strings.HasPrefix(pantriAddress, "http://") {
		return nil, errors.New("pantriAddress did not contain s3://, http://, or https:// prefix")
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(awsRegion),
	)
	if err != nil {
		return &cfg, err
	}

	if strings.HasPrefix(pantriAddress, "https://") || strings.HasPrefix(pantriAddress, "http://") {
		server, _ := path.Split(pantriAddress)
		// https://stackoverflow.com/questions/67575681/is-aws-go-sdk-v2-integrated-with-local-minio-server
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               server,
				SigningRegion:     awsRegion,
				HostnameImmutable: true,
			}, nil
		})
		cfg.EndpointResolverWithOptions = resolver
	}

	return &cfg, nil

}

func (s *Store) init(ctx context.Context, awsRegion string, fsys afero.Fs, sourceRepo string) error {
	c := pantriconfig.Config{
		Type:          "s3",
		PantriAddress: s.PantriAddress,
		Opts:          s.Opts,
		Validate: func() error {
			_, err := s.s3Client.HeadBucket(
				ctx,
				&s3.HeadBucketInput{
					Bucket: s.bucketName(),
				},
			)
			return err
		},
	}

	return c.Write(fsys, sourceRepo)
}

func New(ctx context.Context, awsRegion string, fsys afero.Fs, sourceRepo, pantriAddress string, opts stores.Options) (*Store, error) {
	if opts.RemoveFromSourceRepo == nil {
		b := false
		opts.RemoveFromSourceRepo = &b
	}
	if opts.MetaDataFileExtension == "" {
		e := ".pfile"
		opts.MetaDataFileExtension = e
	}
	cfg, err := getConfig(
		ctx,
		awsRegion,
		pantriAddress,
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

	s := &Store{
		PantriAddress: pantriAddress,
		Opts:          opts,
		awsRegion:     awsRegion,
		s3Client:      s3Client,
		s3Uploader:    s3Uploader,
		s3Downloader:  s3Downloader,
	}

	err = s.init(
		ctx,
		awsRegion,
		fsys,
		sourceRepo,
	)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func Load(m map[string]interface{}) (stores.Store, error) {
	logger.Infof("type %q detected in pantri %q", m["type"], m["pantri_address"])
	var s *Store
	if err := mapstructure.Decode(m, &s); err != nil {
		return nil, err
	}
	return stores.Store(s), nil
}

// TODO(discentem): #34 largely copy-pasted from stores/local/local.go. Can be consolidated
func (s *Store) Upload(ctx context.Context, fsys afero.Fs, sourceRepo string, destination string, objects ...string) error {
	for _, object := range objects {
		logger.Info("starting upload")
		logger.V(2).Info(object)
		f, err := fsys.Open(object)
		if err != nil {
			return err
		}
		defer f.Close()
		fStat, err := f.Stat()
		if err != nil {
			return err
		}
		md, err := metadata.GenerateFromReader(
			object,
			fStat.ModTime(),
			f,
		)
		if err != nil {
			return err
		}

		logger.V(2).Infof("%s has a checksum of %q", object, md.Checksum)

		// Write the metadata file to disk
		err = metadata.WriteToFs(
			fsys,
			sourceRepo,
			*md,
			destination,
			s.Opts.MetaDataFileExtension,
		)
		if err != nil {
			return err
		}

		cfg, err := getConfig(ctx, s.awsRegion, s.PantriAddress)
		if err != nil {
			return err
		}
		s.s3Client = s3.NewFromConfig(*cfg)
		s.s3Uploader = s3manager.NewUploader(
			s.s3Client,
			func(u *s3manager.Uploader) {
				u.PartSize = 64 * 1024 * 1024 // 64MB per part
			},
		)
		// Reset offset after writing metadata, before attempting upload. This prevents a 0 byte upload.
		_, err = f.Seek(0, 0)
		if err != nil {
			return err
		}

		// Begin multipart upload
		_, err = s.s3Uploader.Upload(
			ctx,
			&s3.PutObjectInput{
				Bucket: s.bucketName(),
				Key:    aws.String(filepath.Join(filepath.Dir(destination), filepath.Base(object))),
				Body:   f,
			},
		)
		logger.Infof("uploaded %s to %s", object, destination)
		logger.V(2).Infof("uploaded %s to %s via multipart upload", object, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

//lint:ignore U1000 func unused
func keyFromPfile(fsys afero.Fs, path, ext string) (*string, error) {
	m, err := metadata.ParsePfile(fsys, path, ".pfile")
	if err != nil {
		return nil, err
	}
	return &m.Name, nil
}

func (s *Store) Retrieve(ctx context.Context, fsys afero.Fs, sourceRepo string, pfiles ...string) error {
	for _, o := range pfiles {
		retrievePath := strings.TrimSuffix(filepath.Join(sourceRepo, o), s.Opts.MetaDataFileExtension)
		f, err := fsys.Create(retrievePath)
		if err != nil {
			return err
		}
		defer f.Close()
		obj := &s3.GetObjectInput{
			Bucket: s.bucketName(),
			Key:    aws.String(o),
		}
		_, err = s.s3Downloader.Download(
			ctx,
			f,
			obj,
		)
		if err != nil {
			if rmvFerr := fsys.Remove(retrievePath); rmvFerr != nil {
				return fmt.Errorf("failed to remove %s due %w after encountering download failure: %w", retrievePath, rmvFerr, err)
			}
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

		m, err := metadata.ParsePfile(fsys, pfilePath, ext)
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
		if err := afero.WriteFile(fsys, op, b, 0644); err != nil {
			return err
		}
	}
	return nil
}
