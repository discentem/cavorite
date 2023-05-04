package metadata

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/afero"
)

type ObjectMetaData struct {
	Name         string    `json:"name"`
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

func GenerateFromFile(f os.File) (*ObjectMetaData, error) {
	fstat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return GenerateFromReader(fstat.Name(), fstat.ModTime(), &f)
}

func ParsePfile(fsys afero.Fs, obj string) (*ObjectMetaData, error) {
	pfile, err := fsys.Open(obj)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(pfile)
	if err != nil {
		return nil, err
	}
	var metadata ObjectMetaData
	if err := json.Unmarshal(b, &metadata); err != nil {
		return nil, fmt.Errorf("json marshal failed: %w", err)
	}
	return &metadata, nil
}

func ParsePfileWithExtension(fsys afero.Fs, obj, ext string) (*ObjectMetaData, error) {
	pfile := fmt.Sprintf("%s.%s", obj, ext)
	return ParsePfile(fsys, pfile)
}
