package stores

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	gcsStorage "cloud.google.com/go/storage"
	"github.com/discentem/cavorite/internal/metadata"
	"github.com/google/logger"
	"github.com/spf13/afero"
	"google.golang.org/api/option"
)

type GCSStore struct {
	Options   Options `mapstructure:"options"`
	fsys      afero.Fs
	gcsClient *gcsStorage.Client
}

// NewGCSStoreClient creates a GCS Storage Client utilizing either the default GOOGLE_APPLICATION_CREDENTIAL  env var
// or a json string env var named CAVORITE_GCS_CREDENTIALS
func NewGCSStoreClient(ctx context.Context, fsys afero.Fs, opts Options) (*GCSStore, error) {
	gcsDefault := os.Getenv("GOOGLE_APPLICATION_CREDENTIAL")
	cavoriteGCSCreds := os.Getenv("CAVORITE_GCS_CREDENTIALS")
	var client *gcsStorage.Client
	var err error
	// Look for the default google env var first
	if gcsDefault != "" {
		client, err = gcsStorage.NewClient(ctx, option.WithCredentialsFile(gcsDefault))
		if err != nil {
			return nil, err
		}
	} else if cavoriteGCSCreds != "" {
		// If that doesn't exist, look for our own en var which contains a json string
		credentialsBytes := []byte(cavoriteGCSCreds)
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
		Options:   opts,
		fsys:      fsys,
		gcsClient: client,
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
		f, err := s.fsys.Open(o)
		if err != nil {
			return err
		}
		defer f.Close()

		// Generate cavorite metadata
		m, err := metadata.GenerateFromFile(f)
		if err != nil {
			return err
		}
		logger.V(2).Infof("%s has a checksum of %q", o, m.Checksum)
		// convert metadata to json
		blob, err := json.MarshalIndent(m, "", " ")
		if err != nil {
			return err
		}
		// Write metadata to fsys
		fname := fmt.Sprintf("%s.%s", o, s.Options.MetaDataFileExtension)
		if err := afero.WriteFile(s.fsys, fname, blob, 0644); err != nil {
			return err
		}

		// ToDo(natewalck) Expose this timeout as a setting
		ctx, cancel := context.WithTimeout(ctx, time.Second*1800)
		defer cancel()

		gcsObject := s.gcsClient.Bucket(s.Options.BackendAddress).Object(o)
		// ToDo(natewalck) Maybe expose this as a setting?
		// Only allow the file to be written if it doesn't already exist.
		wc := gcsObject.If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx)

		// Reset to the start of the file because metadata generation has already read it once
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			return err
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
// if err := s.Retrieve(cmd.Context(), objects...); err != nil {
func (s *GCSStore) Retrieve(ctx context.Context, metaObjects ...string) error {
	for _, mo := range metaObjects {
		// For Retrieve, the object is the cfile itself, which we derive the actual filename from
		objectPath := inferObjPath(mo)
		// We will either read the file that already exists or download it because it
		// is missing
		f, err := openOrCreateFile(s.fsys, objectPath)
		if err != nil {
			return err
		}
		fileInfo, err := f.Stat()
		if err != nil {
			return err
		}
		if fileInfo.Size() > 0 {
			logger.Infof("%s already exists", objectPath)
		} else { // Download the file as it doesn't exist on disk
			rc, err := s.gcsClient.Bucket(s.Options.BackendAddress).Object(objectPath).NewReader(ctx)
			if err != nil {
				return err
			}
			defer rc.Close()

			f, err := s.fsys.Create(objectPath)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err = io.Copy(f, rc); err != nil {
				return err
			}
		}
		// Get the hash for the downloaded file
		hash, err := metadata.SHA256FromReader(f)
		if err != nil {
			return err
		}
		// Get the metadata from the metadata file
		m, err := metadata.ParseCfile(s.fsys, mo)
		if err != nil {
			return err
		}
		// If the hash of the downloaded file does not match the retrieved file, return an error
		if hash != m.Checksum {
			logger.V(2).Infof("Hash mismatch, got %s but expected %s", hash, m.Checksum)
			if err := s.fsys.Remove(objectPath); err != nil {
				return err
			}
			return ErrRetrieveFailureHashMismatch
		}
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}
