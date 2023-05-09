package stores

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	gcsStorage "cloud.google.com/go/storage"
	"github.com/discentem/pantri_but_go/internal/metadata"
	"github.com/google/logger"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

type GCSStore struct {
	Options    Options `mapstructure:"options"`
	fsys       afero.Fs
	bucketName string
	gcsClient  *gcsStorage.Client
}

// NewGCSStoreClient creates a GCS Storage Client utilizing either the default GOOGLE_APPLICATION_CREDENTIAL  env var
// or a json string env var named PANTRI_GCS_CREDENTIALS
func NewGCSStoreClient(ctx context.Context, fsys afero.Fs, opts Options) (*GCSStore, error) {
	gcsDefault := viper.GetString("GOOGLE_APPLICATION_CREDENTIAL")
	var credentialsBytes []byte
	var client *gcsStorage.Client
	var err error
	// Look for the default google env var first
	if viper.IsSet("GOOGLE_APPLICATION_CREDENTIAL") {
		client, err = gcsStorage.NewClient(ctx, option.WithCredentialsFile(gcsDefault))
		if err != nil {
			return nil, err
		}
	} else if viper.IsSet("PANTRI_GCS_CREDENTIALS") {
		// If that doesn't exist, look for our own en var which contains a json string
		credentialsBytes = []byte(viper.GetString("PANTRI_GCS_CREDENTIALS"))
		client, err = gcsStorage.NewClient(ctx, option.WithCredentialsJSON(credentialsBytes))
		if err != nil {
			return nil, err
		}
	} else {
		// If we cannot find either, we cannot continue
		return nil, errors.New("No valid GCS credentials found. Exiting...")
	}

	// Create the GCS client with credentials
	defer client.Close()

	return &GCSStore{
		Options:    opts,
		fsys:       fsys,
		bucketName: opts.PantriAddress,
		gcsClient:  client,
	}, nil
}

func (s *GCSStore) GetOptions() Options {
	return s.Options
}

// TODO(discentem): #34 largely copy-pasted from stores/local/local.go. Can be consolidated
// Upload generates the metadata, writes it to disk and uploads the file to the GCS bucket
func (s *GCSStore) Upload(ctx context.Context, objects ...string) error {
	for _, o := range objects {
		logger.V(2).Infof("Object: %s\n", o)
		logger.V(2).Infof("Bucket: %v\n", s.bucketName)
		f, err := os.Open(o)
		if err != nil {
			return err
		}
		defer f.Close()

		// This should be moved to a helper as it is copy pasted from s3 store
		// generate pantri metadata
		m, err := metadata.GenerateFromFile(*f)
		if err != nil {
			return err
		}
		logger.V(2).Infof("%s has a checksum of %q", o, m.Checksum)
		// convert metadata to json
		blob, err := json.MarshalIndent(m, "", " ")
		if err != nil {
			return err
		}
		// Write metadata to disk
		if err := os.WriteFile(fmt.Sprintf("%s.%s", o, s.Options.MetaDataFileExtension), blob, 0644); err != nil {
			return err
		}

		// ToDo(natewalck) Expose this timeout as a setting
		ctx, cancel := context.WithTimeout(ctx, time.Second*1800)
		defer cancel()

		gcsObject := s.gcsClient.Bucket(s.bucketName).Object(o)
		// ToDo(natewalck) Maybe expose this as a setting?
		// Only allow the file to be written if it doesn't already exist.
		wc := gcsObject.If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx)

		// Reset to the start of the file because metadata generation has already read it once
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(wc, f)
		if err != nil {
			logger.V(2).Infof("Failed to upload %s", o)
			return err
		}

		if err := wc.Close(); err != nil {
			// Error will contain this string if the DoesNotExist condition isn't met
			if strings.Contains(err.Error(), "conditionNotMet") {
				logger.Infof("%s already exists, skipping...", o)
			} else {
				return err
			}
		}
	}
	return nil
}

// Retrieve gets the file from the GCS bucket, validates the hash is correct and writes it to disk
func (s *GCSStore) Retrieve(ctx context.Context, objects ...string) error {
	for _, o := range objects {
		objectPath := strings.TrimSuffix(o, filepath.Ext(o))
		logger.V(2).Infof("Retrieving %s", objectPath)

		var f *os.File
		if _, err := os.Stat(objectPath); err == nil {
			logger.V(2).Infof("%s already exists", objectPath)
			f, err := os.Open(objectPath)
			if err != nil {
				return err
			}
			defer f.Close()
		} else {
			// Create local path if it doesn't already exist
			f, err = os.Create(objectPath)
			if err != nil {
				return err
			}
			defer f.Close()
			// Download the file
			rc, err := s.gcsClient.Bucket(s.bucketName).Object(objectPath).NewReader(ctx)
			if err != nil {
				return err
			}
			defer rc.Close()

			file, err := os.Create(objectPath)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err = io.Copy(file, rc); err != nil {
				return err
			}
		}
		// Get the hash for the downloaded file
		hash, err := metadata.SHA256FromReader(f)
		if err != nil {
			return err
		}
		// Get the metadata from the metadata file
		m, err := metadata.ParsePfile(s.fsys, o)
		if err != nil {
			return err
		}
		// If the hash of the downloaded file does not match the retrieved file, return an error
		if hash != m.Checksum {
			logger.V(2).Infof("Hash mismatch, got %s but expected %s", hash, m.Checksum)
			os.Remove(objectPath)
			return ErrRetrieveFailureHashMismatch
		}

		return nil
	}
	return nil
}
