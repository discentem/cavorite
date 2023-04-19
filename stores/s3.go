package stores

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/logger"
	"github.com/spf13/afero"

	"github.com/discentem/pantri_but_go/internal/metadata"

	pantriconfig "github.com/discentem/pantri_but_go/pantri"

	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

type S3Store struct {
	Opts         Options `mapstructure:"options"`
	fs           afero.Fs
	awsRegion    string
	s3Client     *s3.Client
	s3Uploader   *s3manager.Uploader
	s3Downloader *s3manager.Downloader
}

func (s *S3Store) bucketName() *string {
	_, buck := path.Split(s.Opts.PantriAddress)
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

func (s *S3Store) WriteConfig(ctx context.Context, sourceRepo string) error {
	c := pantriconfig.Config{
		Type:          "s3",
		PantriAddress: s.Opts.PantriAddress,
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

	return c.Write(s.fs, sourceRepo)
}

func NewS3StoreClient(ctx context.Context, fs afero.Fs, awsRegion, sourceRepo string, opts Options) (*S3Store, error) {
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
		fs:           fs,
		awsRegion:    awsRegion,
		s3Client:     s3Client,
		s3Uploader:   s3Uploader,
		s3Downloader: s3Downloader,
	}, nil
}

// func Load(m map[string]interface{}) (Store, error) {
// 	logger.Infof("type %q detected in pantri %q", m["type"], m["pantri_address"])
// 	var s *S3Store
// 	if err := mapstructure.Decode(m, &s); err != nil {
// 		return nil, err
// 	}
// 	return Store(s), nil
// }

// TODO(discentem): #34 largely copy-pasted from stores/local/local.go. Can be consolidated
func (s *S3Store) Upload(ctx context.Context, sourceRepo string, destination string, objects ...string) error {
	for _, object := range objects {
		logger.Info("starting upload")
		logger.V(2).Info(object)
		f, err := s.fs.Open(object)
		if err != nil {
			return err
		}
		defer f.Close()
		fStat, err := f.Stat()
		if err != nil {
			return err
		}
		md, err := metadata.GenerateFromReader(
			destination,
			fStat.ModTime(),
			f,
		)
		if err != nil {
			return err
		}

		logger.V(2).Infof("%s has a checksum of %q", object, md.Checksum)

		// Write the metadata file to disk
		err = metadata.WriteToFs(
			s.fs,
			sourceRepo,
			*md,
			filepath.Base(destination),
			s.Opts.MetaDataFileExtension,
		)
		if err != nil {
			return err
		}

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

func (s *S3Store) Retrieve(ctx context.Context, sourceRepo string, pfiles ...string) error {
	log.Println(pfiles)
	log.Println(sourceRepo)
	log.Printf("bucket name: %s", *s.bucketName())
	for v, o := range pfiles {
		fmt.Printf("%d, retrieve path: %s\n",
			v,
			strings.TrimSuffix(filepath.Join(sourceRepo, o), s.Opts.MetaDataFileExtension),
		)
		retrievePath := strings.TrimSuffix(filepath.Join(sourceRepo, o), s.Opts.MetaDataFileExtension)
		f, err := s.fs.Create(retrievePath)
		if err != nil {
			return err
		}
		// fstruct, ok := f.()
		// if !ok {
		// 	return err
		// }
		defer f.Close()

		// test / chrome / googlechromebeta.dmg
		test_key := "chrome/googlechromebeta.dmg"
		obj := &s3.GetObjectInput{
			Bucket: s.bucketName(),
			Key:    aws.String(test_key),
		}
		logger.Infof("f: %+v", f)
		// w := io.WriterAt(f)
		// fmt.Print(w)
		_, err = f.Seek(0, 0)
		if err != nil {
			return err
		}
		//The w io.WriterAt can be satisfied by an os.File to do multipart concurrent downloads, or in memory []byte wrapper using aws.WriteAtBuffer. In case you download files into memory do not forget to pre-allocate memory to avoid additional allocations and GC runs.
		// io.WriterAt()
		numOfBytesDownloaded, err := s.s3Downloader.Download(
			ctx,
			f,
			obj,
		)
		if err != nil {
			if rmvFerr := s.fs.Remove(retrievePath); rmvFerr != nil {
				return fmt.Errorf("failed to remove %s due %w after encountering download failure: %w", retrievePath, rmvFerr, err)
			}
			return err
		}

		log.Printf("file downloaded, bytes transferred: %d", numOfBytesDownloaded)

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
			return ErrRetrieveFailureHashMismatch
		}
		op := path.Join(sourceRepo, o)
		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		if err := afero.WriteFile(s.fs, op, b, 0644); err != nil {
			return err
		}
	}
	return nil
}
