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
	gitRepo string
	pantri  string
	opts    stores.Options
}

func NewWithOptions(gitRepo, pantri string, o stores.Options) (*Store, error) {
	if o.RemoveFromRepo == nil {
		b := false
		o.RemoveFromRepo = &b
	}
	s := &Store{
		gitRepo: gitRepo,
		pantri:  pantri,
		opts:    o,
	}
	err := s.init()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func New(gitRepo, pantri string, removeFromRepo *bool) (*Store, error) {
	return NewWithOptions(gitRepo, pantri, stores.Options{
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
		f, err := os.Open(path.Join(s.gitRepo, o))
		if err != nil {
			return err
		}
		defer f.Close()
		objp := path.Join(s.pantri, o)
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		if err := os.WriteFile(objp, b, 0644); err != nil {
			return err
		}
		mf, err := os.Open(objp)
		if err != nil {
			return err
		}
		defer mf.Close()

		// generate pantri metadata
		m, err := s.generateMetadata(*mf)
		if err != nil {
			return err
		}

		// convert to json
		blob, err := json.MarshalIndent(m, "", " ")
		if err != nil {
			return err
		}
		// write json to pfile
		mpath := path.Join(s.gitRepo, fmt.Sprintf("%s.pfile", o))
		if err := os.WriteFile(mpath, blob, 0644); err != nil {
			return err
		}

		if *s.opts.RemoveFromRepo {
			if err := os.Remove(path.Join(s.gitRepo, o)); err != nil {
				return err
			}
		}
	}
	return nil
}

var (
	ErrRetrieveFailureHashMismatch = errors.New("hashes don't match, Retrieve aborted")
)

func (s *Store) Retrieve(objects []string) error {
	for _, o := range objects {
		pfile, err := os.Open(path.Join(s.gitRepo, fmt.Sprintf("%s.pfile", o)))
		if err != nil {
			return err
		}
		pfileBytes, err := ioutil.ReadAll(pfile)
		if err != nil {
			return err
		}
		var metadata *metadata.ObjectMetaData
		if err := json.Unmarshal(pfileBytes, &metadata); err != nil {
			return err
		}

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
		if hash != metadata.Checksum {
			fmt.Println(hash, metadata.Checksum)
			return ErrRetrieveFailureHashMismatch
		}

		objp := path.Join(s.gitRepo, o)
		if err != nil {
			return err
		}
		if err := os.WriteFile(objp, b, 0644); err != nil {
			return err
		}

	}
	return nil
}
