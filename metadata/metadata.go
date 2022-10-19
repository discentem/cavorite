package metadata

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type ObjectMetaData struct {
	Name         string    `json:"name"`
	Checksum     string    `json:"checksum"`
	DateModified time.Time `json:"date_modified"`
}

func SHA256FromBytes(b []byte) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, bytes.NewReader(b)); err != nil {
		return "", err
	}
	h := hex.EncodeToString(hash.Sum(nil))
	return h, nil
}

func GenerateFromReader(name string, modTime time.Time, r io.Reader) (*ObjectMetaData, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	hash, err := SHA256FromBytes(b)
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

func ParsePfile(obj, ext string) (*ObjectMetaData, error) {
	pfile, err := os.Open(fmt.Sprintf("%s.%s", obj, ext))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(pfile)
	if err != nil {
		return nil, err
	}
	var metadata *ObjectMetaData
	if err := json.Unmarshal(b, &metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}
