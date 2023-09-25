package metadata

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/logger"
	"github.com/spf13/afero"
)

const MetadataFileExtension string = "cfile"

var (
	ErrFileExtensionEmpty          = fmt.Errorf("options.MetadatafileExtension cannot be %q", "")
	ErrRetrieveFailureHashMismatch = errors.New("hashes don't match, Retrieve aborted")
)

type ObjectMetaData struct {
	Name         string    `json:"name"`
	Checksum     string    `json:"checksum"`
	DateModified time.Time `json:"date_modified"`
}

type CfileMetadataMap map[string]ObjectMetaData

func HashFromCfileMatches(fsys afero.Fs, cfile string, expected string) (bool, error) {
	obj := strings.TrimSuffix(cfile, filepath.Ext(cfile))
	f, err := fsys.Open(obj)
	if err != nil {
		return false, err
	}
	actual, err := SHA256FromReader(f)
	if err != nil {
		return false, err
	}
	// If the hash of the downloaded file does not match the retrieved file, return an error
	if actual != expected {
		logger.Infof("Hash mismatch, got %s but expected %s", actual, expected)
		return false, ErrRetrieveFailureHashMismatch
	}
	if err := f.Close(); err != nil {
		return true, err
	}
	return true, nil

}

func SHA256FromReader(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", fmt.Errorf("%v: %s", err, "could not generate sha256 due to io.Copy error")
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func GenerateFromReader(name string, modTime time.Time, r io.Reader) (*ObjectMetaData, error) {
	hash, err := SHA256FromReader(r)
	if err != nil {
		return nil, err
	}
	return &ObjectMetaData{
		Name:         name,
		Checksum:     hash,
		DateModified: modTime,
	}, nil
}

func GenerateFromFile(f afero.File, key string) (*ObjectMetaData, error) {
	fstat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return GenerateFromReader(key, fstat.ModTime(), f)
}

func ParseCfile(fsys afero.Fs, obj string) (*ObjectMetaData, error) {
	cfile, err := fsys.Open(obj)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(cfile)
	if err != nil {
		return nil, err
	}
	var metadata ObjectMetaData
	if err := json.Unmarshal(b, &metadata); err != nil {
		return nil, fmt.Errorf("json marshal failed: %w", err)
	}
	return &metadata, nil
}

func ParseCfileWithExtension(fsys afero.Fs, obj, ext string) (*ObjectMetaData, error) {
	cfile := fmt.Sprintf("%s.%s", obj, ext)
	return ParseCfile(fsys, cfile)
}

type FsysWriteRequest struct {
	Object       string
	Fsys         afero.Fs
	Fi           afero.File
	MetadataPath string
	Extension    string
}

// WriteToFsys generates Cavorite metadata for req.Object and writes it to req.Fsys
func WriteToFsys(req FsysWriteRequest) (err error) {
	if req.MetadataPath == "" {
		return fmt.Errorf("req.MetadataPath cannot be %q", "")
	}
	logger.V(2).Infof("object: %s", req.Object)
	// generate metadata
	m, err := GenerateFromFile(req.Fi, req.Object)
	if err != nil {
		return err
	}
	logger.V(2).Infof("%s has a checksum of %q", req.Object, m.Checksum)
	// convert metadata to json
	blob, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return err
	}
	// Write metadata to disk
	metadataPath := fmt.Sprintf("%s.%s", req.MetadataPath, req.Extension)
	logger.V(2).Infof("writing metadata to %s", metadataPath)
	if err := afero.WriteFile(req.Fsys, metadataPath, blob, 0644); err != nil {
		return err
	}

	return nil
}
