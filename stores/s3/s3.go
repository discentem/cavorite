package s3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	pantriconfig "github.com/discentem/pantri_but_go/config"
	"github.com/discentem/pantri_but_go/metadata"
	"github.com/discentem/pantri_but_go/stores"
)

type Store struct {
	PantriAddress string         `mapstructure:"pantri_address"`
	Opts          stores.Options `mapstructure:"options"`
}

func (s *Store) init(sourceRepo string) error {
	c := pantriconfig.Config{
		Type:          "s3",
		PantriAddress: s.PantriAddress,
		Opts:          s.Opts,
		Validate: func() error {
			// use auth that was configued by aws cli
			sess := session.Must(session.NewSessionWithOptions(
				session.Options{
					// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-the-region
					SharedConfigState: session.SharedConfigEnable,
				},
			))
			uploader := s3manager.NewUploader(sess)
			buck := strings.TrimPrefix(s.PantriAddress, "s3://")
			_, err := uploader.S3.HeadBucket(&s3.HeadBucketInput{
				Bucket: &buck,
			})
			if err != nil {
				log.Print("error from S3.HeadBucket")
				return err
			}
			return nil
		},
	}

	return c.WriteToDisk(sourceRepo)
}

func New(sourceRepo, pantriAddress string, o stores.Options) (*Store, error) {
	if o.RemoveFromSourceRepo == nil {
		b := false
		o.RemoveFromSourceRepo = &b
	}
	s := &Store{
		PantriAddress: pantriAddress,
		Opts:          o,
	}
	err := s.init(sourceRepo)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// TODO(discentem): #34 largely copy-pasted from stores/local/local.go. Can be consolidated
func (s *Store) Upload(sourceRepo string, objects ...string) error {
	// use auth that was configued by aws cli
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-the-region
			SharedConfigState: session.SharedConfigEnable,
		},
	))

	uploader := s3manager.NewUploader(sess)

	for _, o := range objects {
		f, err := os.Open(o)
		if err != nil {
			return err
		}
		// TODO(discentem): probably inefficient, reading same file multiple times
		var b []byte
		fstat, err := f.Stat()
		if err != nil {
			return err
		}
		// get bytes of the file
		b = make([]byte, fstat.Size())
		_, err = f.Read(b)
		if err != nil {
			return err
		}
		defer f.Close()

		// generate pantri metadata
		m, err := metadata.GenerateFromFile(*f)
		if err != nil {
			return err
		}
		// convert to json
		blob, err := json.MarshalIndent(m, "", " ")
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(path.Join(sourceRepo, fmt.Sprintf("%s.pfile", o))), os.ModePerm); err != nil {
			return err
		}
		if err := os.WriteFile(path.Join(sourceRepo, fmt.Sprintf("%s.pfile", o)), blob, 0644); err != nil {
			return err
		}

		if s.Opts.RemoveFromSourceRepo != nil {
			if *s.Opts.RemoveFromSourceRepo {
				if err := os.Remove(o); err != nil {
					return err
				}
			}
		}
		buck := strings.TrimPrefix(s.PantriAddress, "s3://")
		out, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(buck),
			Key:    &o,
			Body:   bytes.NewReader(b),
		})
		if err != nil {
			fmt.Println(out)
			return err
		}
	}
	return nil
}

func (s *Store) Retrieve(sourceRepo string, objects ...string) error {
	return nil
}
