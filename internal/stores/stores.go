package stores

import (
	"context"
	"encoding/json"
	"errors"
	"os"
)

type StoreType int

const (
	StoreTypeUndefined StoreType = iota
	StoreTypeS3
)

var (
	ErrRetrieveFailureHashMismatch = errors.New("hashes don't match, Retrieve aborted")
)

type Store interface {
	Upload(ctx context.Context, objects ...string) error
	Retrieve(ctx context.Context, objects ...string) error
	GetOptions() Options
}

func (s StoreType) String() string {
	switch s {
	case StoreTypeS3:
		return "s3"
	}
	return "undefined"
}

func (s *StoreType) FromString(storeTypeString string) StoreType {
	return map[string]StoreType{
		"undefined": StoreTypeUndefined,
		"s3":        StoreTypeS3,
	}[storeTypeString]
}

func (s StoreType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *StoreType) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	*s = s.FromString(str)
	return nil
}

func openOrCreateFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}
