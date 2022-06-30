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

	pantriconfig "github.com/discentem/pantri_but_go/config"
	"github.com/discentem/pantri_but_go/metadata"
	"github.com/discentem/pantri_but_go/stores"
	"github.com/mitchellh/go-homedir"
)

type Store struct {
	PantriAddress string         `mapstructure:"pantri_address"`
	Opts          stores.Options `mapstructure:"options"`
}

func (s *Store) init(sourceRepo string) error {
	epa, err := homedir.Expand(s.PantriAddress)
	if err != nil {
		return err
	}

	c := pantriconfig.Config{
		Type:          "local",
		PantriAddress: epa,
		Opts:          s.Opts,
		Validate: func() error {
			// Ensure s.PantriAddress exists before c.WriteToDisk()
			if _, err := os.Stat(s.PantriAddress); err != nil {
				fmt.Println(err)
				return fmt.Errorf("specified pantri_address %q does not exist, so we can't make it a pantri repo", epa)
			}
			return nil
		},
	}

	return c.WriteToDisk(sourceRepo)
}

func New(sourceRepo, pantriAddress string, o stores.Options) (*Store, error) {
	if o.RemoveFromSourceRepo == nil {
		b := false
		o.RemoveFromSourceRepo = &b
	}
	s := &Store{
		PantriAddress: pantriAddress,
		Opts:          o,
	}
	err := s.init(sourceRepo)
	if err != nil {
		return nil, err
	}
	return s, nil
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

func (s *Store) Upload(sourceRepo string, objects ...string) error {
	for _, o := range objects {
		objp := path.Join(s.PantriAddress, o)
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
		pfilePaths := []string{
			path.Join(s.PantriAddress, fmt.Sprintf("%s.pfile", o)),
			path.Join(sourceRepo, fmt.Sprintf("%s.pfile", o)),
		}
		for _, p := range pfilePaths {
			if err := os.WriteFile(p, blob, 0644); err != nil {
				return err
			}
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

func (s *Store) Retrieve(sourceRepo string, objects ...string) error {
	for _, o := range objects {
		f, err := os.Open(path.Join(s.PantriAddress, o))
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
		op := path.Join(sourceRepo, o)
		if err := os.WriteFile(op, b, 0644); err != nil {
			return err
		}

	}
	return nil
}
