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
	sourceRepo string
	pantri     string
	opts       stores.Options
}

func NewWithOptions(sourceRepo, pantri string, o stores.Options) (*Store, error) {
	if o.RemoveFromRepo == nil {
		b := false
		o.RemoveFromRepo = &b
	}
	s := &Store{
		sourceRepo: sourceRepo,
		pantri:     pantri,
		opts:       o,
	}
	err := s.init()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func New(sourceRepo, pantri string, removeFromRepo *bool) (*Store, error) {
	return NewWithOptions(sourceRepo, pantri, stores.Options{
		RemoveFromRepo: removeFromRepo,
	})
}

func (s *Store) init() error {
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
		// open real object in repo
		f, err := os.Open(o)
		if err != nil {
			return err
		}
		defer f.Close()
		objp := path.Join(s.pantri, o)
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(path.Dir(objp), os.ModePerm); err != nil {
			return err
		}
		if err := os.WriteFile(objp, b, 0644); err != nil {
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

		if *s.opts.RemoveFromRepo {
			if err := os.Remove(o); err != nil {
				return err
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
		f, err := os.Open(path.Join(s.pantri, o))
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
