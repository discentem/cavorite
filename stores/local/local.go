package local

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/discentem/pantri_but_go/metadata"
	"github.com/discentem/pantri_but_go/stores"
)

type Store struct {
	SourceRepo string         `json:"source_repo"`
	Pantri     string         `json:"pantri"`
	Opts       stores.Options `json:"options"`
}

type PantriConfig struct {
	Type string `json:"type"`
	Store
}

func (c *PantriConfig) MarshalJSON() ([]byte, error) {
	conf := struct {
		Type       string         `json:"type"`
		SourceRepo string         `json:"source_repo"`
		Pantri     string         `json:"pantri"`
		Opts       stores.Options `json:"options"`
	}{
		Type:       c.Type,
		SourceRepo: c.Store.SourceRepo,
		Pantri:     c.Store.Pantri,
		Opts:       c.Store.Opts,
	}
	return json.Marshal(conf)
}

func NewWithOptions(sourceRepo, pantri string, o stores.Options) (*Store, error) {
	if o.RemoveFromSourceRepo == nil {
		b := false
		o.RemoveFromSourceRepo = &b
	}
	s := &Store{
		SourceRepo: sourceRepo,
		Pantri:     pantri,
		Opts:       o,
	}
	err := s.init()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func New(sourceRepo, pantri string, removeFromSourceRepo *bool) (*Store, error) {
	return NewWithOptions(sourceRepo, pantri, stores.Options{
		RemoveFromSourceRepo: removeFromSourceRepo,
	})
}

func (s *Store) config() error {
	localStoreconfig := PantriConfig{
		Type:  "local",
		Store: *s,
	}
	b, err := json.MarshalIndent(localStoreconfig, "", "  ")
	if err != nil {
		return err
	}
	dir := s.SourceRepo
	if err := os.MkdirAll(fmt.Sprintf("%s/.pantri", dir), os.ModePerm); err != nil {
		return err
	}
	cfile := fmt.Sprintf("%s/.pantri/config", dir)
	if err := os.WriteFile(cfile, b, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (s *Store) init() error {
	if err := s.config(); err != nil {
		return err
	}
	return nil
}

func sha256FromBytes(b []byte) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, bytes.NewReader(b)); err != nil {
		return "", err
	}
	h := hex.EncodeToString(hash.Sum(nil))
	return h, nil
}

func (s *Store) generateMetadata(f os.File) (*metadata.ObjectMetaData, error) {
	fstat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(&f)
	if err != nil {
		return nil, err
	}
	hash, err := sha256FromBytes(b)
	if err != nil {
		return nil, err
	}
	m := &metadata.ObjectMetaData{
		Name:         f.Name(),
		Checksum:     hash,
		DateModified: fstat.ModTime(),
	}
	return m, nil
}

func (s *Store) Upload(objects []string) error {
	for _, o := range objects {
		objp := path.Join(s.Pantri, o)
		b, err := os.ReadFile(o)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(path.Dir(objp), os.ModePerm); err != nil {
			return err
		}
		if err := os.WriteFile(objp, b, 0644); err != nil {
			return err
		}

		f, err := os.Open(o)
		if err != nil {
			return err
		}

		// generate pantri metadata
		m, err := s.generateMetadata(*f)
		if err != nil {
			return err
		}

		// convert to json
		blob, err := json.MarshalIndent(m, "", " ")
		if err != nil {
			return err
		}
		// write json to pfile
		mpath := path.Join(fmt.Sprintf("%s.pfile", o))
		if err := os.WriteFile(mpath, blob, 0644); err != nil {
			return err
		}
		if s.Opts.RemoveFromSourceRepo != nil {
			if *s.Opts.RemoveFromSourceRepo {
				if err := os.Remove(o); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

var (
	ErrRetrieveFailureHashMismatch = errors.New("hashes don't match, Retrieve aborted")
)

func (s *Store) metadataFromFile(obj string) (*metadata.ObjectMetaData, error) {
	pfile, err := os.Open(fmt.Sprintf("%s.pfile", obj))
	if err != nil {
		return nil, err
	}
	pfileBytes, err := ioutil.ReadAll(pfile)
	if err != nil {
		return nil, err
	}
	var metadata *metadata.ObjectMetaData
	if err := json.Unmarshal(pfileBytes, &metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

func (s *Store) Retrieve(objects []string) error {
	for _, o := range objects {
		f, err := os.Open(path.Join(s.Pantri, o))
		if err != nil {
			return err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		hash, err := sha256FromBytes(b)
		if err != nil {
			return err
		}
		m, err := s.metadataFromFile(o)
		if err != nil {
			return err
		}
		if hash != m.Checksum {
			fmt.Println(hash, m.Checksum)
			return ErrRetrieveFailureHashMismatch
		}
		if err := os.WriteFile(o, b, 0644); err != nil {
			return err
		}

	}
	return nil
}
