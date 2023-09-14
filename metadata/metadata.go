package metadata

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/logger"
	"github.com/spf13/afero"
)

const MetadataFileExtension string = "cfile"

var ErrFileExtensionEmpty = fmt.Errorf("options.MetadatafileExtension cannot be %q", "")

type ObjectMetaData struct {
	Name         string    `json:"name"`
	Checksum     string    `json:"checksum"`
	DateModified time.Time `json:"date_modified"`
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

func GenerateFromFile(f afero.File) (*ObjectMetaData, error) {
	fstat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return GenerateFromReader(fstat.Name(), fstat.ModTime(), f)
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
	Object    string
	Fsys      afero.Fs
	Fi        afero.File
	Extension string
}

// WriteMetadata generates Cavorite metadata for obj and writes it to s.Fsys
func WriteToFsys(req FsysWriteRequest) (err error) {
	logger.V(2).Infof("object: %s", req.Object)
	// generate metadata
	m, err := GenerateFromFile(req.Fi)
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
	metadataPath := fmt.Sprintf("%s.%s", req.Object, req.Extension)
	logger.V(2).Infof("writing metadata to %s", metadataPath)
	if err := afero.WriteFile(req.Fsys, metadataPath, blob, 0644); err != nil {
		return err
	}

	return nil
}
