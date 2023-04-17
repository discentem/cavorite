package metadata

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/logger"
	"github.com/spf13/afero"
)

const (
	DefaultFileExtension string = "pfile"
)

type Object struct {
	Key          string    `json:"key"`
	Checksum     string    `json:"checksum"`
	DateModified time.Time `json:"date_modified"`
}

func SHA256FromReader(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func GenerateFromReader(key string, modTime time.Time, r io.Reader) (*Object, error) {
	hash, err := SHA256FromReader(r)
	if err != nil {
		return nil, err
	}
	logger.V(2).Infof("key: %s", key)
	return &Object{
		Key:          key,
		Checksum:     hash,
		DateModified: modTime,
	}, nil
}

func GenerateFromFile(f os.File) (*Object, error) {
	fstat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return GenerateFromReader(fstat.Name(), fstat.ModTime(), &f)
}

func ParsePfile(fsys afero.Fs, obj, ext string) (*Object, error) {
	pfile, err := fsys.Open(fmt.Sprintf("%s.%s", obj, ext))
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(pfile)
	if err != nil {
		return nil, err
	}
	var metadata Object
	if err := json.Unmarshal(b, &metadata); err != nil {
		return nil, fmt.Errorf("json marshal failed: %w", err)
	}
	return &metadata, nil
}

func createMetadataFolder(fsys afero.Fs, filename string) error {
	return fsys.MkdirAll(
		filepath.Dir(filename),
		os.ModePerm,
	)
}

// WriteToFs writes object to fs as json in a file named object.Name in destination.
// destination should be a relative path.
func WriteToFs(fsys afero.Fs, sourceRepo string, object Object, destination string, extension string) error {
	filename := fmt.Sprintf("%s.%s", object.Key, extension)
	filepath := filepath.Join(sourceRepo, filepath.Dir(destination), filename)
	logger.V(2).Infof("filepath: %s", filepath)
	logger.V(2).Infof("sourceRepo: %s", sourceRepo)
	err := createMetadataFolder(fsys, filepath)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(object, "", " ")
	if err != nil {
		return err
	}
	logger.V(2).Infof("marshaled %s json into bytes", object.Key)
	logger.V(2).Infof("attempting to remove %s", filepath)
	if err := fsys.Remove(filepath); err != nil {
		if _, ok := err.(*os.PathError); ok {
			// ignore if path doesn't exist
		} else {
			return err
		}
	}
	f, err := fsys.Create(filepath)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}
